package project

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	postgres "github.com/TFTPL/AWS-Cost-Calculator/services/postgres"
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
)

func GetListOfProjectByUserId(userId string, projectParam []string) ([]map[string]string, error) {
	uId, _ := strconv.ParseUint(userId, 0, 64)

	projectByUser, err := postgres.GetProjectsByUserId(uId) // userproject table
	if err != nil {
		return nil, err
	}

	var pIds []int64
	for _, j := range projectByUser {
		pIds = append(pIds, j.ProjectId)
	}

	var projectNameList []string
	projectList, err := postgres.GetBulkProjectsByIds(pIds) // project table
	if err != nil {
		return nil, err
	}

	for _, i := range projectList {
		projectNameList = append(projectNameList, i.Name) // list of proj name
	}
	var finalResult []map[string]string

	for k, j := range projectList {
		v := make(map[string]string)
		v["Serial"] = strconv.Itoa(k + 1)
		v["Id"] = fmt.Sprintf("%d", j.ID)
		v["Description"] = j.Description
		v["StartedOn"] = j.CreatedAt.Format("2006-01-02")

		if j.IsActive{
			v["Name"] = j.Name
			v["IsActive"] = "Active"
			finalResult = append(finalResult, v)
			if j.IsPinned{
				v["IsPinned"] = "Pinned"

			} else {
				v["IsPinned"] = "Not Pinned"

			}
		} else {
			v["Name"] = j.Name
			v["IsActive"] = "Not Active"
			finalResult = append(finalResult, v)
			if j.IsPinned{
				v["IsPinned"] = "Pinned"
			} else {
				v["IsPinned"] = "Not Pinned"
			}

		}
	}

	if len(projectParam) == 1 && projectParam[0] == "None" {
		fmt.Println("If Part")
		sort.Slice(finalResult, func(i, j int) bool {
			return finalResult[i]["Id"] < finalResult[j]["Id"]
		})
	} else {

		// create a map to store the index of each key in list A
		indexMap := make(map[string]int)
		for i, key := range projectParam {
			indexMap[key] = i
		}

		// sort the list of maps by comparing the index of each map's Name key in list A
		sort.Slice(finalResult, func(i, j int) bool {
			name1 := finalResult[i]["Name"]
			name2 := finalResult[j]["Name"]
			index1, ok1 := indexMap[name1]
			index2, ok2 := indexMap[name2]
			if !ok1 {
				return false
			}
			if !ok2 {
				return true
			}
			return index1 < index2
		})

		if len(projectParam) < len(projectNameList) && len(projectParam) > 1 {
			finalResult = finalResult[0:len(projectParam)]
		} else if len(projectParam) == 1 {
			finalResult = finalResult[0:1]
		}
		sort.Slice(finalResult, func(i, j int) bool {
			return finalResult[i]["Id"] < finalResult[j]["Id"]
		})
	}
	return finalResult, nil
}

func NewProjectService(userID, projectname, projectdesc string) (int64, error) {
	uId, _ := strconv.ParseInt(userID, 0, 64)
	if projectname == "" {
		return 0, fmt.Errorf("project name is empty")
	}
	projectname = strings.Trim(projectname, " ")
	projectdesc = strings.Trim(projectdesc, " ")

	projects, err := GetListOfAllProjectNamesByUserId(userID)
	if err != nil {
		log.Println("project not found")
		return 0, errors.New("GetListOfAllProjectNamesByUserId()-projects not found")
	}
	for _, existProject := range projects {
		if existProject == projectname {
			log.Println("project name already exist")
			return 0, errors.New("project name is already exist")
		}
	}
	newProject, err := postgres.CreateProject(projectname, projectdesc)
	if err != nil {
		return 0, err
	}

	_, err = postgres.CreateUserProject(uId, newProject.ID)
	if err != nil {
		return 0, err
	}

	return uId, nil
}

func GetProjectByProjectId(pid string, userId string) (types.GetProjectByProjectIdResponse, error) {
	pId, _ := strconv.ParseInt(pid, 0, 64)
	ProjectDetail, err := postgres.ReadProject(pId)
	if err != nil {
		return types.GetProjectByProjectIdResponse{}, err
	}

	return types.GetProjectByProjectIdResponse{
		Name:        ProjectDetail.Name,
		Description: ProjectDetail.Description,
		UserId:      userId,
		IsPinned:    ProjectDetail.IsPinned,
	}, nil
}

func GetGraphData(pid string, filter string) ([]types.GetGraphDataResponse, error) {
	pId, _ := strconv.ParseInt(pid, 0, 64)
	filterInt, _ := strconv.ParseInt(filter, 0, 64)
	GraphData := []*types.CostSum{}
	now := time.Now()
	currentYear := now.Year()
	startOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)
	diff := now.Sub(startOfYear)
	days := int(diff.Hours() / 24)
	daysStr := strconv.Itoa(days)
	var err error
	if filter == strconv.Itoa(currentYear) {
		GraphData, err = postgres.GetServiceCostByProject(pId, daysStr)
		if err != nil {
			return nil, err
		}
	} else {
		GraphData, err = postgres.GetServiceCostByProject(pId, filter) // For Lsat 1 year
		if err != nil {
			return nil, err
		}
	}

	var result []types.GetDateandCost

	for _, j := range GraphData {
		result = append(result, types.GetDateandCost{Time: j.UsageDate,
			Cost: float64(j.AccruedCostMicrodollar) / 1000000000})
	}

	today := time.Now().AddDate(0, 0, -1)

	var dates []time.Time
	for i := 0; i < int(filterInt); i++ {
		dates = append(dates, today.AddDate(0, 0, -i))
	}

	for _, j := range result { // Removing Date that already exsist in result
		for i, d := range dates {
			if d.Format("2006-01-02") == j.Time.Format("2006-01-02") {
				dates = append(dates[:i], dates[i+1:]...)
			}
		}
	}

	for _, j := range dates {
		result = append(result, types.GetDateandCost{Time: j, Cost: 0})
	}

	sort.Slice(result, func(i, j int) bool { // Sorting result by time
		return result[i].Time.Before(result[j].Time)
	})

	var finalResult []types.GetGraphDataResponse
	for _, j := range result {
		finalResult = append(finalResult, types.GetGraphDataResponse{Time: "cost" + j.Time.Format("2006-01-02") + "date", Cost: j.Cost})
	}

	return finalResult, nil
}

func GetGraphDataByServiceName(pid string, servicesParam []string, filter string) (map[string][]types.GetGraphDataResponse, error) {
	pId, _ := strconv.ParseInt(pid, 0, 64)
	var err error
	serviceGraphdata := make(map[string][]types.GetGraphDataResponse)
	filterInt, _ := strconv.ParseInt(filter, 0, 64)
	GraphData := []*types.CostSum{}

	for _, serviceName := range servicesParam {
		if serviceName != "Combined Cost" {
			GraphData, err = postgres.GetServiceCostByProjectByServiceName(pId, serviceName, filter) //GetDateCostByService(pId)
			if err != nil {
				fmt.Println("Data Not Found")
			}

			var temp []types.GetDateandCost
			for _, j := range GraphData {
				temp = append(temp, types.GetDateandCost{Time: j.UsageDate, Cost: float64(j.AccruedCostMicrodollar)})
			}
			today := time.Now().AddDate(0, 0, -1)
			var dates []time.Time
			for i := 0; i < int(filterInt); i++ {
				dates = append(dates, today.AddDate(0, 0, -i))
			}
			for _, j := range temp { // Removing Date that already exsist in result
				for i, d := range dates {
					if d.Format("2006-01-02") == j.Time.Format("2006-01-02") {
						dates = append(dates[:i], dates[i+1:]...)
					}
				}
			}

			for _, j := range dates {
				temp = append(temp, types.GetDateandCost{Time: j, Cost: 0})
			}

			sort.Slice(temp, func(i, j int) bool { // Sorting result by time
				return temp[i].Time.Before(temp[j].Time)
			})
			var result []types.GetGraphDataResponse
			for _, i := range temp {
				result = append(result, types.GetGraphDataResponse{Time: i.Time.Format("2006-01-02"), Cost: i.Cost / 1000000000})
			}

			serviceGraphdata[serviceName] = result
		} else {
			resGraphData, err := GetGraphData(pid, filter)
			if err != nil {
				fmt.Print(errors.New("graphPageHandler: graph data not found"))
			}
			var tempCombined []types.GetGraphDataResponse
			for _, i := range resGraphData {
				newStr := strings.ReplaceAll(i.Time, "cost", "")
				temp_ := strings.ReplaceAll(newStr, "date", "")
				tempCombined = append(tempCombined, types.GetGraphDataResponse{Time: temp_, Cost: i.Cost})
			}
			serviceGraphdata[serviceName] = tempCombined
		}
	}

	return serviceGraphdata, err
}

func GetServiceCostMinMAxTotalAvg(pid string, filter string) (types.GetServiceCostMinMAxTotalAvgResponse, error) {
	pId, _ := strconv.ParseInt(pid, 0, 64)
	Filter, _ := strconv.ParseInt(filter, 0, 64)
	var err error
	var getValues types.GetServiceCostMinMAxTotalAvg
	getValues, err = postgres.GetServiceCostMinMAxTotalAvg(pId, filter)
	if err != nil {
		return types.GetServiceCostMinMAxTotalAvgResponse{}, err
	}

	totalCost := float64(getValues.AccruedCostMicrodollar) / 1000000000
	minimumCost := float64(getValues.MinAccruedCostMicrodollar) / 1000000000
	maximumCost := float64(getValues.MaxAccruedCostMicrodollar) / 1000000000
	averageCost := totalCost / float64(Filter)

	var result = types.GetServiceCostMinMAxTotalAvgResponse{
		TotalCost:   fmt.Sprintf("%.2f $", totalCost),
		MinimumCost: fmt.Sprintf("%.2f $", minimumCost),
		MaximumCost: fmt.Sprintf("%.2f $", maximumCost),
		AverageCost: fmt.Sprintf("%.2f $", averageCost),
	}
	return result, err
}

func GetServiceCostMinMAxTotalAvgByServiceName(pid string, serviceName string) (types.GetServiceCostMinMAxTotalAvgResponse, error) {
	pId, _ := strconv.ParseInt(pid, 0, 64)

	getValues, err := postgres.GetServiceCostMinMAxTotalAvgByServiceName(pId, serviceName)
	if err != nil {
		return types.GetServiceCostMinMAxTotalAvgResponse{}, err
	}

	totalCost := float64(getValues.AccruedCostMicrodollar) / 1000000000
	minimumCost := float64(getValues.MinAccruedCostMicrodollar) / 1000000000
	maximumCost := float64(getValues.MaxAccruedCostMicrodollar) / 1000000000
	averageCost := totalCost / 30

	var result = types.GetServiceCostMinMAxTotalAvgResponse{
		TotalCost:   fmt.Sprintf("%.2f $", totalCost),
		MinimumCost: fmt.Sprintf("%.2f $", minimumCost),
		MaximumCost: fmt.Sprintf("%.2f $", maximumCost),
		AverageCost: fmt.Sprintf("%.2f $", averageCost),
	}

	return result, err
}

func ProjectExport(pid, filter, cloudServices string) (string, []byte, error) {
	pId, _ := strconv.ParseInt(pid, 0, 64)
	interval, _ := strconv.ParseInt(filter, 0, 64)

	project, err := postgres.ReadProject(pId)
	if err != nil {
		fmt.Println("Error in getting Project details for project Id", pId)
		return "", nil, err
	}

	newFileName := fmt.Sprintf("%s-%s-Cost-Last-%s-days.csv", project.Name, cloudServices, filter)

	buf := new(bytes.Buffer)

	writer := csv.NewWriter(buf)

	header := []string{"Services", fmt.Sprintf("%s-Cost", project.Name), "Date"}
	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error in creating header:", err)
		return "", nil, err
	}
	var csvData []*types.DateServiceCost

	if cloudServices == "AWS" {
		csvData, err = postgres.GetServiceCostDateByProjectId(pId, interval)
		if err != nil {
			fmt.Println("Error in getting data for requested project:", err)
			return "", nil, err
		}
	} else if cloudServices == "MongoDb" {
		csvData, err = postgres.GetMongoServiceCostDateByProjectId(pId, interval)
		if err != nil {
			fmt.Println("Error in getting data for requested project:", err)
			return "", nil, err
		}
	}

	for _, row := range csvData {
		record := []string{
			row.ServiceTitle,
			strconv.FormatFloat(float64(row.AccruedCostMicrodollar)/1000000000, 'f', 14, 64),
			row.UsageDate.Format("2006-01-02"),
		}
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error in writing record to file:", err)
			return "", nil, err
		}

	}

	// Flush the data to the file
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing data:", err)
		return "", nil, err
	}

	//return filename, buf.Bytes(), nil
	return newFileName, buf.Bytes(), nil
}

// export graph by list of services
func ProjectExportByService(pid string, filter string, servicesParam []string) (string, []byte, error) {

	var csvData []*types.DateServiceCost

	newServices := []string{}
	pId, _ := strconv.ParseInt(pid, 0, 64)
	interval, _ := strconv.ParseInt(filter, 0, 64)

	project, err := postgres.ReadProject(pId)
	if err != nil {
		fmt.Println("Error in getting Project details for project Id", pId)
		return "", nil, err
	}

	newFileName := fmt.Sprintf("%s-Services-Aws-Cost-Last-%s-days.csv", project.Name, filter)

	buf := new(bytes.Buffer)

	writer := csv.NewWriter(buf)

	header := []string{"Services", fmt.Sprintf("%s-Cost", project.Name), "Date"}

	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error in creating header:", err)
		return "", nil, err
	}

	if len(servicesParam) == 1 && servicesParam[0] == "Combined Cost" {
		csvData, err = postgres.GetServiceCostDateByProjectId(pId, interval)
		if err != nil {
			fmt.Println("Error in getting data for requested project:", err)
			return "", nil, err
		}
	}

	for _, val := range servicesParam {
		if val != "Combined Cost" {
			newServices = append(newServices, val)
			csvData, err = postgres.GetCostDateByProjectIdAndService(pId, interval, newServices)
			if err != nil {
				fmt.Println("Error in getting data for requested services:", err)
				return "", nil, err
			}
		}
	}

	for _, row := range csvData {
		record := []string{
			row.ServiceTitle,
			strconv.FormatFloat(float64(row.AccruedCostMicrodollar)/1000000000, 'f', 14, 64),
			row.UsageDate.Format("2006-01-02"),
		}
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error in writing record to file:", err)
			return "", nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing data:", err)
		return "", nil, err
	}

	return newFileName, buf.Bytes(), nil
}

func ProjectExportByDate(pid, fromdate, toDate, cloudServices string) (string, []byte, error) {
	pId, _ := strconv.ParseInt(pid, 0, 64)
	date, _ := time.Parse("2006-01-02", toDate)
	nextDate := date.AddDate(0, 0, 1)
	todate := nextDate.Format("2006-01-02")

	var csvData []*types.DateServiceCost

	project, err := postgres.ReadProject(pId)
	if err != nil {
		fmt.Println("Error in getting Project details for project Id", pId)
		return "", nil, err
	}

	newFileName := fmt.Sprintf("%s-%s-Aws-Cost-From-%s-To-%s.csv", project.Name, cloudServices, fromdate, toDate)

	buf := new(bytes.Buffer)

	writer := csv.NewWriter(buf)

	header := []string{"Services", fmt.Sprintf("%s-Cost", project.Name), "Date"}

	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error in creating header:", err)
		return "", nil, err
	}

	if cloudServices == "AWS" {
		csvData, err = postgres.GetServiceCostDateByProjectIdAndDate(pId, fromdate, todate)
		if err != nil {
			fmt.Println("Error in getting data:", err)
			return "", nil, err
		}
	} else if cloudServices == "MongoDb" {
		csvData, err = postgres.GetMongoServiceCostDateByProjectIdAndDate(pId, fromdate, todate)
		if err != nil {
			fmt.Println("Error in getting data:", err)
			return "", nil, err
		}
	}

	for _, row := range csvData {
		record := []string{
			row.ServiceTitle,
			strconv.FormatFloat(float64(row.AccruedCostMicrodollar)/1000000000, 'f', 14, 64),
			row.UsageDate.Format("2006-01-02"),
		}
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error in writing record to file:", err)
			return "", nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing data:", err)
		return "", nil, err
	}

	return newFileName, buf.Bytes(), nil
}

func ProjectExportByYear(pid, year, cloudServices string) (string, []byte, error) {
	pId, _ := strconv.ParseInt(pid, 0, 64)

	var csvData []*types.DateServiceCost

	project, err := postgres.ReadProject(pId)
	if err != nil {
		fmt.Println("Error in getting Project details for project Id", pId)
		return "", nil, err
	}

	newFileName := fmt.Sprintf("%s-%s-Aws-Cost-Of-Year-%s.csv", project.Name, cloudServices, year)

	buf := new(bytes.Buffer)

	writer := csv.NewWriter(buf)

	header := []string{"Services", fmt.Sprintf("%s-Cost", project.Name), "Date"}

	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error in creating header:", err)
		return "", nil, err
	}

	if cloudServices == "AWS" {
		csvData, err = postgres.GetServiceCostDateByProjectIdAndYear(pId, year)
		if err != nil {
			fmt.Println("Error in getting data:", err)
			return "", nil, err
		}
	} else if cloudServices == "MongoDb" {
		csvData, err = postgres.GetMongoServiceCostDateByProjectIdAndYear(pId, year)
		if err != nil {
			fmt.Println("Error in getting data:", err)
			return "", nil, err
		}
	}

	for _, row := range csvData {
		record := []string{
			row.ServiceTitle,
			strconv.FormatFloat(float64(row.AccruedCostMicrodollar)/1000000000, 'f', 14, 64),
			row.UsageDate.Format("2006-01-02"),
		}
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error in writing record to file:", err)
			return "", nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing data:", err)
		return "", nil, err
	}

	return newFileName, buf.Bytes(), nil
}

func UpdateProjectsDetailsByProjectId(pid, projectname, description, awsaccesskey, awssecretkey string) error {
	pId, _ := strconv.ParseInt(pid, 0, 64)
	ProjectDetail, err := postgres.ReadProject(pId)
	if err != nil {
		return err
	}
	AwsCredentials, err := postgres.ReadAwsCredentials(pId)
	if err != nil {
		return err
	}
	if len(projectname) == 0 {
		projectname = ProjectDetail.Name
	}
	if len(description) == 0 {
		description = ProjectDetail.Description
	}
	if len(awsaccesskey) == 0 {
		awsaccesskey = AwsCredentials.AccessKey
	}
	if len(awssecretkey) == 0 {
		awssecretkey = AwsCredentials.AwsSecretKey
	}

	_, err = postgres.UpdateProject(pId, projectname, description)
	if err != nil {
		return err
	}
	_, err = postgres.UpdateAwsCredentials(pId, awssecretkey, awsaccesskey)
	if err != nil {
		return err
	}

	return nil
}

// update pinned project
func UpdatePinProjects(Pid, IsPinned string) error {

	pid, _ := strconv.ParseInt(Pid, 0, 64)
	var valIsPinned bool
	if IsPinned == "Pinned" {
		valIsPinned = false
	} else {
		valIsPinned = true
	}

	err := postgres.UpdateProjectbyPin(pid, valIsPinned)
	if err != nil {
		return err
	}

	return nil

}

// get all Active pinned project
func GetAllActivePinnedProjects(uid string) ([]types.PinnedProjectResponse, error) {
	userId, _ := strconv.ParseInt(uid, 0, 64)

	projects, err := postgres.GetProjectsByUserId(uint64(userId))
	if err != nil {
		return nil, errors.New("pid not found")
	}
	var pIds []int64
	for _, j := range projects {
		pIds = append(pIds, j.ProjectId)
	}
	pinProjectResponse, err := postgres.GetAllPinnedProjects(pIds)

	if err != nil {
		return nil, errors.New("pinned project not found")
	}

	return pinProjectResponse, nil
}

func GetDistinctServiceTitlesForProject(pId string) ([]string, error) {
	projectId, _ := strconv.ParseInt(pId, 0, 64)

	projectServices, err := postgres.GetDistinctServiceTitlesForProject(projectId) // For Distinct Service Titles of the Project
	if err != nil {
		return nil, errors.New("pid not found")
	}

	projectServices = append(projectServices, "Combined Cost")

	return projectServices, err
}

func DailyTotalCostOfAllProjectsOfUser(userid string, daysFilter string) ([]types.TotalCostSumResonse, error) {
	uid, err := strconv.ParseInt(userid, 0, 64)
	if err != nil {
		return nil, err
	}
	dayFilterInt, err := strconv.ParseInt(daysFilter, 0, 64)
	if err != nil {
		return nil, err
	}

	resUserProject, err := postgres.GetProjectsIdsByUserId(uid)
	if err != nil {
		return nil, err
	}

	// Create a map to keep track of the daily cost sums for each project
	dailyCostSums := make(map[string]int64)

	// Get the date range for the last 30 days
	today := time.Now().AddDate(0, 0, -1)
	var dates []time.Time
	for i := 0; i < int(dayFilterInt); i++ {
		dates = append(dates, today.AddDate(0, 0, -i))
	}

	// Loop through each project and accumulate the daily cost sums
	for _, userProject := range resUserProject {
		resProjectCost, err := postgres.GetTotalCostOfProjects(userProject.ProjectId, daysFilter)
		if err != nil {
			return nil, err
		}

		for _, costSum := range resProjectCost {
			dateStr := costSum.UsageDate.Format("2006-01-02")
			dailyCostSums[dateStr] += costSum.AccruedCostMicrodollar
		}
	}

	// Create the final result array with daily cost sums for each date
	var Result []types.TotalCostSumResonse
	for _, date := range dates {
		dateStr := date.Format("2006-01-02")

		cost, ok := dailyCostSums[dateStr]
		if !ok {
			// If there is no data for the date, set the cost to zero
			cost = 0
		}

		Result = append(Result, types.TotalCostSumResonse{
			UsageDate:              "cost" + dateStr + "date",
			AccruedCostMicrodollar: float64(cost) / 1000000000,
		})
	}

	var finalResult []types.TotalCostSumResonse
	for i := len(Result) - 1; i >= 0; i-- {
		finalResult = append(finalResult, Result[i])
	}

	return finalResult, nil
}

func DeleteProject(pId string) error {
	projectId, _ := strconv.ParseInt(pId, 0, 64)

	_, err := postgres.DeleteProject(projectId)
	if err != nil {
		return err
	}

	// _, err = postgres.DeleteAwsCredentials(projectId)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func UpdatePinProjectsByProjectName(uid, pName, IsPinned string) error {
	var valIsPinned bool
	if IsPinned == "true" {
		valIsPinned = false
	} else {
		valIsPinned = true
	}

	err := postgres.UpdateProjectbyProjectNmae(pName, valIsPinned)
	if err != nil {
		return err
	}

	return nil

}

func ProjectExportPageHandlerByServiceAndDate(pidRequest, fromDate, toDate string, services []string) (string, []byte, error) {

	pId, _ := strconv.ParseInt(pidRequest, 0, 64)
	date, _ := time.Parse("2006-01-02", toDate)
	nextDate := date.AddDate(0, 0, 1)
	todate := nextDate.Format("2006-01-02")

	project, err := postgres.ReadProject(pId)
	if err != nil {
		fmt.Println("Error in getting Project details for project Id", pId)
		return "", nil, err
	}

	newFileName := fmt.Sprintf("%s-Services-Aws-Cost-From-%s-To-%s.csv", project.Name, fromDate, toDate)

	buf := new(bytes.Buffer)

	writer := csv.NewWriter(buf)

	header := []string{"Services", fmt.Sprintf("%s-Cost", project.Name), "Date"}

	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error in creating header:", err)
		return "", nil, err
	}

	csvData, err := postgres.GetServiceCostDateByProjectIdAndServicesAndCustomDate(pId, fromDate, todate, services)

	if err != nil {
		fmt.Println("Error in getting data:", err)
		return "", nil, err
	}

	for _, row := range csvData {
		record := []string{
			row.ServiceTitle,
			strconv.FormatFloat(float64(row.AccruedCostMicrodollar)/1000000000, 'f', 14, 64),
			row.UsageDate.Format("2006-01-02"),
		}
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error in writing record to file:", err)
			return "", nil, err
		}
	}

	// Flush the data to the file
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing data:", err)
		return "", nil, err
	}

	return newFileName, buf.Bytes(), nil
}

func ProjectExportPageHandlerByServiceWithYear(pidRequest, year string, services []string) (string, []byte, error) {
	pId, _ := strconv.ParseInt(pidRequest, 0, 64)

	project, err := postgres.ReadProject(pId)
	if err != nil {
		fmt.Println("Error in getting Project details for project Id", pId)
		return "", nil, err
	}

	newFileName := fmt.Sprintf("%s-Services-Aws-Cost-Of-Year-%s.csv", project.Name, year)

	buf := new(bytes.Buffer)

	writer := csv.NewWriter(buf)

	header := []string{"Services", fmt.Sprintf("%s-Cost", "project.Name"), "Date"}

	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error in creating header:", err)
		return "", nil, err
	}

	csvData, err := postgres.GetServiceCostDateByProjectIdWithServiceAndYear(pId, year, services)
	if err != nil {
		fmt.Println("Error in getting data:", err)
		return "", nil, err
	}

	for _, row := range csvData {
		record := []string{
			row.ServiceTitle,
			strconv.FormatFloat(float64(row.AccruedCostMicrodollar)/1000000000, 'f', 14, 64),
			row.UsageDate.Format("2006-01-02"),
		}
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error in writing record to file:", err)
			return "", nil, err
		}

	}

	// Flush the data to the file
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing data:", err)
		return "", nil, err
	}

	return newFileName, buf.Bytes(), nil

}

// func for getting total cost of last 30 days of list project by users.
func GetTotalCostOfProjectList(userid string) ([]types.TotalCostProject, error) {
	uId, _ := strconv.ParseUint(userid, 0, 64)

	projectByUser, err := postgres.GetProjectsByUserId(uId)
	if err != nil {
		return nil, err
	}

	var pIds []int64
	for _, j := range projectByUser {
		pIds = append(pIds, j.ProjectId)
	}

	activeProjectList, err := postgres.GetAllActiveProjectsByProjectIds(pIds)
	if err != nil {
		return nil, err
	}
	var projectwithCost []types.TotalCostProject
	for _, pid := range activeProjectList {

		totalCostOfAProject, err := postgres.GetTotalCostOfProjectList(pid.ID)
		if err != nil {
			return nil, err
		}
		s := fmt.Sprintf("%v", float64(totalCostOfAProject.AccruedCostMicrodollar)/1000000000)
		projectwithCost = append(projectwithCost, types.TotalCostProject{
			ProjectName:            pid.Name,
			AccruedCostMicrodollar: s + "split",
		})

	}
	return projectwithCost, err

}

func UpdateProjectDetails(pid, projectname, description string) error {
	pId, _ := strconv.ParseInt(pid, 0, 64)
	ProjectDetail, err := postgres.ReadProject(pId)
	if err != nil {
		return err
	}
	if len(projectname) == 0 {
		projectname = ProjectDetail.Name
	}
	if len(description) == 0 {
		description = ProjectDetail.Description
	}

	_, err = postgres.UpdateProject(pId, projectname, description)
	if err != nil {
		return err
	}

	return nil
}

// get aws credentials from database [access key, secretkey]
func GetAWSCredentials(pid string) (bool, error) {
	pId, _ := strconv.ParseInt(pid, 0, 64)

	AwsKeys, err := postgres.ReadAwsCredentials(pId)

	if err != nil {
		return false, err
	}
	if len(AwsKeys.AccessKey) != 0 && len(AwsKeys.AwsSecretKey) != 0 {
		return true, nil
	}
	return false, err
}

// currently unused but used in later
func UpsertAwsCredentials(pid, accessKey, secretKey string) error {
	pId, _ := strconv.ParseInt(pid, 0, 64)

	authkeys, err := CheckAWSKeys(accessKey, secretKey)

	if !authkeys && err != nil {
		return err
	}

	AwsKeys, err := postgres.ReadAwsCredentials(pId)
	if err != nil {
		_, err = postgres.CreateAwsCredentials(pId, accessKey, secretKey)
		if err != nil {
			return err
		}
		return nil
	}
	if AwsKeys.ProjectId == pId {
		if len(accessKey) == 0 {
			accessKey = AwsKeys.AccessKey
		}
		if len(secretKey) == 0 {
			secretKey = AwsKeys.AwsSecretKey
		}
		_, err = postgres.UpdateAwsCredentials(pId, secretKey, accessKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetListOfAllProjectNamesByUserId(userId string) ([]string, error) {
	uId, _ := strconv.ParseUint(userId, 0, 64)

	projectByUser, err := postgres.GetProjectsByUserId(uId)
	if err != nil {
		return nil, err
	}

	var pIds []int64
	for _, j := range projectByUser {
		pIds = append(pIds, j.ProjectId)
	}

	projectList, err := postgres.GetBulkProjectsByIds(pIds)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, i := range projectList {
		result = append(result, i.Name)
	}

	return result, err

}

// GetCostOfProjectsForDashboard takes the projctnames from dashboard , and return the totalcost of these project

func stringInSlice(str string, list []string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}

func GetCostOfProjectsForDashboard(userId string, filter string, pJName []string) (map[string][]types.GetGraphDataResponse, error) {

	var (
		resultdetails []types.Projectdetail
		elementChk    = make(map[string]bool)
		finalresult   = make(map[string][]types.GetGraphDataResponse)
	)

	if len(pJName) == 1 && pJName[0] == "All Projects" {
		resDailyCostSum, err := DailyTotalCostOfAllProjectsOfUser(userId, filter)
		if err != nil {
			fmt.Print(errors.New("graphPageHandler: data not found"))
		}
		var result []types.GetGraphDataResponse
		for _, j := range resDailyCostSum {
			result = append(result, types.GetGraphDataResponse{Time: j.UsageDate, Cost: j.AccruedCostMicrodollar})
		}

		finalresult["All Projects"] = result
	} else {
		uId, _ := strconv.ParseUint(userId, 0, 64)
		dayFilter, _ := strconv.ParseInt(filter, 0, 64)

		projectByUser, err := postgres.GetProjectsByUserId(uId)
		if err != nil {
			fmt.Println("Error in fetching project list for user", err)
			return nil, err
		}

		var projectIds []int64
		for _, pJ := range projectByUser {
			projectIds = append(projectIds, pJ.ProjectId)
		}

		projectDetails, err := postgres.GetBulkProjectsByIds(projectIds)
		if err != nil {
			fmt.Println("Error in fetching project details for projectsIds", err)
			return nil, err
		}

		for _, val := range pJName {
			elementChk[val] = true
		}

		for _, proj := range projectDetails {
			if elementChk[proj.Name] {
				detail := types.Projectdetail{Name: proj.Name, Id: proj.ID}
				resultdetails = append(resultdetails, detail)
			}
		}

		for _, proj := range resultdetails {
			res, err := postgres.GetTotalCostOfProjects(proj.Id, filter)
			if err != nil {
				fmt.Printf("Error in getting project details for project: %d error is: %v", proj.Id, err)
			}
			var temp []types.GetDateandCost

			for _, j := range res {
				temp = append(temp, types.GetDateandCost{Time: j.UsageDate, Cost: float64(j.AccruedCostMicrodollar)})
			}

			today := time.Now().AddDate(0, 0, -1)
			var dates []time.Time
			for i := 0; i < int(dayFilter); i++ {
				dates = append(dates, today.AddDate(0, 0, -i))
			}
			for _, j := range temp {
				for i, d := range dates {
					if d.Format("2006-01-02") == j.Time.Format("2006-01-02") {
						dates = append(dates[:i], dates[i+1:]...)
					}
				}
			}

			for _, j := range dates {
				temp = append(temp, types.GetDateandCost{Time: j, Cost: 0})
			}

			sort.Slice(temp, func(i, j int) bool {
				return temp[i].Time.Before(temp[j].Time)
			})
			var result []types.GetGraphDataResponse

			for _, i := range temp {
				result = append(result, types.GetGraphDataResponse{Time: i.Time.Format("2006-01-02"), Cost: i.Cost / 1000000000})
				finalresult[proj.Name] = result
			}

		}
	}

	if stringInSlice("All Projects", pJName) {
		resDailyCostSum, err := DailyTotalCostOfAllProjectsOfUser(userId, filter)
		if err != nil {
			fmt.Print(errors.New("graphPageHandler: data not found"))
		}
		var result []types.GetGraphDataResponse
		for _, j := range resDailyCostSum {
			result = append(result, types.GetGraphDataResponse{Time: j.UsageDate, Cost: j.AccruedCostMicrodollar})
		}

		finalresult["All Projects"] = result
	}
	return finalresult, nil
}

func UpdateActiveProject(Pid, IsActive string) error {
	pid, _ := strconv.ParseInt(Pid, 0, 64)
	var valIsActive bool
	if IsActive == "Active" {
		valIsActive = false
	} else {
		valIsActive = true
	}

	err := postgres.UpdateProjectbyActiveStatus(pid, valIsActive)
	if err != nil {
		return err
	}

	return nil

}

func CreateAwsCredentials(pid, accessKey, secretKey string) error {
	pId, _ := strconv.ParseInt(pid, 0, 64)

	authkeys, err := CheckAWSKeys(accessKey, secretKey)

	if !authkeys && err != nil {
		return err
	}

	_, err = postgres.CreateAwsCredentials(pId, accessKey, secretKey)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAwsCredentials(pId string) error {
	projectId, err := strconv.ParseInt(pId, 0, 64)
	if err != nil {
		log.Println("type conversion failed", err)
		return errors.New("DeleteAwsCredentials(): type conversion failed")
	}
	err = postgres.DeleteAwsCredentials(projectId)
	if err != nil {
		log.Println("postgres:DeleteAwsCredentials() delete credentials failed", err)
		return err
	}
	return nil
}

func UpdateByGraphPinned(Pid, IsPinned string) error {

	pid, _ := strconv.ParseInt(Pid, 0, 64)
	var valIsPinned bool
	if IsPinned == "true" {
		valIsPinned = false
	} else {
		valIsPinned = true
	}

	err := postgres.UpdateProjectbyPin(pid, valIsPinned)
	if err != nil {
		return err
	}

	return nil

}

func GetMongoGraphData(pid string, filter string) ([]types.GetGraphDataResponse, error) {
	pId, _ := strconv.ParseInt(pid, 0, 64)
	filterInt, _ := strconv.ParseInt(filter, 0, 64)
	GraphData := []*types.CostSum{}
	now := time.Now()
	currentYear := now.Year()
	startOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)
	diff := now.Sub(startOfYear)
	days := int(diff.Hours() / 24)
	daysStr := strconv.Itoa(days)
	var err error
	if filter == strconv.Itoa(currentYear) {
		GraphData, err = postgres.GetMongoServiceCostByProject(pId, daysStr)
		if err != nil {
			return nil, err
		}
	} else {
		GraphData, err = postgres.GetMongoServiceCostByProject(pId, filter) // For Lsat 1 year
		if err != nil {
			return nil, err
		}
	}

	var result []types.GetDateandCost

	for _, j := range GraphData {
		result = append(result, types.GetDateandCost{Time: j.UsageDate,
			Cost: float64(j.AccruedCostMicrodollar) / 1000000000})
	}

	today := time.Now().AddDate(0, 0, -1)

	var dates []time.Time
	for i := 0; i < int(filterInt); i++ {
		dates = append(dates, today.AddDate(0, 0, -i))
	}

	for _, j := range result { // Removing Date that already exsist in result
		for i, d := range dates {
			if d.Format("2006-01-02") == j.Time.Format("2006-01-02") {
				dates = append(dates[:i], dates[i+1:]...)
			}
		}
	}

	for _, j := range dates {
		result = append(result, types.GetDateandCost{Time: j, Cost: 0})
	}

	sort.Slice(result, func(i, j int) bool { // Sorting result by time
		return result[i].Time.Before(result[j].Time)
	})

	var finalResult []types.GetGraphDataResponse
	for _, j := range result {
		finalResult = append(finalResult, types.GetGraphDataResponse{Time: "cost" + j.Time.Format("2006-01-02") + "date", Cost: j.Cost})
	}

	return finalResult, nil
}

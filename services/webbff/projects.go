package webbff

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	services "github.com/TFTPL/AWS-Cost-Calculator/services/coreservices/project"
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/gorilla/mux"
)

func dashboardHandler(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	userId, ok := val["uid"]
	daysFilter, ok := val["filter"]
	projectParam := request.URL.Query()["projects"]
	if projectParam == nil {
		projectParam = append(projectParam, "All Projects")
	}

	if ok {
		// TODO show welcome username
		// TODO display pinned projects
		// TODO display user all projects resource graph

		pinnedProjectResponse, err := services.GetAllActivePinnedProjects(userId)
		if err != nil {
			fmt.Print(errors.New("update pin: error to fetch data"))
		}

		resDailyCostSum, err := services.DailyTotalCostOfAllProjectsOfUser(userId, daysFilter)
		if err != nil {
			fmt.Print(errors.New("graphPageHandler: data not found"))
		}
		GetTotalCostOfProjectList, err := services.GetTotalCostOfProjectList(userId)
		if err != nil {
			fmt.Print(errors.New("graphPageHandler: data not found"))
			fmt.Println(err)
		}

		projectListResponse, err := services.GetListOfAllProjectNamesByUserId(userId)
		if err != nil {
			fmt.Print(errors.New("update pin: error to fetch data"))
		}

		CostOfReqProjects, err := services.GetCostOfProjectsForDashboard(userId, daysFilter, projectParam)
		if err != nil {
			fmt.Println("Error in generating Totalcost of required Projects", err)
		}

		projectListResponse = append(projectListResponse, "All Projects")

		var Filter string
		var check = ""
		if daysFilter == "7" {
			Filter = "Last 7 days"
			check = "Last 7 days"
		}

		if daysFilter == "30" {
			Filter = "Last 30 days"
			check = "Last 30 days"
		}

		now := time.Now()
		currentYear := now.Year()
		lastYear := now.Year() - 1
		lastYearStr := strconv.Itoa(lastYear)
		currentYearStr := strconv.Itoa(currentYear)

		if daysFilter == "365" {

			Filter = "Last Year(" + lastYearStr + "-" + currentYearStr + ")"
			check = "Last 1 Year"
		}

		t := template.Must(template.ParseFiles(
			"services/webbff/templates/dashboard/dashboard.html",
			"services/webbff/templates/head/head.html",
			"services/webbff/templates/head/nucleo.html",
			"services/webbff/templates/dashboard/import.html",
			"services/webbff/templates/dashboard/title.html",
			"services/webbff/templates/dashboard/header.html",
			"services/webbff/templates/dashboard/main.html",
			"services/webbff/templates/footer/footer.html",
			"services/webbff/templates/script/dashboardGraphScript.html",
			"services/webbff/templates/script/script.html"))
		t.Execute(response, struct {
			UserId                       string
			PinnedProjectResponse        []types.PinnedProjectResponse
			DailyCostSum                 []types.TotalCostSumResonse
			TotalCostOfProjectList       []types.TotalCostProject
			GraphTag                     string
			Check                        string
			ProjectList                  []string
			FilterDays                   string
			SelectedProjects             []string
			TotalCostOfRequestedProjects map[string][]types.GetGraphDataResponse
		}{
			UserId:                       userId,
			PinnedProjectResponse:        pinnedProjectResponse,
			DailyCostSum:                 resDailyCostSum,
			TotalCostOfProjectList:       GetTotalCostOfProjectList,
			GraphTag:                     Filter,
			Check:                        check,
			ProjectList:                  projectListResponse,
			FilterDays:                   daysFilter,
			SelectedProjects:             projectParam,
			TotalCostOfRequestedProjects: CostOfReqProjects,
		})

	}

}

// func newProjectPageHandler(response http.ResponseWriter, request *http.Request) {

// 	if request.Method == "POST" {
// 		projectname := request.FormValue("projectname")
// 		projectdesc := request.FormValue("projectdesc")

// 		val := mux.Vars(request)
// 		userId := val["uid"]

// 		_, err := services.NewProjectService(userId, projectname, projectdesc)

// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		http.Redirect(response, request, fmt.Sprintf("/%s/projects/", userId), http.StatusFound)

// 	}
// }

func projectSettingsPageHandler(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	userId := val["uid"]
	projectId := val["pid"]
	if request.Method == "POST" {
		projectname := request.FormValue("projectname")
		projectdesc := request.FormValue("projectdesc")
		awsaccesskey := request.FormValue("awsaccesskey")
		awssecretkey := request.FormValue("awssecretkey")

		err := services.UpdateProjectsDetailsByProjectId(projectId, projectname, projectdesc, awsaccesskey, awssecretkey)

		if err != nil {
			fmt.Println(err)
		}

		http.Redirect(response, request, fmt.Sprintf("/%s/projects", userId), http.StatusFound)

	}
}

func getProjectByUserIdHandler(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	userId, ok := val["uid"]
	projectParam := request.URL.Query()["projects"]
	var searchProjects []string
	if projectParam != nil {
		queries := strings.Split(projectParam[0], ",")

		searchProjects = append(searchProjects, queries...)
	} else {
		searchProjects = append(searchProjects, "None")
		projectParam = append(projectParam, "None")
	}
	//fmt.Println(len(searchProjects[1]))
	if ok {

		res, err := services.GetListOfProjectByUserId(userId, searchProjects)
		if err != nil {
			fmt.Println(err)
		}
		t := template.Must(template.ParseFiles(
			"services/webbff/templates/projects/project.html",
			"services/webbff/templates/head/head.html",
			"services/webbff/templates/head/nucleo.html",
			"services/webbff/templates/projects/import.html",
			"services/webbff/templates/projects/title.html",
			"services/webbff/templates/projects/header.html",
			"services/webbff/templates/projects/main.html",
			"services/webbff/templates/footer/footer.html"))
		t.Execute(response, struct {
			Projects         []map[string]string
			UserId           string
			SelectedProjects string
		}{
			Projects:         res,
			UserId:           userId,
			SelectedProjects: projectParam[0],
		},
		)
	}
}

func graphPageHandler(response http.ResponseWriter, request *http.Request) {
	// response.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	// response.Header().Set("Access-Control-Allow-Origin", "*")
	// response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	val := mux.Vars(request)
	pidRequest := val["pid"]
	uIdRequest := val["uid"]
	filter := val["filter"]
	var projectParam []string
	projectParam = append(projectParam, "None")

	resProjectDetail, err := services.GetProjectByProjectId(pidRequest, uIdRequest)
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: project data not found"))
	}

	resGraphData, err := services.GetGraphData(pidRequest, filter)
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: graph data not found"))
	}

	resUserProjects, err := services.GetListOfProjectByUserId(uIdRequest, projectParam)
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: data not found"))
	}

	resServiceCostDifferentAspects, err := services.GetServiceCostMinMAxTotalAvg(pidRequest, filter)
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: data not found"))
	}

	resGetProjectServices, err := services.GetDistinctServiceTitlesForProject(pidRequest) // For Distinct Service Titles of the Project
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: data not found"))
	}

	resMongoGraphData, err := services.GetMongoGraphData(pidRequest, filter)
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: data not found"))
	}

	var Filter string
	var check = ""
	if filter == "7" {
		Filter = "Last 7 days"
		check = "Last 7 days"
	}

	if filter == "30" {
		Filter = "Last 30 days"
		check = "Last 30 days"
	}

	now := time.Now()
	currentYear := now.Year()
	lastYear := now.Year() - 1
	lastYearStr := strconv.Itoa(lastYear)
	currentYearStr := strconv.Itoa(currentYear)

	if filter == currentYearStr {
		Filter = "Current Year"
		check = "Current Year"
	}

	if filter == "365" {

		Filter = "Last Year(" + lastYearStr + "-" + currentYearStr + ")"
		check = "Last 1 Year"
	}

	t := template.Must(template.ParseFiles(
		"services/webbff/templates/graph/pgraph.html",
		"services/webbff/templates/head/head.html",
		"services/webbff/templates/graph/title.html",
		"services/webbff/templates/graph/import.html",
		"services/webbff/templates/head/nucleo.html",
		"services/webbff/templates/graph/aside.html",
		"services/webbff/templates/graph/pgmain.html",
		"services/webbff/templates/footer/footer.html",
		"services/webbff/templates/script/graphScript.html",
	))
	t.Execute(response, struct {
		ProjectDetail                types.GetProjectByProjectIdResponse
		GraphData                    []types.GetGraphDataResponse
		UserId                       string
		ProjectId                    string
		UserProjects                 []map[string]string
		GetServiceCostMinMAxTotalAvg types.GetServiceCostMinMAxTotalAvgResponse
		ProjectServiceList           []string
		GraphTag                     string
		Year                         string
		Check                        string
		Url                          string
		MongoUrl                     string
		MongoGraphData               []types.GetGraphDataResponse
	}{ProjectDetail: resProjectDetail,
		UserId:                       uIdRequest,
		GraphData:                    resGraphData,
		ProjectId:                    pidRequest,
		UserProjects:                 resUserProjects,
		GetServiceCostMinMAxTotalAvg: resServiceCostDifferentAspects,
		ProjectServiceList:           resGetProjectServices,
		GraphTag:                     Filter,
		Year:                         currentYearStr,
		Check:                        check,
		Url:                          fmt.Sprintf("/%s/projects/%s/%s/export", uIdRequest, pidRequest, filter),
		MongoGraphData:               resMongoGraphData,
		MongoUrl:                     fmt.Sprintf("/%s/projects/%s/%s/mongoExport", uIdRequest, pidRequest, filter),
	})

}

func graphPageServiceHandler(response http.ResponseWriter, request *http.Request) {
	servicesParam := request.URL.Query()["services"]
	val := mux.Vars(request)
	pidRequest := val["pid"]
	uIdRequest := val["uid"]
	filter := val["filter"]
	var projectParam []string
	projectParam = append(projectParam, "None")

	if servicesParam == nil || (len(servicesParam) == 1 && servicesParam[0] == "Combined Cost") {
		http.Redirect(response, request, fmt.Sprintf("/%s/projects/%s/30", uIdRequest, pidRequest), http.StatusFound)

	}

	resProjectDetail, err := services.GetProjectByProjectId(pidRequest, uIdRequest)
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: project data not found"))
	}
	resGraphData, err := services.GetGraphDataByServiceName(pidRequest, servicesParam, filter)
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: graph data not found"))
	}
	resUserProjects, err := services.GetListOfProjectByUserId(uIdRequest, projectParam)
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: data not found"))
	}

	resServiceCostDifferentAspects, err := services.GetServiceCostMinMAxTotalAvg(pidRequest, "30")
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: data not found"))
	}

	resGetProjectServices, err := services.GetDistinctServiceTitlesForProject(pidRequest) // For Distinct Service Titles of the Project
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: data not found"))
	}
	var testService []string

	for i := range servicesParam {
		testService = append(testService, servicesParam[i]+"/")

	}

	var Filter string
	var check = ""
	if filter == "7" {
		Filter = "Last 7 days"
		check = "Last 7 days"
	}

	if filter == "30" {
		Filter = "Last 30 days"
		check = "Last 30 days"
	}

	now := time.Now()
	currentYear := now.Year()
	lastYear := now.Year() - 1
	lastYearStr := strconv.Itoa(lastYear)
	currentYearStr := strconv.Itoa(currentYear)

	if filter == currentYearStr {
		Filter = "Current Year"
		check = "Current Year"
	}

	if filter == "365" {

		Filter = "Last Year(" + lastYearStr + "-" + currentYearStr + ")"
		check = "Last 1 Year"
	}

	t := template.Must(template.ParseFiles(
		"services/webbff/templates/graph/pgraphByService.html",
		"services/webbff/templates/head/head.html",
		"services/webbff/templates/head/nucleo.html",
		"services/webbff/templates/graph/import.html",
		"services/webbff/templates/graph/title.html",
		"services/webbff/templates/graph/asideService.html",
		"services/webbff/templates/graph/pgServicemain.html",
		"services/webbff/templates/footer/footer.html",
		"services/webbff/templates/script/serviceGraphScript.html",
	))
	t.Execute(response, struct {
		ProjectDetail                types.GetProjectByProjectIdResponse
		GraphData                    map[string][]types.GetGraphDataResponse
		UserId                       string
		ProjectId                    string
		UserProjects                 []map[string]string
		GetServiceCostMinMAxTotalAvg types.GetServiceCostMinMAxTotalAvgResponse
		ServiceName                  []string
		ProjectServiceList           []string
		SelectedServices             []string
		TestService                  []string
		GraphTag                     string
		Year                         string
		Check                        string
		ServiceUrl                   string
	}{ProjectDetail: resProjectDetail,
		UserId:                       uIdRequest,
		GraphData:                    resGraphData,
		ProjectId:                    pidRequest,
		UserProjects:                 resUserProjects,
		GetServiceCostMinMAxTotalAvg: resServiceCostDifferentAspects,
		ServiceName:                  servicesParam,
		ProjectServiceList:           resGetProjectServices,
		SelectedServices:             servicesParam,
		TestService:                  testService,
		GraphTag:                     Filter,
		Year:                         currentYearStr,
		Check:                        check,
		ServiceUrl:                   fmt.Sprintf("/%s/projects/%s/%s/export/service/%s", uIdRequest, pidRequest, filter, testService),
	})

}

func projectExportPageHandler(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	pidRequest := val["pid"]

	filter := val["filter"]

	filename, data, err := services.ProjectExport(pidRequest, filter, "AWS")
	if err != nil {
		fmt.Println("Error in project export", err)
	}

	response.Header().Set("Content-Type", "text/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+filename)
	response.Write(data)
}

func projectExportPageHandlerByService(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	pidRequest := val["pid"]
	serviceName := val["value"]

	filter := val["filter"]

	res := strings.Replace(serviceName, "[", "", -1)
	res1 := strings.Replace(res, "]", "", -1)
	res1 = res1[:len(res1)-1]
	newService := strings.Split(res1, "/ ")

	filename, data, err := services.ProjectExportByService(pidRequest, filter, newService)
	if err != nil {
		fmt.Println("Error in project export", err)
	}
	response.Header().Set("Content-Type", "text/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+filename)
	response.Write(data)

}

func projectExportPageHandlerByDate(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	pidRequest := val["pid"]
	fromDate := request.FormValue("from-date")
	todate := request.FormValue("to-date")
	filename, data, err := services.ProjectExportByDate(pidRequest, fromDate, todate, "AWS")
	if err != nil {
		fmt.Println("Error in project export", err)
	}

	response.Header().Set("Content-Type", "text/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+filename)
	response.Write(data)
}

func projectExportPageHandlerByYear(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	pidRequest := val["pid"]
	year := request.FormValue("year")

	fmt.Println(year)
	filename, data, err := services.ProjectExportByYear(pidRequest, year, "AWS")
	if err != nil {
		fmt.Println("Error in project export", err)
	}

	response.Header().Set("Content-Type", "text/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+filename)
	response.Write(data)
}

// update pinned project
func updatePinnedProject(response http.ResponseWriter, request *http.Request) {
	// check data from how to received.
	val := mux.Vars(request)
	pidRequest := val["pid"]
	ispinned := val["pin"]

	err := services.UpdatePinProjects(pidRequest, ispinned)
	if err != nil {
		fmt.Print(errors.New("update pin: error in update"))
	}

}

func DailyCostSumHandler(response http.ResponseWriter, request *http.Request) {
	val := mux.Vars(request)
	uidRequest := val["uid"]
	dayFilter := val["filter"]
	resDailyCostSum, err := services.DailyTotalCostOfAllProjectsOfUser(uidRequest, dayFilter)
	if err != nil {
		fmt.Print(errors.New("graphPageHandler: data not found"))
	}

	t := template.Must(template.ParseFiles("services/webbff/templates/pgraph.html"))
	t.Execute(response, struct {
		DailyCostSum []types.TotalCostSumResonse
		UserId       string
	}{
		DailyCostSum: resDailyCostSum,
		UserId:       uidRequest,
	})
}

// delete Project
func deleteProjectHandler(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	userId := val["uid"]
	projectId := val["pid"]
	err := services.DeleteProject(projectId)

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(response, request, fmt.Sprintf("/%s/projects", userId), http.StatusFound)

}

// update pinned project by Project Name
func updatePinnedProjectByProjectName(response http.ResponseWriter, request *http.Request) {
	// check data from how to received.
	val := mux.Vars(request)
	uidRequest := val["uid"]
	pName := val["pName"]
	ispinned := val["pin"]

	err := services.UpdatePinProjectsByProjectName(uidRequest, pName, ispinned)
	if err != nil {
		fmt.Print(errors.New("update pin: error in update"))
	}

}
func projectExportPageHandlerByServiceAndCustomDate(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/csv")

	val := mux.Vars(request)
	pidRequest := val["pid"]
	serviceName := val["value"]

	res := strings.Replace(serviceName, "[", "", -1)
	res1 := strings.Replace(res, "]", "", -1)
	res1 = res1[:len(res1)-1]
	newService := strings.Split(res1, "/ ")

	fromDate := request.FormValue("from-date")
	todate := request.FormValue("to-date")

	filename, data, err := services.ProjectExportPageHandlerByServiceAndDate(pidRequest, fromDate, todate, newService)
	if err != nil {
		fmt.Println("err in proj export")
	}
	response.Header().Set("Content-Type", "text/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+filename)
	response.Write(data)

}

func projectExportPageHandlerByServiceWithYear(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/csv")

	val := mux.Vars(request)
	pidRequest := val["pid"]
	serviceName := val["value"]

	fmt.Println(serviceName)
	res := strings.Replace(serviceName, "[", "", -1)
	res1 := strings.Replace(res, "]", "", -1)
	res1 = res1[:len(res1)-1]
	newService := strings.Split(res1, "/ ")

	year := request.FormValue("year")
	fmt.Println(newService)

	filename, data, err := services.ProjectExportPageHandlerByServiceWithYear(pidRequest, year, newService)
	if err != nil {
		fmt.Println("projectExportPageHandlerByServiceWithYear(): error is ", err)
	}

	response.Header().Set("Content-Type", "text/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+filename)
	response.Write(data)

}

func projectSettingHandler(response http.ResponseWriter, request *http.Request) {
	val := mux.Vars(request)
	userId := val["uid"]
	projectId := val["pid"]

	// Initialize error message variable
	var errorMessage string

	awsKeys, err := services.GetAWSCredentials(projectId)
	if err != nil {
		fmt.Println("Error getting AWS credentials: " + err.Error())
	}

	if !awsKeys {
		if request.Method == http.MethodPost {
			if request.FormValue("action") == "submit" {
				err := request.ParseForm()
				if err != nil {
					fmt.Println("Error parsing form value: " + err.Error())
				}

				// handle submit button click  | stay on the same page
				accesskey := request.FormValue("awsaccesskey")
				secretkey := request.FormValue("awssecretkey")

				err = services.CreateAwsCredentials(projectId, accesskey, secretkey)
				if err != nil {
					errorMessage = "Error creating AWS credentials: " + err.Error()
					fmt.Println(errorMessage)
				}
			}
		}

		awsKeys, _ = services.GetAWSCredentials(projectId)
	}

	if request.FormValue("action") == "save_changes" {
		// handle save changes button click | redirect to previous page
		err := request.ParseForm()
		if err != nil {
			fmt.Println("Error parsing form value: " + err.Error())
		}

		projectname := request.FormValue("projectname")
		projectdesc := request.FormValue("projectdesc")

		err = services.UpdateProjectDetails(projectId, projectname, projectdesc)
		if err != nil {
			fmt.Println("Error updating project details: " + err.Error())
		} else {
			redirectUrl := fmt.Sprintf("/%s/projects/%s/30", userId, projectId)
			http.Redirect(response, request, redirectUrl, http.StatusFound)
		}
	}

	// render your template here
	resp, err := services.GetProjectByProjectId(projectId, userId)
	if err != nil {
		fmt.Println("Error getting project details: " + err.Error())
	}

	t := template.Must(template.ParseFiles(
		"services/webbff/templates/settings/projectSetting.html",
		"services/webbff/templates/head/head.html",
		"services/webbff/templates/head/nucleo.html",
		"services/webbff/templates/settings/import.html",
		"services/webbff/templates/settings/title.html",
		"services/webbff/templates/settings/header.html",
		"services/webbff/templates/settings/psmain.html",
		"services/webbff/templates/footer/footer.html",
	))

	// add context variable to indicate if form has been submitted
	t.Execute(response, struct {
		ProjectDetails types.GetProjectByProjectIdResponse
		Id             string
		AwsKeys        bool
		ErrorMessage   string // Pass error message variable to the template
	}{
		ProjectDetails: resp,
		Id:             projectId,
		AwsKeys:        awsKeys,
		ErrorMessage:   errorMessage, // Set error message variable
	})
}

func updateActiveProject(_ http.ResponseWriter, request *http.Request) {
	val := mux.Vars(request)
	pidRequest := val["pid"]
	IsActive := val["active"]

	err := services.UpdateActiveProject(pidRequest, IsActive)
	if err != nil {
		fmt.Println(err)
	}
}

func deleteAwsCredentialsHandler(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	userId := val["uid"]
	projectId := val["pid"]
	err := services.DeleteAwsCredentials(projectId)

	if err != nil {
		log.Println("services: DeleteAwsCredentials()", err)
	}

	http.Redirect(response, request, fmt.Sprintf("/%s/projects/%s/30", userId, projectId), http.StatusFound)

}

func newProjectHandler(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	userId := val["uid"]

	if request.Method == "POST" {

		projectname := request.FormValue("projectname")
		projectdesc := request.FormValue("projectdesc")

		_, err := services.NewProjectService(userId, projectname, projectdesc)

		if err != nil {
			// Render an error message on the frontend
			t := template.Must(template.ParseFiles(
				"services/webbff/templates/newProject/project.html",
				"services/webbff/templates/head/head.html",
				"services/webbff/templates/head/nucleo.html",
				"services/webbff/templates/newProject/import.html",
				"services/webbff/templates/newProject/title.html",
				"services/webbff/templates/newProject/header.html",
				"services/webbff/templates/newProject/pmain.html",
				"services/webbff/templates/footer/footer.html",
			))

			t.Execute(response, struct {
				UserId string
				Error  string
			}{
				UserId: userId,
				Error:  err.Error(),
			})
			return
		}

		http.Redirect(response, request, fmt.Sprintf("/%s/projects/", userId), http.StatusFound)

	}

	t := template.Must(template.ParseFiles(
		"services/webbff/templates/newProject/project.html",
		"services/webbff/templates/head/head.html",
		"services/webbff/templates/head/nucleo.html",
		"services/webbff/templates/newProject/import.html",
		"services/webbff/templates/newProject/title.html",
		"services/webbff/templates/newProject/header.html",
		"services/webbff/templates/newProject/pmain.html",
		"services/webbff/templates/footer/footer.html",
	))

	t.Execute(response, struct {
		UserId string
		Error  string
	}{
		UserId: userId,
		Error:  "",
	})
}

// for graph pinned button
func updateByGraphPinned(response http.ResponseWriter, request *http.Request) {
	// check data from how to received.

	val := mux.Vars(request)
	pidRequest := val["pid"]
	isPinned := val["pin"]

	err := services.UpdateByGraphPinned(pidRequest, isPinned)
	if err != nil {
		fmt.Print(errors.New("update pin: error in update"))
	}

}

// mongo export csv
func mongoProjectExportPageHandler(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	pidRequest := val["pid"]

	filter := val["filter"]
	filename, data, err := services.ProjectExport(pidRequest, filter, "MongoDb")
	if err != nil {
		fmt.Println("Error in project export", err)
	}

	response.Header().Set("Content-Type", "text/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+filename)
	response.Write(data)
}

func mongoProjectExportPageHandlerByDate(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	pidRequest := val["pid"]
	fromDate := request.FormValue("from-date")
	todate := request.FormValue("to-date")
	filename, data, err := services.ProjectExportByDate(pidRequest, fromDate, todate, "MongoDb")
	if err != nil {
		fmt.Println("Error in project export", err)
	}

	response.Header().Set("Content-Type", "text/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+filename)
	response.Write(data)
}

func mongoProjectExportPageHandlerByYear(response http.ResponseWriter, request *http.Request) {

	val := mux.Vars(request)
	pidRequest := val["pid"]
	year := request.FormValue("year")

	fmt.Println(year)
	filename, data, err := services.ProjectExportByYear(pidRequest, year, "MongoDb")
	if err != nil {
		fmt.Println("Error in project export", err)
	}

	response.Header().Set("Content-Type", "text/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+filename)
	response.Write(data)
}

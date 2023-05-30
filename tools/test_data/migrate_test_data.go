package main

// postgres.CreateUser("", "", "")
//create user
// add one project for that user
// add usage data for that project

import (
	"crypto/sha256"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/postgres"
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/tealeg/xlsx"
)

const (
	userName     = "Vijay"
	userEmail    = "social@tftus.com"
	userPassword = "Tftus@1234"
)
const (
	projectName      = "Flipkart"
	projectInfo      = "Ecommerce site"
	projectSecretKey = "AKIA6RTLPWSTCQKARSX5"
	projectAccessKey = "0ZDcHgvdt3Cg1DTiU4M+5T5RJLcsapFathBTlDqr"
)
const (
	MICRODOLLAR = 1000000000
)

func HashData(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ExcelData() []string {
	excelFileName := "tools/test_data/bill.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		log.Fatal(err)
	}

	var values []string
	for _, sheet := range xlFile.Sheets {
		// Loop through the rows in the sheet
		for i, row := range sheet.Rows {
			// If it's the first row, assume it's the header row and skip it
			if i == 0 {
				continue
			}

			// Get the values from the row

			for _, cell := range row.Cells {
				value, _ := cell.FormattedValue()
				values = append(values, value)
			}

			if i == 967 {
				break
			}
		}
	}
	return values
}

// newUserProject: does create new User, new Project and mapUserToProject
func newUserProject() (int64, error) {
	// create user
	newUser, err := postgres.CreateUser(userName, userEmail, HashData(userPassword))
	if err != nil {
		return 0, err
	}

	// create project
	newProject, err := postgres.CreateProject(projectName, projectInfo)
	if err != nil {
		return 0, err
	}

	// create mapping user-project
	newUserProject, err := postgres.CreateUserProject(newUser.ID, newProject.ID)
	if err != nil {
		return 0, err
	}

	_, err = postgres.CreateAwsCredentials(newUserProject.ID, projectSecretKey, projectAccessKey)
	if err != nil {
		return 0, err
	}

	return newUserProject.ProjectId, nil
}

var digitCheck = regexp.MustCompile(`^[0-9]+`)

func testAwsUsageData() ([]*types.DbServiceCostAws, error) {

	xldata := ExcelData()
	var Result []*types.DbServiceCostAws
	now := time.Now()
	dateCounter := 0

	for i := 0; i < len(xldata)-3; i = i + 3 {

		// service name
		serviceName := xldata[i]
		if digitCheck.MatchString(serviceName) {
			fmt.Printf("%q looks like a number.\n", serviceName)
			continue
		}
		if serviceName == "" {
			fmt.Printf("%q looks like serviceName is nil.\n", serviceName)
			continue
		}

		// cost
		cost, _ := strconv.ParseFloat(xldata[i+1], 64)
		costMicro := int64(cost * MICRODOLLAR)

		// new date from current date
		dateCounter--
		dates := now.AddDate(0, 0, dateCounter)

		serviceCostAws := &types.DbServiceCostAws{
			ServiceTitle:           serviceName,
			UsageDate:              dates,
			AccruedCostMicrodollar: costMicro,
		}
		Result = append(Result, serviceCostAws)
	}
	return Result, nil
}

func main() {
	pid, _ := newUserProject()

	list, _ := testAwsUsageData()

	for _, v := range list {

		_, err := postgres.CreateServiceCostAws(pid, v.ServiceTitle, v.UsageDate, v.AccruedCostMicrodollar)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}

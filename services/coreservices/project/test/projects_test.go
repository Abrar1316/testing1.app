package test

import (
	"fmt"
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/coreservices/project"
	"github.com/TFTPL/AWS-Cost-Calculator/services/postgres"
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	userId    int64
	projectId int64
	filter    int64
)

// testuser details
var (
	testname     = "testUser"
	testemail    = "socialtest@tftus.com"
	testpassword = "Test@12345"
	testId       = "123456"
)

// testproject details
var (
	testprojectname       = "billingapp"
	testprojectdesc       = "generates bills"
	// testawsaccesskey      = "12345"
	// testawssecretkey      = "qwerty"
	testserviceTitle      = "aws"
	testcostInMicrodollar = 100000000000
)

var _ = Describe("ProjectService", func() {
	Describe("NewProjectService", func() {
		It("should return no error and create a project if all required fields are validated", func() {

			name := testname
			email := testemail
			password := testpassword
			user, err := postgres.CreateUser(name, email, password)
			if err != nil {
				fmt.Println("error is", err)
			}
			userId = user.ID

			userID := fmt.Sprintf("%d", userId)
			projectname := testprojectname
			projectdesc := testprojectdesc

			_, err = project.NewProjectService(userID, projectname, projectdesc)

			// err1 := project.CreateAwsCredentials(pid, testawsaccesskey, testawssecretkey)

			Expect(err).To(BeNil())
			// Expect(err1).To(BeNil())
		})
		It("should return an error if required fields are missing ", func() {
			userID := testId
			projectname := testprojectname
			projectdesc := testprojectdesc

			_, err := project.NewProjectService(userID, projectname, projectdesc)

			Expect(err).ToNot(BeNil())
		})
	})
	
	
	Describe("GetProjectByProjectId", func() {
		It("should return project info and no error for existing projectid", func() {
			projects, _ := postgres.GetProjectsByUserId(uint64(userId))
			projectId = projects[0].ProjectId

			pId := fmt.Sprintf("%d", projectId)
			uid := fmt.Sprintf("%d", userId)

			resp, err := project.GetProjectByProjectId(pId, uid)
			expectedresp := types.GetProjectByProjectIdResponse{
				UserId:      uid,
				Name:        testprojectname,
				Description: testprojectdesc,
			}

			Expect(resp).To(Equal(expectedresp))
			Expect(err).To(BeNil())
		})
		It("should return an error if project id not exists ", func() {

			_, err := project.GetProjectByProjectId(testId, "")
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("GetListOfProjectByUserId", func() {
		It("should return project info and no error for existing projectid", func() {

			uid := fmt.Sprintf("%d", userId)
			projectName := []string{"projectName"}

			_, err := project.GetListOfProjectByUserId(uid, projectName)
			Expect(err).To(BeNil())
		})
		// It("should return empty response non existing userId", func() {
		// 	var expectedresp []map[string]string
		// 	projectName := []string{"projectName"}

		// 	resp, err := project.GetListOfProjectByUserId(testId, projectName)

		// 	Expect(resp).To(Equal(expectedresp))
		// 	Expect(err).ToNot(BeNil())
		// })
	})
	Describe("GetServiceCostMinMAxTotalAvgByServiceName", func() {
		It("should return required costs for given project and its service", func() {

			// Need to insert data into AWS first
			pid := projectId
			serviceTitle := testserviceTitle
			usageDate := time.Now().AddDate(0, 0, -1)
			costInMicrodollar := testcostInMicrodollar

			_, err := postgres.CreateServiceCostAws(pid, serviceTitle, usageDate, int64(costInMicrodollar))
			if err != nil {
				fmt.Println("error in inserting into CreateServiceCostAws ", err)
			}

			expectedresp := types.GetServiceCostMinMAxTotalAvgResponse{
				TotalCost:   "100.00 $",
				AverageCost: "3.33 $",
				MinimumCost: "100.00 $",
				MaximumCost: "100.00 $",
			}

			projid := fmt.Sprintf("%d", pid)
			resp, err := project.GetServiceCostMinMAxTotalAvgByServiceName(projid, serviceTitle)

			Expect(resp).To(Equal(expectedresp))
			Expect(err).To(BeNil())
		})
		It("should return all costs as 0$ for a project and its service if it doesnot exists", func() {
			expectedresp := types.GetServiceCostMinMAxTotalAvgResponse{
				TotalCost:   "0.00 $",
				AverageCost: "0.00 $",
				MinimumCost: "0.00 $",
				MaximumCost: "0.00 $",
			}

			resp, err := project.GetServiceCostMinMAxTotalAvgByServiceName(testId, "")
			fmt.Println("Respinse is ", resp)

			Expect(resp).To(Equal(expectedresp))
			Expect(err).To(BeNil())
		})
	})
	Describe("GetServiceCostMinMAxTotalAvg", func() {
		It("should return required costs for given project", func() {

			expectedresp := types.GetServiceCostMinMAxTotalAvgResponse{
				TotalCost:   "100.00 $",
				AverageCost: "3.33 $",
				MinimumCost: "100.00 $",
				MaximumCost: "100.00 $",
			}

			projid := fmt.Sprintf("%d", projectId)
			resp, err := project.GetServiceCostMinMAxTotalAvg(projid, "30")

			Expect(resp).To(Equal(expectedresp))
			Expect(err).To(BeNil())
		})
		It("should return all costs as 0$ for a project if it doesnot exists", func() {
			expectedresp := types.GetServiceCostMinMAxTotalAvgResponse{
				TotalCost:   "0.00 $",
				AverageCost: "0.00 $",
				MinimumCost: "0.00 $",
				MaximumCost: "0.00 $",
			}

			resp, err := project.GetServiceCostMinMAxTotalAvg(testId, "30")

			Expect(resp).To(Equal(expectedresp))
			Expect(err).To(BeNil())
		})
	})
	Describe("GetGraphData", func() {
		It("should return no error and return graph data", func() {

			projid := fmt.Sprintf("%d", projectId)
			_, err := project.GetGraphData(projid, "30")

			Expect(err).To(BeNil())
		})
	})
	Describe("GetGraphDataByServiceName", func() {
		It("should return no error and return graph data by service name", func() {

			projid := fmt.Sprintf("%d", projectId)
			servicename := []string{testserviceTitle}
			_, err := project.GetGraphDataByServiceName(projid, servicename, "30")

			Expect(err).To(BeNil())
		})
	})
	Describe("DailyTotalCostOfAllProjectsOfUser", func() {
		It("should return no error and return graph data", func() {

			uid := fmt.Sprintf("%d", userId)
			daysFilter := fmt.Sprintf("%d", filter)

			_, err := project.DailyTotalCostOfAllProjectsOfUser(uid, daysFilter)

			Expect(err).To(BeNil())
		})
		It("should return error for invalid userid", func() {

			_, err := project.DailyTotalCostOfAllProjectsOfUser("uid", "daysFilter")
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("GetDistinctServiceTitlesForProject", func() {
		It("should return no error and return service title", func() {

			expectedresp := []string{"aws", "Combined Cost"}

			pid := fmt.Sprintf("%d", projectId)
			resp, err := project.GetDistinctServiceTitlesForProject(pid)

			Expect(resp).To(Equal(expectedresp))
			Expect(err).To(BeNil())
		})
		It("should return nil response for invalid projectid ", func() {

			expectedresp := []string{"Combined Cost"}

			resp, err := project.GetDistinctServiceTitlesForProject(testId)
			Expect(resp).To(Equal(expectedresp))
			Expect(err).To(BeNil())
		})
	})
	Describe("ProjectExport", func() {
		It("should return no error and return file name", func() {

			expectedresp := fmt.Sprintf("%s-%s-Cost-Last-%s-days.csv", testprojectname, "AWS", "30")

			pid := fmt.Sprintf("%d", projectId)
			resp, _, err := project.ProjectExport(pid, "30", "AWS")

			Expect(resp).To(Equal(expectedresp))
			Expect(err).To(BeNil())
		})
	})
	Describe("ProjectExportByService", func() {
		It("should return no error and return file name", func() {

		
			service := []string{testserviceTitle}
			expectedresp := fmt.Sprintf("%s-Services-Aws-Cost-Last-%s-days.csv", testprojectname, "30")

			pid := fmt.Sprintf("%d", projectId)
			resp, _, err := project.ProjectExportByService(pid, "30", service)

			Expect(resp).To(Equal(expectedresp))
			Expect(err).To(BeNil())
		})
	})
	Describe("UpdateProjectsDetailsByProjectId", func() {
		// It("should return no error and updates the project name", func() {

		// 	pid := fmt.Sprintf("%d", projectId)
		// 	updatedname := "updated billing app"
		// 	err := project.UpdateProjectsDetailsByProjectId(pid, updatedname, "","", "")

		// 	Expect(err).To(BeNil())
		// })
		It("should return error and if project not exists", func() {

			err := project.UpdateProjectsDetailsByProjectId("1234567", "", "", "", "")
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("UpdatePinProjects", func() {
		It("should return no error and updates the pin status of the project ", func() {

			pid := fmt.Sprintf("%d", projectId)
			status := "Unpinned"
			err := project.UpdatePinProjects(pid, status)

			Expect(err).To(BeNil())
		})
	})
	Describe("GetAllActivePinnedProjects", func() {
		It("should return no error and return all pinned projects", func() {

			uid := fmt.Sprintf("%d", userId)

			_, err := project.GetAllActivePinnedProjects(uid)
			Expect(err).To(BeNil())
		})
		It("should return error and if project not exists", func() {

			_, err := project.GetAllActivePinnedProjects(testId)
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("GetCostOfProjectsForDashboard", func() {
		It("should return no error and return total cost of all the projects", func() {

			uid := fmt.Sprintf("%d", userId)
			filter := "30"
			projNames := []string{testprojectname}

			_, err := project.GetCostOfProjectsForDashboard(uid, filter, projNames)

			Expect(err).To(BeNil())

		})

	})
	Describe("UpdateProjectDetails", func() {
		It("should update project details with new name and description", func() {
			pid := fmt.Sprintf("%d", projectId)
			updatedName := "Updated Test Project"
			updatedDescription := "Updated Description"
			err := project.UpdateProjectDetails(pid, updatedName, updatedDescription)
			Expect(err).To(BeNil())
		})
		It("should update project name only if description is not provided", func() {
			pid := fmt.Sprintf("%d", projectId)
			updatedName := "Updated Test Project Name Only"
			err := project.UpdateProjectDetails(pid, updatedName, "")
			Expect(err).To(BeNil())
		})
		It("should update project description only if name is not provided", func() {
			pid := fmt.Sprintf("%d", projectId)
			updatedDescription := "Updated Description Only"
			err := project.UpdateProjectDetails(pid, "", updatedDescription)
			Expect(err).To(BeNil())
		})
		It("should return an error when project does not exist", func() {
			err := project.UpdateProjectDetails("123456", "Updated Test Project", "Updated Description")
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("DeleteProject", func() {
		It("should return error  if project not exists", func() {

			err := project.DeleteProject(testId)
			Expect(err).ToNot(BeNil())
		})
		It("should return no error  and deletes the project", func() {
			pid := fmt.Sprintf("%d", projectId)

			err := project.DeleteProject(pid)
			Expect(err).To(BeNil())
		})
	})
	Describe("Deleteuser", func() {
		It("Delete the test user that has been created for tests", func() {

			user, _ := postgres.GetUserByEmail(testemail)
			_, err := postgres.DeleteUser(user.ID)
			if err != nil {
				fmt.Println("error in deleting the test user", err)
			}
			Expect(err).To(BeNil())
		})
	})
})

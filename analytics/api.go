package main

import (
	"log"
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/postgres"
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/go-co-op/gocron"
)

func getAnalytics() {

	dailyCostData := []types.DbAnalytics{}

	users, err := postgres.GetAllUsers()
	if err != nil {
		log.Println("Error in getting all user details", err)
		return
	}

	for uID := range users {
		projects, err := postgres.GetProjectDetailsByUserId(users[uID].ID)
		if err != nil {
			log.Println("Error in getting project details of user:", users[uID].ID, "Error is :", err)
		}
		for i := range projects {
			cost, err := postgres.GetDailyCostOfProject(projects[i].ID)
			if err != nil {
				log.Println("Error in getting daily cost of project", projects[i].ID, "Error is :", err)
			}

			dailydata := types.DbAnalytics{
				ID:                   projects[i].ID,
				Name:                 projects[i].Name,
				Email:                users[uID].Email,
				DailyCostMicrodollar: cost.AccruedCostMicrodollar,
			}
			dailyCostData = append(dailyCostData, dailydata)
		}
	}

	err = postgres.UpsertAnalyticsData(dailyCostData)
	if err != nil {
		log.Println("Error in upserting data into analytics", err)
		return
	}

}

func main() {
	s := gocron.NewScheduler(time.Local)

	s.Every(1).Day().At("02:00").Do(func() {

		log.Println("Inserting daily cost of projects into analytics table ", time.Now())
		getAnalytics()
	})

	s.StartBlocking()
}

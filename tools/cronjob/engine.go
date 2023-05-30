package main

import (
	"sync"
	"time"

	aws "github.com/TFTPL/AWS-Cost-Calculator/tools/cronjob/cronjob_aws"
	mongo "github.com/TFTPL/AWS-Cost-Calculator/tools/cronjob/cronjob_mongoCloud"
	"github.com/go-co-op/gocron"
)

func main() {
	var wg sync.WaitGroup
	crontime := "19:11"

	s := gocron.NewScheduler(time.Local)

	// Schedule the Mongo cron job
	s.Every(1).Day().At(crontime).Do(func() {
		wg.Add(1)
		go mongo.FetchBillCronJobMongo(&wg)
	})

	// Schedule the AWS cron job
	s.Every(1).Day().At(crontime).Do(func() {
		wg.Add(1)
		go aws.FetchBillCronJobAWS(&wg)
	})

	s.StartBlocking()

	// Wait for both cron jobs to complete
	wg.Wait()
}

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"sync"
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/postgres"
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/go-co-op/gocron"
)

const (
	// needs to be modified
	from = "dummytestbillingapp@gmail.com"
	pass = "obtujhsyfdazrdzv"

	// smtp server configuration.
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"

	divider = 1000000000
)

func sendEmail(wg *sync.WaitGroup, data types.DbAnalytics) {

	defer wg.Done()

	to := []string{
		data.Email,
	}

	// Authentication.
	auth := smtp.PlainAuth("", from, pass, smtpHost)

	t, err := template.ParseFiles("emailalert/emailtemplate.html")
	if err != nil {
		log.Println("sendEmail: Error in parsing Email tamplate ", err)
		return
	}

	// Get the yesterday date
	date := time.Now().AddDate(0, 0, -1)

	generatedcost := float64(data.DailyCostMicrodollar / divider)
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Daily limit exceeded for project %s \n%s\n\n", data.Name, mimeHeaders)))

	t.Execute(&body, struct {
		Name string
		Date string
		Cost float64
	}{
		Name: data.Name,
		Date: date.Format("02-01-2006"),
		Cost: generatedcost,
	})

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		log.Printf("sendEmail: Error in sending Email of Project %s: Error is %v", data.Name, err)
		return
	}

	log.Println("Email sent to :", data.Email, "for project:", data.Name)
}

func getAnalyticsReport() {

	report, err := postgres.ReadAnalyticsData()
	if err != nil {
		log.Println("getAnalyticsReport: Error in reading analytics data :", err)
		return
	}

	if len(report) == 0 {
		log.Println("No project exceeds its daily limit")
		return
	}

	wg := new(sync.WaitGroup)

	for _, val := range report {
		wg.Add(1)
		go sendEmail(wg, val)
	}

	wg.Wait()

}

func main() {

	s := gocron.NewScheduler(time.Local)

	s.Every(1).Day().At("03:00").Do(func() {

		log.Println("Reading daily cost of projects from analytics table ", time.Now())
		getAnalyticsReport()
	})

	s.StartBlocking()

}

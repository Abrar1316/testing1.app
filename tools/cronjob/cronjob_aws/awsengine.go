package aws

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/postgres"
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/jackc/pgx/v5"
)

const (
	BLENDED_COST_CAPS      = "BLENDED_COST"
	BLENDED_COST_CAMELCASE = "BlendedCost"
	UNBLENDED_COST         = "UnblendedCost"
	DAILY                  = "DAILY"
	GROUP_DIMENSION        = "DIMENSION"
	GROUP_SERVICE          = "SERVICE"
	AWS_REGION             = "us-west-2"
	DATE_FORMAT            = "2006-01-02"
	TEST_ACCESS_KEY        = "AKIA6RTLPWSTCQKARSX5"
	TEST_SECRET_KEY        = "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"
	MICRODOLLAR            = 1000000000
	TAX                    = "TAX"
)

func FetchBillCronJobAWS(wg *sync.WaitGroup) {
	defer wg.Done()

	// Your logic for fetching AWS service data
	log.Println("Fetching AWS service data at time ", time.Now())

retry:
	err := fetchBillAWS()
	if err != nil && err == pgx.ErrNoRows {
		log.Println("Error in cron job, retrying...", err)
		time.Sleep(1 * time.Hour)
		goto retry
	}
}

func processAwsData(awsData *costexplorer.GetCostAndUsageOutput) ([]*types.DbServiceCostAws, error) {
	result := make([]*types.DbServiceCostAws, 0)

	for _, results := range awsData.ResultsByTime {
		startDate := *results.TimePeriod.Start
		for _, groups := range results.Groups {
			for _, metrics := range groups.Metrics {
				costFloat, _ := strconv.ParseFloat(*metrics.Amount, 64)
				dateTime, _ := time.Parse(DATE_FORMAT, startDate)
				x := &types.DbServiceCostAws{
					ServiceTitle:           *groups.Keys[0],
					UsageDate:              dateTime,
					AccruedCostMicrodollar: int64(costFloat * MICRODOLLAR),
				}
				result = append(result, x)
			}
		}
	}
	return result, nil
}

func fetchBillAWS() error {
	projectList, err := postgres.GetProjectsByIsActive(true)
	if err != nil {
		return err
	}
	for _, v := range projectList {
		awskeys, err := postgres.ReadAwsCredentials(v.ID)
		if err != nil {
			return err
		}
		syncProjectAwsUsage(awskeys.AccessKey, awskeys.AwsSecretKey, awskeys.ProjectId)
	}
	err = postgres.DeleteServiceCostAwsbyServiceName(TAX)
	if err != nil {
		return err
	}
	return nil
}

func syncProjectAwsUsage(accessKey, secretKey string, projectId int64) error {
	awsData, err := getCost(accessKey, secretKey)
	if err != nil {
		return err
	}
	processedData, err := processAwsData(awsData)
	if err != nil {
		return err
	}

	for _, v := range processedData {
		_, err := postgres.CreateServiceCostAws(projectId, v.ServiceTitle, v.UsageDate, v.AccruedCostMicrodollar)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCost(ACCESS_KEY, SECRET_KEY string) (*costexplorer.GetCostAndUsageOutput, error) {
	awsSession := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(ACCESS_KEY, SECRET_KEY, ""),
		Region:      aws.String(AWS_REGION),
	}))

	svc := costexplorer.New(awsSession)
	// to get, -1 for last day
	now := getDates(0, 0, -1)
	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(*now.Start),
			End:   aws.String(*now.End),
		},
		Granularity: aws.String(DAILY),
		Metrics:     []*string{aws.String(BLENDED_COST_CAPS)},
		GroupBy: []*costexplorer.GroupDefinition{
			{
				Type: aws.String(GROUP_DIMENSION),
				Key:  aws.String(GROUP_SERVICE),
			},
		},
	}

	result, err := svc.GetCostAndUsage(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// getDates format (year, month, day), -1 for yesterday, -2 for last two days
func getDates(year, month, day int) *costexplorer.DateInterval {
	now := time.Now()
	then := now.AddDate(year, month, day)
	dateRange := costexplorer.DateInterval{}
	dateRange.SetEnd(now.Format(DATE_FORMAT))
	dateRange.SetStart(then.Format(DATE_FORMAT))
	return &dateRange
}

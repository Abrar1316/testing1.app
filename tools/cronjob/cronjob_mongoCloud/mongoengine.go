package mongo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	auth "github.com/Snawoot/go-http-digest-auth-client"
	"github.com/TFTPL/AWS-Cost-Calculator/services/postgres"
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/jackc/pgx"
)

var date = time.Now().AddDate(0, 0, -1).Format("2006-01-02T00:00:00Z")

type LineItem struct {
	GroupName        string  `json:"groupName"`
	UnitPriceDollars float64 `json:"unitPriceDollars"`
	Unit             string  `json:"unit"`
	TotalPriceCents  float64 `json:"totalPriceCents"`
	SKU              string  `json:"sku"`
	Quantity         float64 `json:"quantity"`
	EndDate          string  `json:"endDate"`
}

func FetchBillCronJobMongo(wg *sync.WaitGroup) {
	defer wg.Done()

	// Your logic for fetching Mongo service data
	log.Println("Fetching Mongo service data at time ", time.Now())

retry:
	err := fetchBillMongo()
	if err != nil && err == pgx.ErrNoRows {
		log.Println("Error in cron job, retrying...", err)
		time.Sleep(1 * time.Hour)
		goto retry
	}
}

func fetchBillMongo() error {

	listmongokeys, err := postgres.ReadAllMongoCredentials()

	if err != nil {
		log.Println(err)
		return err
	}
	for _, mongokeys := range listmongokeys {
		projects, err := postgres.ReadProject(mongokeys.ProjectId)
		err = syncMongoUsage(projects.Name, mongokeys.MongoOrganizationId, mongokeys.PrivateKey, mongokeys.PublicKey, mongokeys.ProjectId)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func syncMongoUsage(ProjectName, OrganizationId, PrivateKey, PublicKey string, projectId int64) error {
	// fmt.Println(OrganizationId, PrivateKey, PublicKey)
	mongoData, err := getMongoUsage(OrganizationId, PrivateKey, PublicKey)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	processMongoData, err := processMongoData(ProjectName, mongoData)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	for _, v := range processMongoData {
		_, err := postgres.CreateServiceCostMongo(projectId, v.MongoProjectName, v.ServiceTitle,
			v.AccruedCostMicrodollar, v.UsageDate, v.Unit, v.UnitPrice, v.Quantity)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil

}
func getMongoUsage(OrganizationId, PrivateKey, PublicKey string) (map[string]interface{}, error) {

	url := fmt.Sprintf("https://cloud.mongodb.com/api/atlas/v1.0/orgs/%s/invoices/pending/?granularity=DAILY", OrganizationId)

	client := &http.Client{
		Transport: auth.NewDigestTransport(PublicKey, PrivateKey, http.DefaultTransport),
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	return result, nil
}

func processMongoData(projectName string, result map[string]interface{}) ([]*types.DbServiceCostMongo, error) {
	finalresult := make([]*types.DbServiceCostMongo, 0)

	lineItemsValue, ok := result["lineItems"]
	if !ok || lineItemsValue == nil {
		// handle the case where the "lineItems" field is nil or not found
		return finalresult, nil
	}

	lineItems := result["lineItems"].([]interface{})
	// fmt.Println(lineItems)
	for _, item := range lineItems {
		itemMap := item.(map[string]interface{})
		if itemMap["endDate"].(string) == date && itemMap["groupName"].(string) == projectName {
			// fmt.Printf("groupName: %s, unitPriceDollars: %v, unit: %s, totalPriceCents: %v, sku: %s, quantity: %v, date : %s\n",
			// 	itemMap["groupName"], itemMap["unitPriceDollars"], itemMap["unit"],
			// 	itemMap["totalPriceCents"], itemMap["sku"], itemMap["quantity"], itemMap["endDate"])
			fmt.Println(itemMap["groupName"].(string),projectName)

			dateTime, _ := time.Parse("2006-01-02T00:00:00Z", date)
			x := &types.DbServiceCostMongo{
				MongoProjectName:       itemMap["groupName"].(string),
				ServiceTitle:           itemMap["sku"].(string),
				AccruedCostMicrodollar: int64(itemMap["totalPriceCents"].(float64)),
				UsageDate:              dateTime,
				Unit:                   itemMap["unit"].(string),
				UnitPrice:              int64(itemMap["unitPriceDollars"].(float64)),
				Quantity:               int64(itemMap["quantity"].(float64)),
			}
			finalresult = append(finalresult, x)
		}
	}
	return finalresult, nil
}

//depricated
// func fetchBillCronJob() error {
// 	projectList, err := postgres.GetProjectsByIsActive(true)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(&projectList)

// 		for _, v := range projectList {
// 			mongokeys, err := postgres.ReadMongoCredentials(v.ID)
// 			// if err == pgx.ErrNoRows {
// 			// 	i++
// 			// }
// 			if err != nil  {
// 				return err
// 			}
// 			fmt.Println(mongokeys)
// 			syncMongoUsage(mongokeys.MongoOrganizationId, mongokeys.PublicKey, mongokeys.PrivateKey, mongokeys.ProjectId)
// 		}
// 		return nil
// 	}

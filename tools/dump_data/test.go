package main

import (
	"fmt"

	"github.com/TFTPL/AWS-Cost-Calculator/services/postgres"
)

func main() {
	// date := time.Now().Format("2006-01-02")
	// postgres.CreateServiceCostAws(2, "ec2", date, 10000)
	// postgres.CreateServiceCostAws(2, "ec2", date, 12000)
	// postgres.CreateServiceCostAws(2, "ec2", date, 13500)
	list, _ := postgres.GetServiceCostByProject(2, "")
	for _, v := range list {
		fmt.Println(v)
	}
	// postgres.UpdateProject(2, "updated val", "description", "acces", "secret", true)
}

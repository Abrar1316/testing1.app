package postgres

import (
	"log"

	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
)

const (
	dailylimit = 10000000000
)

func UpsertAnalyticsData(data []types.DbAnalytics) error {

	_, err := GetDB().Model(&data).
		OnConflict("(id) DO UPDATE").
		Set("daily_cost_microdollar = EXCLUDED.daily_cost_microdollar , email = EXCLUDED.email").
		Insert()

	if err != nil {
		log.Println("Error in upserting records to analytics", err)
		return err
	}

	return nil
}

func ReadAnalyticsData() ([]types.DbAnalytics, error) {

	projects := []types.DbAnalytics{}

	err := GetDB().Model(&projects).
		Where("daily_cost_microdollar > ?", dailylimit).
		Select()

	if err != nil {
		log.Println("Error in getting project datails that exceeds daily limit", err)
		return nil, err
	}

	return projects, nil

}

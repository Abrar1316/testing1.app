package postgres

import (
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

// CreateServiceCostMongo create a ServiceCost
func CreateServiceCostMongo(projectId int64, mongoProjectName, serviceTitle string, costInMicrodollar int64, usageDate time.Time, unit string, unitPrice int64, quantity int64) (*types.DbServiceCostMongo, error) {
	serviceCost := &types.DbServiceCostMongo{
		ProjectId:              projectId,
		MongoProjectName:       mongoProjectName,
		ServiceTitle:           serviceTitle,
		AccruedCostMicrodollar: costInMicrodollar,
		UsageDate:              usageDate,
		Unit:                   unit,
		UnitPrice:              unitPrice,
		Quantity:               quantity,
	}

	err := Insert(serviceCost)

	if err != nil {
		return nil, errors.Wrap(err, "pgServiceCostMongo:CreateServiceCostMongo failed - could not insert row")
	}

	return serviceCost, nil
}

// ReadServiceCostMongo reads a ServiceCost
func ReadServiceCostMongo(id int64) (*types.DbServiceCostMongo, error) {
	serviceCost := &types.DbServiceCostMongo{
		ID: id,
	}

	err := getDB().Model(serviceCost).Where("id = ?", id).Select()

	if err != nil {
		return nil, errors.Wrap(err, "pgServiceCostMongo:ReadServiceCostMongo failed - could not query row")
	}

	return serviceCost, nil
}

// UpdateServiceCostMongo updates a ServiceCost
func UpdateServiceCostMongo(id, projectId int64, mongoProjectName, serviceTitle string, costInMicrodollar int64, usageDate time.Time) (*types.DbServiceCostMongo, error) {
	serviceCost := &types.DbServiceCostMongo{
		ID:                     id,
		ProjectId:              projectId,
		MongoProjectName:       mongoProjectName,
		ServiceTitle:           serviceTitle,
		AccruedCostMicrodollar: costInMicrodollar,
		UsageDate:              usageDate,
	}

	err := Update(serviceCost)

	if err != nil {
		return nil, errors.Wrap(err, "pgServiceCostMongo:UpdateServiceCostMongo failed - could not update row")
	}

	return serviceCost, nil
}

// DeleteServiceCostAws deletes a ServiceCost
func DeleteServiceCostMongo(id int64) (*types.DbServiceCostMongo, error) {
	serviceCost := &types.DbServiceCostMongo{
		ID: id,
	}

	err := Delete(serviceCost)

	if err != nil {
		return nil, errors.Wrap(err, "pgServiceCostMongo:DeleteServiceCostMongo failed - could not delete row")
	}

	return serviceCost, nil
}

func GetMongoServiceCostByProject(pid int64, filter string) ([]*types.CostSum, error) {
	res := make([]*types.CostSum, 0)

	query := GetDB().
		Model(&types.DbServiceCostMongo{}).
		ColumnExpr("sum(accrued_cost_microdollar) as accrued_cost_microdollar").
		Column("usage_date").
		Where("project_id = ? AND usage_date > current_date - interval ? day", pid, filter).
		Group("usage_date").
		Order("usage_date DESC")

	err := query.Select(&res)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return res, nil
}

func GetMongoServiceCostDateByProjectId(pid, interval int64) ([]*types.DateServiceCost, error) {
	res := make([]*types.DateServiceCost, 0)

	query := GetDB().
		Model(&types.DbServiceCostMongo{}).
		Column("usage_date", "service_title", "accrued_cost_microdollar").
		Where("project_id = ? AND usage_date > current_date - interval '? day'", pid, interval).
		// Group("usage_date").
		Order("usage_date ASC")

	err := query.Select(&res)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return res, nil
}

func GetMongoServiceCostDateByProjectIdAndDate(pid int64, fromDate, toDate string) ([]*types.DateServiceCost, error) {
	res := make([]*types.DateServiceCost, 0)

	query := GetDB().
		Model(&types.DbServiceCostMongo{}).
		Column("usage_date", "service_title", "accrued_cost_microdollar").
		Where("project_id = ? AND usage_date >= ? AND usage_date <= ?", pid, fromDate, toDate).
		Order("usage_date ASC")

	err := query.Select(&res)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	// fmt.Println(res)
	return res, nil
}

// search by year
func GetMongoServiceCostDateByProjectIdAndYear(pid int64, year string) ([]*types.DateServiceCost, error) {
	res := make([]*types.DateServiceCost, 0)

	query := GetDB().
		Model(&types.DbServiceCostMongo{}).
		Column("usage_date", "service_title", "accrued_cost_microdollar").
		Where("project_id = ? AND EXTRACT(year FROM usage_date) = ?", pid, year).
		Order("usage_date ASC")

	err := query.Select(&res)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return res, nil
}

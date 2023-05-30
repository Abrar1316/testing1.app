package postgres

import (
	"fmt"
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

// CreateServiceCostAws create a ServiceCost
func CreateServiceCostAws(projectId int64, serviceTitle string, usageDate time.Time, costInMicrodollar int64) (*types.DbServiceCostAws, error) {
	serviceCost := &types.DbServiceCostAws{
		ProjectId:              projectId,
		ServiceTitle:           serviceTitle,
		UsageDate:              usageDate,
		AccruedCostMicrodollar: costInMicrodollar,
	}

	err := Insert(serviceCost)

	if err != nil {
		return nil, errors.Wrap(err, "pgServiceCostAws:CreateServiceCostAws failed - could not insert row")
	}

	return serviceCost, nil
}

// ReadServiceCostAws reads a ServiceCost
func ReadServiceCostAws(id int64) (*types.DbServiceCostAws, error) {
	serviceCost := &types.DbServiceCostAws{
		ID: id,
	}

	err := getDB().Model(serviceCost).Where("id = ?", id).Select()

	if err != nil {
		return nil, errors.Wrap(err, "pgServiceCostAws:ReadServiceCostAws failed - could not query row")
	}

	return serviceCost, nil
}

// UpdateServiceCostAws updates a ServiceCost
func UpdateServiceCostAws(id int64, projectId int64, serviceTitle string, usageDate time.Time, costInMicrodollar int64) (*types.DbServiceCostAws, error) {
	serviceCost := &types.DbServiceCostAws{
		ID:                     id,
		ProjectId:              projectId,
		ServiceTitle:           serviceTitle,
		UsageDate:              usageDate,
		AccruedCostMicrodollar: costInMicrodollar,
	}
	err := Update(serviceCost)

	if err != nil {
		return nil, errors.Wrap(err, "pgServiceCostAws:UpdateServiceCostAws failed - could not update row")
	}

	return serviceCost, nil
}

// DeleteServiceCostAws deletes a ServiceCost
func DeleteServiceCostAws(id int64) (*types.DbServiceCostAws, error) {
	serviceCost := &types.DbServiceCostAws{
		ID: id,
	}

	err := Delete(serviceCost)

	if err != nil {
		return nil, errors.Wrap(err, "pgServiceCostAws:DeleteServiceCostAws failed - could not delete row")
	}

	return serviceCost, nil
}

func GetDateCostByService(pid int64) ([]types.DbServiceCostAws, error) {
	var result []types.DbServiceCostAws

	err := GetDB().Model(&result).Column("usage_date", "accrued_cost_microdollar").Where("project_id=?", pid).Select()
	if err != nil {
		return nil, err
	}

	return result, err
}

func GetServiceCostByProject(pid int64, filter string) ([]*types.CostSum, error) {
	res := make([]*types.CostSum, 0)

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
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

func GetServiceCostByProjectByServiceName(pid int64, serviceName string, filter string) ([]*types.CostSum, error) {
	res := make([]*types.CostSum, 0)

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
		ColumnExpr("sum(accrued_cost_microdollar) as accrued_cost_microdollar").
		Column("usage_date").
		Where("project_id = ? AND usage_date > current_date - interval ? day", pid, filter).Where("service_title=?", serviceName).
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

func GetServiceCostMinMAxTotalAvg(pid int64, filter string) (types.GetServiceCostMinMAxTotalAvg, error) {
	var res types.GetServiceCostMinMAxTotalAvg

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
		ColumnExpr("sum(accrued_cost_microdollar) as accrued_cost_microdollar").
		ColumnExpr("min(accrued_cost_microdollar) as min_accrued_cost_microdollar").
		ColumnExpr("count(accrued_cost_microdollar) as count").
		ColumnExpr("max(accrued_cost_microdollar) as max_accrued_cost_microdollar").
		Where("project_id = ? AND usage_date > current_date - interval ? day", pid, filter)

	err := query.Select(&res)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return types.GetServiceCostMinMAxTotalAvg{}, nil
		} else {
			return types.GetServiceCostMinMAxTotalAvg{}, err
		}
	}

	return res, nil
}

func GetServiceCostMinMAxTotalAvgByServiceName(pid int64, serviceName string) (types.GetServiceCostMinMAxTotalAvg, error) {
	var res types.GetServiceCostMinMAxTotalAvg

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
		ColumnExpr("sum(accrued_cost_microdollar) as accrued_cost_microdollar").
		ColumnExpr("min(accrued_cost_microdollar) as min_accrued_cost_microdollar").
		ColumnExpr("count(accrued_cost_microdollar) as count").
		ColumnExpr("max(accrued_cost_microdollar) as max_accrued_cost_microdollar").
		Where("project_id = ?", pid).
		Where("service_title = ?", serviceName).
		Where("usage_date > current_date - interval '30 day'")

	err := query.Select(&res)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return types.GetServiceCostMinMAxTotalAvg{}, nil
		} else {
			return types.GetServiceCostMinMAxTotalAvg{}, err
		}
	}

	return res, nil
}

// DeleteServiceCostAwsbyServiceName deletes a ServiceCost
func DeleteServiceCostAwsbyServiceName(service string) error {
	var services []types.DbServiceCostAws
	_, err := GetDB().Model(&services).Where("service_title = ?", service).Delete()

	if err != nil {
		return fmt.Errorf("invalid email or password")
	}

	return nil
}

func GetServiceCostDateByProjectId(pid, interval int64) ([]*types.DateServiceCost, error) {
	res := make([]*types.DateServiceCost, 0)

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
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

func GetCostDateByProjectIdAndService(pid int64, interval int64, services []string) ([]*types.DateServiceCost, error) {
	res := make([]*types.DateServiceCost, 0)

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
		Column("usage_date", "service_title", "accrued_cost_microdollar").
		Where("project_id = ? AND service_title = ANY(?) AND usage_date > current_date - interval '? day'", pid, pg.Array(services), interval).
		Order("service_title ASC").Order("usage_date ASC")

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
func GetServiceCostDateByProjectIdAndDate(pid int64, fromDate, toDate string) ([]*types.DateServiceCost, error) {
	res := make([]*types.DateServiceCost, 0)

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
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
func GetServiceCostDateByProjectIdAndYear(pid int64, year string) ([]*types.DateServiceCost, error) {
	res := make([]*types.DateServiceCost, 0)

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
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

func GetDistinctServiceTitlesForProject(pId int64) ([]string, error) { // For Distinct Service Titles of the Project
	var res []string

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
		ColumnExpr("DISTINCT service_title").
		Where("project_id = ?", pId)

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

func GetTotalCostOfProjects(pid int64, daysFilter string) ([]types.CostSum, error) {

	res := []types.CostSum{}

	err := GetDB().Model(&types.DbServiceCostAws{}).
		ColumnExpr("date_trunc('day', usage_date)::date as usage_date").
		ColumnExpr("sum(accrued_cost_microdollar) as accrued_cost_microdollar").
		Where("usage_date >= CURRENT_DATE - INTERVAL ? day", daysFilter).
		Where("project_id = ?", pid).
		GroupExpr("date_trunc('day', usage_date)::date").
		Select(&res)

	if err != nil {
		fmt.Println("Error in GetTotalCostOfProjects: error is ", err)
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return res, nil
}

func GetServiceCostDateByProjectIdAndServicesAndCustomDate(pid int64, fromDate, todate string, services []string) ([]*types.DateServiceCost, error) {
	res := make([]*types.DateServiceCost, 0)

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
		Column("usage_date", "service_title", "accrued_cost_microdollar").
		Where("project_id = ? AND usage_date >= ? AND usage_date <= ? AND service_title = ANY(?)", pid, fromDate, todate, pg.Array(services)).
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

func GetServiceCostDateByProjectIdWithServiceAndYear(pId int64, year string, services []string) ([]*types.DateServiceCost, error) {
	res := make([]*types.DateServiceCost, 0)

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
		Column("usage_date", "service_title", "accrued_cost_microdollar").
		Where("project_id = ? AND EXTRACT(year FROM usage_date) = ? AND service_title = ANY(?)", pId, year, pg.Array(services)).
		Order("service_title").
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

func GetTotalCostOfProjectList(pId int64) (*types.TotalCost, error) {
	var totalcost types.TotalCost

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
		ColumnExpr("SUM(accrued_cost_microdollar) as accrued_cost_microdollar").
		Where("usage_date >= CURRENT_DATE - INTERVAL '30 days'").
		Where("project_id = ?", pId)

	err := query.Select(&totalcost)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &totalcost, nil
}

func GetDailyCostOfProject(pId int64) (types.TotalCost, error) {
	var dailycost types.TotalCost

	query := GetDB().
		Model(&types.DbServiceCostAws{}).
		ColumnExpr("SUM(accrued_cost_microdollar) as accrued_cost_microdollar").
		Where("usage_date >= CURRENT_DATE - INTERVAL '1 days'").
		Where("project_id = ?", pId)

	err := query.Select(&dailycost)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return types.TotalCost{}, nil
		} else {
			return types.TotalCost{}, err
		}
	}
	return dailycost, nil
}

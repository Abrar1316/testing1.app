package main

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20230202012125_add_service_cost_aws",
		up20230202012125AddServiceCostAws,
		down20230202012125AddServiceCostAws,
	)
}

func up20230202012125AddServiceCostAws(tx *pg.Tx) error {
	_, err := tx.Exec(`
		create table if not exists service_cost_aws(
			id bigserial primary key,
			project_id bigint references projects(id) on delete cascade,
			service_title character varying(512),
			usage_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			accrued_cost_microdollar bigint
		);
	`)
	return err
}

func down20230202012125AddServiceCostAws(tx *pg.Tx) error {
	_, err := tx.Exec(`
		drop table if exists service_cost_aws;
	`)
	return err
}

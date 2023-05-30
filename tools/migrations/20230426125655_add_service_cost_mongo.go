package main

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20230426125655_add_service_cost_mongo",
		up20230426125655AddServiceCostMongo,
		down20230426125655AddServiceCostMongo,
	)
}

func up20230426125655AddServiceCostMongo(tx *pg.Tx) error {
	_, err := tx.Exec(`
		create table if not exists service_cost_mongo(
			id bigserial primary key,
			project_id bigint references projects(id) on delete cascade,
			mongo_project_name character varying(512),
			service_title character varying(512),
			accrued_cost_microdollar bigint,
			usage_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			unit character varying(512),
			unit_price bigint,
			quantity bigint
		);
	`)
	return err
}

func down20230426125655AddServiceCostMongo(tx *pg.Tx) error {
	_, err := tx.Exec(`
		drop table if exists service_cost_mongo;
	`)
	return err
}

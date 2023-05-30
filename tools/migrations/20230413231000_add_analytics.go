package main

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20230413231000_add_analytics",
		up20230413231000AddAnalytics,
		down20230413231000AddAnalytics,
	)
}

func up20230413231000AddAnalytics(tx *pg.Tx) error {
	_, err := tx.Exec(`
	    create table if not exists analytics(
		  	id bigserial primary key,
			name character varying(512),
			email text ,
			daily_cost_microdollar bigint
		);
   `)
	return err
}

func down20230413231000AddAnalytics(tx *pg.Tx) error {
	_, err := tx.Exec(`
		drop table if exists analytics;
	`)
	return err
}

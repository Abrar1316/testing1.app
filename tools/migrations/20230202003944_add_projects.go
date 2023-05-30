package main

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20230202003944_add_projects",
		up20230202003944AddProjects,
		down20230202003944AddProjects,
	)
}

func up20230202003944AddProjects(tx *pg.Tx) error {
	_, err := tx.Exec(`
		create table if not exists projects(
			id bigserial primary key,
			name character varying(512),
			description varchar(255),
			aws_secret_key varchar(255),
			access_key varchar(255),
			is_active bool,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`)
	return err
}

func down20230202003944AddProjects(tx *pg.Tx) error {
	_, err := tx.Exec(`
		drop table if exists projects;
	`)
	return err
}

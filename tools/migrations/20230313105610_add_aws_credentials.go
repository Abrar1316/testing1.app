package main

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20230313105610_add_aws_credentials",
		up20230313105610AddAWSCredentials,
		down20230313105610AddAWSCredentials,
	)
}

func up20230313105610AddAWSCredentials(tx *pg.Tx) error {
	_, err := tx.Exec(`
		create table if not exists aws_credentials (
			id bigserial primary key,
			project_id bigint references projects(id) on delete cascade,
			aws_secret_key varchar(255) not null,
			access_key varchar(255) not null,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	
		insert into aws_credentials (project_id, aws_secret_key, access_key)
		select id, aws_secret_key, access_key from projects;
	
	`)
	return err
}

func down20230313105610AddAWSCredentials(tx *pg.Tx) error {
	_, err := tx.Exec(`
		drop table if exists aws_credentials;
	`)
	return err
}

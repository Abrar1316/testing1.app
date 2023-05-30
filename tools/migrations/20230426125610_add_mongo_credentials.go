package main

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20230426125610_add_mongo_credentials",
		up20230426125610AddMongoCredentials,
		down20230426125610AddMongoCredentials,
	)
}

func up20230426125610AddMongoCredentials(tx *pg.Tx) error {
	_, err := tx.Exec(`
		create table if not exists mongo_credentials (
			id bigserial primary key,
			project_id bigint references projects(id) on delete cascade,
			mongo_organization_id varchar(255) not null,
			public_key varchar(255) not null,
			private_key varchar(255) not null,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	
	`)
	return err
}

func down20230426125610AddMongoCredentials(tx *pg.Tx) error {
	_, err := tx.Exec(`
		drop table if exists mongo_credentials;
	`)
	return err
}

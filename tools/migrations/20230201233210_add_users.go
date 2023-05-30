package main

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20230201233210_add_users",
		up20230201233210AddUsers,
		down20230201233210AddUsers,
	)
}

func up20230201233210AddUsers(tx *pg.Tx) error {
	_, err := tx.Exec(`
		create table if not exists users(
			id bigserial primary key,
			name character varying(512),
			email text UNIQUE,
			password varchar(255),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`)
	return err
}

func down20230201233210AddUsers(tx *pg.Tx) error {
	_, err := tx.Exec(`
		drop table if exists users;
	`)
	return err
}

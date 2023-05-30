package main

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20230202005532_add_user_projects",
		up20230202005532AddUserProjects,
		down20230202005532AddUserProjects,
	)
}

func up20230202005532AddUserProjects(tx *pg.Tx) error {
	_, err := tx.Exec(`
		create table if not exists user_projects(
			id bigserial primary key,
			user_id bigint references users(id) on delete cascade,
			project_id bigint references projects(id) on delete cascade
		);
	`)
	return err
}

func down20230202005532AddUserProjects(tx *pg.Tx) error {
	_, err := tx.Exec(`
		drop table if exists user_projects;
	`)
	return err
}

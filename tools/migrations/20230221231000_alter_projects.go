package main

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20230221231000_alter_projects",
		up20230221231000AlterProject,
		down20230221231000AlterProject,
	)
}

func up20230221231000AlterProject(tx *pg.Tx) error {
	_, err := tx.Exec(`
			Alter table projects
				Add column is_pinned bool DEFAULT false;
	`)
	return err
}

func down20230221231000AlterProject(tx *pg.Tx) error {
	_, err := tx.Exec(`
			Alter table projects
				drop column if exists is_pinned;
	`)
	return err
}

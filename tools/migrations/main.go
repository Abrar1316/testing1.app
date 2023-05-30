package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TFTPL/AWS-Cost-Calculator/services/helpers/util"
	"github.com/TFTPL/AWS-Cost-Calculator/services/migrations"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

const usageText = `This program runs command on the database.
Supported commands are:
  - init - runs the specified initial migration as a batch on it's own.
  - migrate - runs all available migration.
  - rollback - reverts the last batch of migration.
  - create - creates a migration file.
Usage:
  migrations -command=[command] -name=[name of migration] -template=[filename of the template]`

func main() {

	var db *pg.DB
	util.InitZapLogger()
	flag.Usage = usage
	cmd := flag.String("command", "migrate", "Command that should be executed on the migration engine. Supported commands are: init, create, migrate and rollback")
	name := flag.String("name", "", "Name of the migration to be created")
	templateName := flag.String("template", "", "The filename of the template that should be used to create the migration")
	extra := flag.String("extra", "", "Extra parameters to pass to the command. Currently only migrate has an extra parameter caller one-by-one which runs the migrations in batches of one")
	test := flag.Bool("test", false, "Is migration for test db")

	flag.Parse()

	migrations.SetMigrationTableName("public.billapp_migrations")
	migrations.SetInitialMigration("000000000000_init")
	migrations.SetMigrationNameConvention(migrations.SnakeCase)

	template := ""

	if len(*templateName) > 0 {
		pwd, _ := os.Getwd()

		path := fmt.Sprintf("%s/%s/%s", pwd, "templates", *templateName)

		buf, err := os.ReadFile(path)

		if err != nil {
			zap.L().Fatal("file not found", zap.String("path", path))
		}

		template = string(buf)
	}

	if *cmd != "create" {
		cfg := util.GetBillAppConfigInstance()

		// config.ini has database credentials
		if !(*test) {
			fmt.Print("Loading config of dev db \n")
			cfg.LoadConfig(filepath.Join("..", "..", "config", "config.ini"))
		} else {
			fmt.Print("Loading config of test db \n")
			cfg.LoadConfig(filepath.Join("..", "..", "test", "config.ini"))
		}

		dbCnf := cfg.GetPostgresIni(false, "pg")

		db = pg.Connect(&pg.Options{
			Addr:     dbCnf.GetHost(),
			User:     dbCnf.GetUsername(),
			Database: dbCnf.GetDatabase(),
			Password: dbCnf.GetPassword(),
		})
	}

	var err error

	switch *cmd {
	case "create":
		err = migrations.Run(db, *cmd, *name, template)
	case "migrate":
		err = migrations.Run(db, *cmd, *extra)
	default:
		err = migrations.Run(db, *cmd)
	}

	if err != nil {
		zap.L().Fatal("Run", zap.Error(err))
	}
}

func usage() {
	fmt.Println(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

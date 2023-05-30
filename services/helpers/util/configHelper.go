package util

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	types "github.com/TFTPL/AWS-Cost-Calculator/services/types/ini"

	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

type BillAppConfig struct {
	configLoad sync.Once
	config     *ini.File
	relative   string
	loadedPath string
	testConfig bool
}

var instance *BillAppConfig
var once sync.Once

// GetBillAppConfigInstance returns go representation of config
func GetBillAppConfigInstance() *BillAppConfig {
	once.Do(func() {
		instance = &BillAppConfig{}
		instance.config = nil
	})
	return instance
}

func (m *BillAppConfig) LoadConfig(file string) error {
	// We need to block `sync.Once` otherwise subsequent
	// calls to `getConfig` will overwrite what we've loaded.
	m.configLoad.Do(func() {})

	config, err := ini.InsensitiveLoad(file)
	if err == nil {
		m.config = config
		m.loadedPath = file
	}
	return err
}

func (m *BillAppConfig) GetConfig() {
	m.configLoad.Do(func() {
		zap.L().Info("configHelper:getConfig - Std config.ini mode")
		path := filepath.Join("config", "config.ini")
		config, err := ini.InsensitiveLoad(path)

		if err == nil {
			m.config = config
			return
		}

		zap.L().Info("configHelper:getConfig:Load - Failed.", zap.Error(err))

		// If we get here there are tests or no config

		path = filepath.Join("..", "test", "config.ini")
		relative := filepath.Join("..")
		config, err = ini.InsensitiveLoad(path)

		for counter := 0; err != nil && counter < 5; counter++ {
			path = filepath.Join("..", path)
			relative = filepath.Join("..", relative)
			config, err = ini.InsensitiveLoad(path)
		}

		if err == nil {
			zap.L().Info("configHelper:getconfig - Loaded", zap.String("path", path))
			m.config = config
			m.relative = relative
			m.testConfig = true
			m.loadedPath = path
			return
		}

		zap.L().Fatal("configHelper:getConfig - Failed to find usable config")
	})
}

// testconfig returns testconfig bool
func (m *BillAppConfig) TestConfig() bool {
	return m.testConfig
}

func validTestServer(svr string) bool {
	return svr == "127.0.0.1" || svr == "localhost" || svr == "postgres"
}

func validTestPort(port string) bool {
	return true
}

func validTestDB(db string) bool {
	return strings.HasPrefix(db, "billapp-test") || db == "test"
}

// GetPostgresIni returns the Postgres config values
func (m *BillAppConfig) GetPostgresIni(fromTests bool, section string) *types.PostgresIni {
	m.GetConfig()

	host := m.config.Section(section).Key("host").String()
	db := m.config.Section(section).Key("database").String()

	hostParts := strings.Split(host, ":")

	if (m.testConfig || fromTests) && len(hostParts) > 1 {
		testServerValid := validTestServer(hostParts[0])
		testPortValid := validTestPort(hostParts[1])
		testDbValid := validTestDB(db)

		if !testServerValid || !testPortValid || !testDbValid {
			panic(fmt.Sprintf("You have loaded: %s\nYou are running tests\n[pg] host: %s\n[pg] database: %s\n\nYou can only use:\nhost: localhost or 127.0.0.1\n port: 5432 or 5433\n database: `billapp-test.*` or `test`\n\nwhen running tests\n\n", m.loadedPath, host, db))
		}
	}

	return types.NewPostgresIni(
		host,
		m.config.Section(section).Key("username").String(),
		m.config.Section(section).Key("password").String(),
		db,
	)
}

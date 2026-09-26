package db

import (
	"fmt"
	"sync"

	"github.com/CakeForKit/rsoi-lab1/internal/common/config"
	log "github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var postgresDB *gorm.DB
var postgresDBMutes sync.Mutex

func GetPostgresDataSource() *gorm.DB {
	postgresDBMutes.Lock()
	defer postgresDBMutes.Unlock()

	if postgresDB != nil {
		return postgresDB
	}

	postgresDB = newDataSource(config.CoreConfig.Postgres)
	return postgresDB
}

func newDataSource(properties config.GormProperty) *gorm.DB {
	dsn := properties.DSN
	if dsn == "" {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", properties.Username, properties.Password, properties.Host, properties.Port, properties.Database)
	}
	namingStrategy := schema.NamingStrategy{
		TablePrefix:   config.CoreConfig.SchemaDB + ".",
		SingularTable: true,
	}

	logLevel := logger.Warn
	if properties.ShowSql {
		logLevel = logger.Info
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logLevel),
		NamingStrategy: namingStrategy,
	})
	if err != nil {
		log.WithError(err).Error("Couldn't connect to database postgres")
		return nil
	}

	closer.Bind(func() {
		sqlDB, err := db.DB() // чтобы возвращался актуальый пул
		if err != nil {
			log.WithError(err).Errorf("Couldn't get sql db for postgres")
		}
		_ = sqlDB.Close()
		log.Debugf("Connection to database postgres closed")
	})
	return db
}

package casedb

import (
	"os"
	"time"

	llog "log"

	"github.com/case-management-suite/common"
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getDialector(conf config.DatabaseConfig) gorm.Dialector {
	switch conf.DatabaseType {
	case config.Sqlite:
		return sqlite.Open(conf.Address)
	case config.Postgres:
		return postgres.Open(conf.Address)
	default:
		return sqlite.Open("./default_cases.db")
	}
}

func NewSQLCaseStorageService(config config.DatabaseConfig) CaseStorageService {
	logLevel := logger.Error

	if config.LogSQL {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		llog.New(os.Stdout, "\r\n", llog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	db, _ := gorm.Open(getDialector(config), &gorm.Config{Logger: newLogger})

	db.AutoMigrate(models.CaseRecord{})
	db.AutoMigrate(models.CaseAction{CaseRecord: models.CaseRecord{}})
	return CaseStorageServiceImpl{DB: db, ServerUtils: common.NewTestServerUtils()}
}

type CaseStorageServiceImpl struct {
	DB *gorm.DB
	common.ServerUtils
}

func (c CaseStorageServiceImpl) SaveNewCase(id models.Identifier) error {
	c.Logger.Debug().Str("UUID", id).Msg("DB: Saving new case")
	newCase := models.CaseRecord{ID: id, Status: models.CaseStatusList.NewCase}
	result := c.DB.Create(&newCase)
	if result.Error != nil {
		c.Logger.Warn().Err(result.Error).Str("UUID", id).Msg("DB: Failed to save case")
	} else {
		c.Logger.Debug().Str("UUID", id).Msg("DB: Saved new case")
	}
	return result.Error
}

func (c CaseStorageServiceImpl) FindAllCases(spec models.CaseRecordSpec) ([]models.CaseRecord, error) {
	results := []models.CaseRecord{}
	status := c.DB.Select(spec.GetIncludedFieldsInSnakeCase()).Find(&results)
	if status.Error != nil {
		c.Logger.Warn().Err(status.Error).Msg("DB: Failed to find all cases")
	} else {
		c.Logger.Debug().Msg("DB: Found all cases")
	}
	return results, status.Error
}

func (c CaseStorageServiceImpl) FindCase(id models.Identifier, spec models.CaseRecordSpec) (models.CaseRecord, error) {
	record := models.CaseRecord{ID: id}
	status := c.DB.Select(spec.GetIncludedFieldsInSnakeCase()).First(&record)
	if status.Error != nil {
		c.Logger.Warn().Err(status.Error).Str("UUID", id).Msg("DB: Failed to find case")
	} else {
		c.Logger.Debug().Str("UUID", record.ID).Str("status", record.Status).Msg("DB: Found case")
	}
	return record, status.Error
}

func (c CaseStorageServiceImpl) UpdateCase(caseRecord *models.CaseRecord) error {
	result := c.DB.Save(&caseRecord)
	if result.Error != nil {
		c.Logger.Warn().Err(result.Error).Str("UUID", caseRecord.ID).Msg("DB: Failed to update case")
	} else {
		c.Logger.Debug().Str("UUID", caseRecord.ID).Str("status", caseRecord.Status).Msg("DB: Updated case")
	}
	return result.Error
}

func (c CaseStorageServiceImpl) SaveCaseContext(caseContext *models.CaseAction) error {
	result := c.DB.Create(&caseContext)
	if result.Error != nil {
		c.Logger.Warn().Err(result.Error).Str("UUID", caseContext.ID).Msg("DB: Failed to save case context")
	} else {
		c.Logger.Debug().Str("UUID", caseContext.ID).Msg("DB: Saved case context")
	}
	return result.Error
}

func (c CaseStorageServiceImpl) GetCaseContext(id models.Identifier, spec models.CaseActionSpec) (models.CaseAction, error) {
	record := models.CaseAction{ID: id}
	status := c.DB.Select(spec.GetIncludedFieldsInSnakeCase()).Find(&record).First(&record)
	return record, status.Error
}

func (c CaseStorageServiceImpl) GetContextForCase(caseId models.Identifier) ([]models.CaseAction, error) {
	contextList := []models.CaseAction{}
	status := c.DB.Where("case_record_id", caseId).Find(&contextList)
	return contextList, status.Error
}

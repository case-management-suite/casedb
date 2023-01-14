package casedb

import (
	"github.com/case-management-suite/common/config"
	"go.uber.org/fx"
)

func NewCaseDBConfig(params CaseDBParams) config.DatabaseConfig {
	return params.AppConfig.CasesStorage
}

type CaseDBParams struct {
	fx.In
	AppConfig config.AppConfig
}

type CaseDBResult struct {
	fx.Out
	CaseStorageService CaseStorageService
}

func ProvideCaseDBResult(dbConfig config.DatabaseConfig) CaseDBResult {
	storageService := NewSQLCaseStorageService(dbConfig)
	return CaseDBResult{CaseStorageService: storageService}
}

func NewFxCaseDBService() fx.Option {
	return fx.Module("casedb",
		fx.Provide(
			NewCaseDBConfig,
			fx.Private,
		),
		fx.Provide(ProvideCaseDBResult),
	)

}

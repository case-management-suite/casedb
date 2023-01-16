package casedb

import "github.com/case-management-suite/common/config"

type CaseStorageServiceFactory func(config.DatabaseConfig) CaseStorageService

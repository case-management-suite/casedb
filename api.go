package casedb

import "github.com/case-management-suite/models"

type CaseStorageService interface {
	SaveNewCase(id models.Identifier) error
	FindCase(id models.Identifier, spec models.CaseRecordSpec) (models.CaseRecord, error)
	FindAllCases(spec models.CaseRecordSpec) ([]models.CaseRecord, error)
	UpdateCase(caseRecord *models.CaseRecord) error
	SaveCaseContext(caseContext *models.CaseAction) error
	GetCaseContext(id models.Identifier, spec models.CaseActionSpec) (models.CaseAction, error)
	GetContextForCase(caseId models.Identifier) ([]models.CaseAction, error)
}

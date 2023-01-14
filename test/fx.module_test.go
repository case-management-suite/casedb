package casedb_test

import (
	"testing"

	"github.com/case-management-suite/casedb"
	"github.com/case-management-suite/testutil"
	"go.uber.org/fx/fxtest"
)

func TestCaseDBModule(t *testing.T) {
	testutil.AppFx(t, casedb.NewFxCaseDBService(), func(a *fxtest.App) {})
}

package dba

import "github.com/patrickarmengol/coffeetanuki/internal/errs"

func errInvalidFK(tableName string, fkName string, fkValue int64) *errs.Error {
	return errs.Errorf(errs.ERRUNPROCESSABLE, "invalid foreign key on table [%s] for field [%s] with value [%d]", tableName, fkName, fkValue)
}

func errRecordNotFound(tableName string, notFoundKey int64) *errs.Error {
	return errs.Errorf(errs.ERRNOTFOUND, "record not found on table [%s] for index key [%d]", tableName, notFoundKey)
}

func errEditConflict(tableName string, conflictKey int64) *errs.Error {
	return errs.Errorf(errs.ERRCONFLICT, "edit conflict on table [%s] for index key [%d]", tableName, conflictKey)
}

func errDuplicate(tableName string, dupeFieldName string, dupeValue any) *errs.Error {
	return errs.Errorf(errs.ERRCONFLICT, "duplicate on table [%s] for field [%s] with value [%q]", tableName, dupeFieldName, dupeValue)
}

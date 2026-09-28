package dberrx

import (
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

// ErrRecordNotFound 记录不存在。与 gorm.ErrRecordNotFound 为同一 sentinel，可用 errors.Is 匹配。
var ErrRecordNotFound = gorm.ErrRecordNotFound

// IsRecordNotFound 判断错误是否为记录不存在（含 wrap）。
func IsRecordNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows)
}

// IgnoreRecordNotFound 记录不存在时返回 nil，其它错误原样返回。
func IgnoreRecordNotFound(err error) error {
	if IsRecordNotFound(err) {
		return nil
	}
	return err
}

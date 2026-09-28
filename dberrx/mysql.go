package dberrx

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const (
	mysqlDeadlockErrorNumber        uint16 = 1213
	mysqlLockWaitTimeoutErrorNumber uint16 = 1205
	mysqlDuplicateEntryErrorNumber  uint16 = 1062
	mysqlSerializationSQLState             = "40001"
)

const (
	mysqlDeadlockErrorText       = "Error 1213 (40001):"
	mysqlDuplicateEntryErrorText = "Error 1062"
	sqliteUniqueConstraintText   = "UNIQUE constraint failed"
)

// IsMySQLRetryableTransactionError 判断错误是否为 MySQL 可重试事务错误。
func IsMySQLRetryableTransactionError(err error) bool {
	if mysqlErr, ok := AsMySQLError(err); ok {
		return isRetryableMySQLError(mysqlErr)
	}
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), mysqlDeadlockErrorText)
}

// IsDuplicateKey 判断错误是否为唯一键冲突。
func IsDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	// TranslateError 开启时 GORM 转成 ErrDuplicatedKey；本项目未开启，仍以 1062 为主。
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	if mysqlErr, ok := AsMySQLError(err); ok {
		return mysqlErr.Number == mysqlDuplicateEntryErrorNumber
	}
	// wrap 丢失类型时退回文本匹配，覆盖 MySQL 1062 与 SQLite unique constraint。
	msg := err.Error()
	return strings.Contains(msg, mysqlDuplicateEntryErrorText) ||
		strings.Contains(msg, sqliteUniqueConstraintText)
}

// AsMySQLError 从错误链取出 *mysql.MySQLError。
func AsMySQLError(err error) (*mysql.MySQLError, bool) {
	if err == nil {
		return nil, false
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr, true
	}
	return nil, false
}

// isRetryableMySQLError 判断 MySQL 原始错误是否适合事务重试。
func isRetryableMySQLError(err *mysql.MySQLError) bool {
	if err.Number == mysqlDeadlockErrorNumber || err.Number == mysqlLockWaitTimeoutErrorNumber {
		return true
	}
	return string(err.SQLState[:]) == mysqlSerializationSQLState
}

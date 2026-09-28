package dberrx

import (
	"fmt"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestIsMySQLRetryableTransactionError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "deadlock number",
			err: &mysql.MySQLError{
				Number:   1213,
				SQLState: [5]byte{'4', '0', '0', '0', '1'},
				Message:  "Deadlock found when trying to get lock",
			},
			want: true,
		},
		{
			name: "lock wait timeout number",
			err:  &mysql.MySQLError{Number: 1205, Message: "Lock wait timeout exceeded"},
			want: true,
		},
		{
			name: "serialization sql state",
			err:  &mysql.MySQLError{Number: 9999, SQLState: [5]byte{'4', '0', '0', '0', '1'}},
			want: true,
		},
		{
			name: "wrapped deadlock text",
			err:  fmt.Errorf("draw failed: %w", fmt.Errorf("Error 1213 (40001): Deadlock found when trying to get lock")),
			want: true,
		},
		{
			name: "duplicate key",
			err:  &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsMySQLRetryableTransactionError(tt.err))
		})
	}
}

func TestIsDuplicateKey(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{
			name: "mysql 1062",
			err:  &mysql.MySQLError{Number: 1062, Message: "Duplicate entry '1' for key 'uk_name'"},
			want: true,
		},
		{
			name: "gorm duplicated key",
			err:  gorm.ErrDuplicatedKey,
			want: true,
		},
		{
			name: "wrapped mysql 1062",
			err: fmt.Errorf("create record: %w", &mysql.MySQLError{
				Number:  1062,
				Message: "Duplicate entry",
			}),
			want: true,
		},
		{
			name: "wrapped gorm duplicated key",
			err:  fmt.Errorf("insert: %w", gorm.ErrDuplicatedKey),
			want: true,
		},
		{
			name: "wrapped 1062 text",
			err:  fmt.Errorf("create failed: %v", fmt.Errorf("Error 1062 (23000): Duplicate entry '1' for key 'uk_name'")),
			want: true,
		},
		{
			name: "sqlite unique constraint",
			err:  fmt.Errorf("UNIQUE constraint failed: pay_notify_events.event_id"),
			want: true,
		},
		{
			name: "deadlock",
			err:  &mysql.MySQLError{Number: 1213, Message: "Deadlock found when trying to get lock"},
			want: false,
		},
		{
			name: "record not found",
			err:  ErrRecordNotFound,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsDuplicateKey(tt.err))
		})
	}
}

func TestAsMySQLError(t *testing.T) {
	mysqlErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
	tests := []struct {
		name   string
		err    error
		wantOK bool
		number uint16
	}{
		{name: "nil", err: nil, wantOK: false},
		{name: "mysql error", err: mysqlErr, wantOK: true, number: 1062},
		{name: "wrapped mysql error", err: fmt.Errorf("insert: %w", mysqlErr), wantOK: true, number: 1062},
		{name: "non mysql error", err: fmt.Errorf("boom"), wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := AsMySQLError(tt.err)
			require.Equal(t, tt.wantOK, ok)
			if !tt.wantOK {
				require.Nil(t, got)
				return
			}
			require.Equal(t, tt.number, got.Number)
		})
	}
}

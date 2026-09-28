package dberrx

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestIsRecordNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "gorm record not found", err: gorm.ErrRecordNotFound, want: true},
		{name: "dberrx sentinel", err: ErrRecordNotFound, want: true},
		{name: "sql no rows", err: sql.ErrNoRows, want: true},
		{name: "wrapped gorm", err: fmt.Errorf("find user: %w", gorm.ErrRecordNotFound), want: true},
		{name: "wrapped sql no rows", err: fmt.Errorf("query: %w", sql.ErrNoRows), want: true},
		{name: "other error", err: errors.New("connection refused"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsRecordNotFound(tt.err))
		})
	}
}

func TestIgnoreRecordNotFound(t *testing.T) {
	other := errors.New("connection refused")
	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{name: "nil", err: nil, wantErr: nil},
		{name: "record not found", err: gorm.ErrRecordNotFound, wantErr: nil},
		{name: "sql no rows", err: sql.ErrNoRows, wantErr: nil},
		{name: "wrapped record not found", err: fmt.Errorf("find: %w", gorm.ErrRecordNotFound), wantErr: nil},
		{name: "other error", err: other, wantErr: other},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IgnoreRecordNotFound(tt.err)
			if tt.wantErr == nil {
				require.NoError(t, got)
				return
			}
			require.ErrorIs(t, got, tt.wantErr)
		})
	}
}

func TestErrRecordNotFoundAlias(t *testing.T) {
	require.ErrorIs(t, ErrRecordNotFound, gorm.ErrRecordNotFound)
	require.ErrorIs(t, gorm.ErrRecordNotFound, ErrRecordNotFound)
}

package user_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	Repo "github.com/maithuc2003/Test_GIN_golang/internal/repositories/user"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gdb, mock
}

func TestGetByUsername(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		mockExpectFn func(sqlmock.Sqlmock, string)
		expectErr    bool
	}{
		{
			name:     "user found",
			username: "john",
			mockExpectFn: func(mock sqlmock.Sqlmock, username string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).
					AddRow(1, username, "hashed_pw")

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"\."id" LIMIT \$2`).
					WithArgs(username, 1).
					WillReturnRows(rows)
			},
			expectErr: false,
		},
		{
			name:     "user not found",
			username: "unknown",
			mockExpectFn: func(mock sqlmock.Sqlmock, username string) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"\."id" LIMIT \$2`).
					WithArgs(username, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newMockDB(t)
			tt.mockExpectFn(mock, tt.username)
			repo := Repo.NewRepository(db)

			_, err := repo.GetByUsername(tt.username)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		mockExpectFn func(sqlmock.Sqlmock, string)
		expectErr    bool
	}{
		{
			name:     "user found for login",
			username: "john",
			mockExpectFn: func(mock sqlmock.Sqlmock, username string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"}).
					AddRow(1, username, "hashed_pw")

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"\."id" LIMIT \$2`).
					WithArgs(username, 1).
					WillReturnRows(rows)
			},
			expectErr: false,
		},
		{
			name:     "user not found for login",
			username: "unknown",
			mockExpectFn: func(mock sqlmock.Sqlmock, username string) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"\."id" LIMIT \$2`).
					WithArgs(username, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := newMockDB(t)
			tt.mockExpectFn(mock, tt.username)
			repo := Repo.NewRepository(db)

			_, err := repo.LoginUser(tt.username, "any_password")
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

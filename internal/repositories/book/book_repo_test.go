package book_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/maithuc2003/Test_GIN_golang/internal/repositories/book"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gdb, mock
}

func TestCreateBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		book         *models.Book
		mockExpectFn func(sqlmock.Sqlmock, *models.Book)
		expectErr    bool
	}{
		{
			name: "success",
			book: &models.Book{
				Title:    "Book 1",
				Stock:    10,
				AuthorID: 1,
			},
			mockExpectFn: func(mock sqlmock.Sqlmock, book *models.Book) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "books" ("title","stock","author_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
					WithArgs(book.Title, book.Stock, book.AuthorID, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "insert error",
			book: &models.Book{
				Title:    "Book 2",
				Stock:    0,
				AuthorID: 2,
			},
			mockExpectFn: func(mock sqlmock.Sqlmock, book *models.Book) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO "books" ("title","stock","author_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
					WithArgs(book.Title, book.Stock, book.AuthorID, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(gorm.ErrInvalidData)
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture loop variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := newMockDB(t)
			tt.mockExpectFn(mock, tt.book)
			repo := book.NewRepository(db)

			err := repo.CreateBook(tt.book)
			if (err != nil) != tt.expectErr {
				t.Errorf("CreateBook() error = %v, expectErr %v", err, tt.expectErr)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAllBooks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockExpectFn func(sqlmock.Sqlmock)
		wantCount    int
		expectErr    bool
	}{
		{
			name: "success with 2 books",
			mockExpectFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "stock", "author_id", "created_at", "updated_at"}).
					AddRow(1, "Book One", 5, 1, time.Now(), time.Now()).
					AddRow(2, "Book Two", 3, 2, time.Now(), time.Now())

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books"`)).
					WillReturnRows(rows)
			},
			wantCount: 2,
			expectErr: false,
		},
		{
			name: "query error",
			mockExpectFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books"`)).
					WillReturnError(gorm.ErrInvalidDB)
			},
			wantCount: 0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := newMockDB(t)
			tt.mockExpectFn(mock)
			repo := book.NewRepository(db)

			books, err := repo.GetAllBooks()
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAllBooks() error = %v, expectErr %v", err, tt.expectErr)
			}
			if len(books) != tt.wantCount {
				t.Errorf("GetAllBooks() got %d books, want %d", len(books), tt.wantCount)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestBookRepo_GetByBookID(t *testing.T) {
	t.Parallel()

	type args struct {
		id int
	}

	tests := []struct {
		name         string
		args         args
		mockExpectFn func(sqlmock.Sqlmock, int)
		wantTitle    string
		expectErr    bool
		errMsg       string
	}{
		{
			name: "success book found",
			args: args{id: 1},
			mockExpectFn: func(mock sqlmock.Sqlmock, id int) {
				now := time.Now()
				rows := sqlmock.NewRows([]string{"id", "title", "stock", "author_id", "created_at", "updated_at"}).
					AddRow(id, "Golang Mastery", 10, 2, now, now)

				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnRows(rows)
			},
			wantTitle: "Golang Mastery",
			expectErr: false,
		},
		{
			name: "book not found",
			args: args{id: 99},
			mockExpectFn: func(mock sqlmock.Sqlmock, id int) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectErr: true,
			errMsg:    "not found",
		},
		{
			name: "database error",
			args: args{id: 2},
			mockExpectFn: func(mock sqlmock.Sqlmock, id int) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnError(fmt.Errorf("db connection lost"))
			},
			expectErr: true,
			errMsg:    "db connection lost",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := newMockDB(t)
			tt.mockExpectFn(mock, tt.args.id)

			repo := book.NewRepository(db)
			gotBook, err := repo.GetByBookID(tt.args.id)

			if tt.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantTitle, gotBook.Title)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBookRepo_DeleteById(t *testing.T) {
	t.Parallel()

	type args struct {
		id int
	}

	tests := []struct {
		name         string
		args         args
		mockExpectFn func(mock sqlmock.Sqlmock, id int)
		expectErr    bool
		errMsg       string
		wantTitle    string
	}{
		{
			name: "success - delete book",
			args: args{id: 1},
			mockExpectFn: func(mock sqlmock.Sqlmock, id int) {
				now := time.Now()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "stock", "author_id", "created_at", "updated_at"}).
						AddRow(id, "Delete Me", 5, 1, now, now))

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "books" WHERE "books"."id" = $1`)).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantTitle: "Delete Me",
			expectErr: false,
		},
		{
			name: "book not found",
			args: args{id: 99},
			mockExpectFn: func(mock sqlmock.Sqlmock, id int) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectErr: true,
			errMsg:    "book with ID 99 not found",
		},
		{
			name: "delete failed",
			args: args{id: 2},
			mockExpectFn: func(mock sqlmock.Sqlmock, id int) {
				now := time.Now()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE "books"."id" = $1 ORDER BY "books"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "stock", "author_id", "created_at", "updated_at"}).
						AddRow(id, "Will Fail", 3, 1, now, now))

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "books" WHERE "books"."id" = $1`)).
					WithArgs(id).
					WillReturnError(fmt.Errorf("db error"))
				mock.ExpectRollback()
			},
			expectErr: true,
			errMsg:    "failed to delete book",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock := newMockDB(t)
			repo := book.NewRepository(db)

			tt.mockExpectFn(mock, tt.args.id)

			deletedBook, err := repo.DeleteById(tt.args.id)

			if tt.expectErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantTitle, deletedBook.Title)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBookRepo_UpdateById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	require.NoError(t, err)
	repo := book.NewRepository(gormDB)
	now := time.Now()
	tests := []struct {
		name         string
		book         *models.Book
		mockExpect   func()
		expectedErr  string
		expectResult bool
	}{
		{
			name: "success update",
			book: &models.Book{ID: 1, Title: "Updated Title", AuthorID: 1, Stock: 10, UpdatedAt: time.Now()},
			mockExpect: func() {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "authors" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				mock.ExpectBegin() // 👈 THÊM DÒNG NÀY

				mock.ExpectExec(`UPDATE "books" SET "title"=\$1,"stock"=\$2,"author_id"=\$3,"updated_at"=\$4 WHERE id = \$5`).
					WithArgs("Updated Title", 10, 1, sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit() // 👈 THÊM DÒNG NÀY
			},
			expectedErr:  "",
			expectResult: true,
		},
		{
			name: "author not found",
			book: &models.Book{ID: 2, Title: "New Title", AuthorID: 99, Stock: 10, UpdatedAt: now},
			mockExpect: func() {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "authors" WHERE id = \$1`).
					WithArgs(99).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedErr:  "author not found",
			expectResult: false,
		},
		{
			name: "book not found",
			book: &models.Book{ID: 999, Title: "Title", AuthorID: 1, Stock: 5, UpdatedAt: now},
			mockExpect: func() {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "authors" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "books" SET "title"=\$1,"stock"=\$2,"author_id"=\$3,"updated_at"=\$4 WHERE id = \$5`).
					WithArgs("Title", 5, 1, sqlmock.AnyArg(), 999).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectCommit()
			},
			expectedErr:  "no book updated",
			expectResult: false,
		},
		{
			name: "update error",
			book: &models.Book{ID: 3, Title: "Error Title", AuthorID: 1, Stock: 5, UpdatedAt: now},
			mockExpect: func() {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "authors" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "books" SET "title"=\$1,"stock"=\$2,"author_id"=\$3,"updated_at"=\$4 WHERE id = \$5`).
					WithArgs("Error Title", 5, 1, sqlmock.AnyArg(), 3).
					WillReturnError(errors.New("failed to update book"))

				mock.ExpectRollback()
			},
			expectedErr:  "failed to update book",
			expectResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpect()
			result, err := repo.UpdateById(tt.book)

			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			if tt.expectResult {
				require.NotNil(t, result)
				require.Equal(t, tt.book.Title, result.Title)
			} else {
				require.Nil(t, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

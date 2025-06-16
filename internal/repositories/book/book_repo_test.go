package book_test

import (
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

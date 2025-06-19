package author_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/maithuc2003/Test_GIN_golang/internal/repositories/author"

	sqlitedriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"

	// Import driver SQLite thuần Go
	_ "modernc.org/sqlite"
)

// setupTestDB dùng modernc.org/sqlite driver, không cần cgo
func setupTestDB(t *testing.T) *gorm.DB {
	// Tạo config cho sqlite driver với driverName "sqlite" là driver của modernc
	db, err := gorm.Open(sqlitedriver.New(sqlitedriver.Config{
		DSN:        "file::memory:?cache=shared",
		DriverName: "sqlite",
	}), &gorm.Config{})

	require.NoError(t, err)

	err = db.AutoMigrate(&models.Author{})
	require.NoError(t, err)

	return db
}

func TestAuthorRepo_GetAllAuthors(t *testing.T) {
	db := setupTestDB(t)
	repo := author.NewAuthorRepo(db)

	// Seed data
	authors := []models.Author{
		{Name: "Author1", Nationality: "US"},
		{Name: "Author2", Nationality: "UK"},
	}
	for _, a := range authors {
		err := repo.CreateAuthor(&a)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Get all authors success",
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetAllAuthors()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, tt.wantCount)
		})
	}
}
func TestAuthorRepo_GetByAuthorID(t *testing.T) {
	db := setupTestDB(t)
	repo := author.NewAuthorRepo(db)

	author := models.Author{Name: "Author1", Nationality: "US"}
	err := repo.CreateAuthor(&author)
	require.NoError(t, err)

	tests := []struct {
		name     string
		id       int
		wantErr  bool
		wantName string
	}{
		{"Existing author", author.ID, false, "Author1"},
		{"Non-existing author", 999, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByAuthorID(tt.id)
			if tt.wantErr {
				// Kiểm tra có lỗi, và object trả về phải nil hoặc empty
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got) // got không được nil nếu không lỗi
			assert.Equal(t, tt.wantName, got.Name)
		})
	}
}

func TestAuthorRepo_CreateAuthor(t *testing.T) {
	db := setupTestDB(t)
	repo := author.NewAuthorRepo(db)

	tests := []struct {
		name    string
		input   models.Author
		wantErr bool
	}{
		{"Valid author", models.Author{Name: "Author1", Nationality: "US"}, false},
		{"Empty name", models.Author{Name: "", Nationality: "VN"}, false}, // depends on validation
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateAuthor(&tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotZero(t, tt.input.ID)
		})
	}
}

func TestAuthorRepo_DeleteById(t *testing.T) {
	db := setupTestDB(t)
	repo := author.NewAuthorRepo(db)

	author := models.Author{Name: "AuthorToDelete", Nationality: "US"}
	err := repo.CreateAuthor(&author)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"Delete existing author", author.ID, false},
		{"Delete non-existing author", 999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.DeleteById(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.id, got.ID)
		})
	}
}

func TestAuthorRepo_UpdateById(t *testing.T) {
	db := setupTestDB(t)
	repo := author.NewAuthorRepo(db)

	author := models.Author{Name: "Original", Nationality: "US", UpdatedAt: time.Now()}
	err := repo.CreateAuthor(&author)
	require.NoError(t, err)

	tests := []struct {
		name       string
		updateData models.Author
		wantErr    bool
		wantName   string
		wantNation string
	}{
		{
			name:       "Update existing author",
			updateData: models.Author{ID: author.ID, Name: "Updated", Nationality: "UK", UpdatedAt: time.Now()},
			wantErr:    false,
			wantName:   "Updated",
			wantNation: "UK",
		},
		{
			name:       "Update non-existing author",
			updateData: models.Author{ID: 9999, Name: "NoOne"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.UpdateById(&tt.updateData)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantName, got.Name)
			assert.Equal(t, tt.wantNation, got.Nationality)
		})
	}
}

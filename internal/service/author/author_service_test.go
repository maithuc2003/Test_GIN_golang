package author_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	mockrepo "github.com/maithuc2003/Test_GIN_golang/internal/mocks/repositories"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/maithuc2003/Test_GIN_golang/internal/service/author"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuthor(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		input         *models.Author
		mockAuthors   []*models.Author
		mockGetAllErr error
		mockCreateErr error
		expectError   bool
		errorMessage  string
	}{
		{
			name:         "nil author",
			input:        nil,
			expectError:  true,
			errorMessage: "author is nil",
		},
		{
			name:         "empty author name",
			input:        &models.Author{Name: "   ", Nationality: "VN"},
			expectError:  true,
			errorMessage: "author name cannot be empty",
		},
		{
			name:  "duplicate author name",
			input: &models.Author{Name: "John", Nationality: "US"},
			mockAuthors: []*models.Author{
				{Name: "john"},
			},
			expectError:  true,
			errorMessage: "author with the same name already exists",
		},
		{
			name:          "get all authors error",
			input:         &models.Author{Name: "NewAuthor"},
			mockGetAllErr: errors.New("DB error"),
			expectError:   true,
			errorMessage:  "failed to fetch authors for validation: DB error",
		},
		{
			name:          "create author error",
			input:         &models.Author{Name: "NewAuthor"},
			mockAuthors:   []*models.Author{},
			mockCreateErr: errors.New("insert error"),
			expectError:   true,
			errorMessage:  "failed to create author: insert error",
		},
		{
			name:        "successful create",
			input:       &models.Author{Name: "Unique", Nationality: "JP", CreatedAt: now, UpdatedAt: now},
			mockAuthors: []*models.Author{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockAuthorRepository)
			svc := author.NewAuthorService(mockRepo)

			if tt.input != nil {
				// Không gọi GetAll nếu tên rỗng (early return)
				if tt.input.Name != "   " {
					mockRepo.On("GetAllAuthors").Return(tt.mockAuthors, tt.mockGetAllErr)

					if tt.mockGetAllErr == nil {
						// Check nếu không trùng tên và không expect lỗi → gọi CreateAuthor
						isDuplicate := false
						for _, a := range tt.mockAuthors {
							if a.Name == tt.input.Name || a.Name == "john" && tt.input.Name == "John" {
								isDuplicate = true
								break
							}
						}

						if !isDuplicate {
							mockRepo.On("CreateAuthor", tt.input).Return(tt.mockCreateErr)
						}
					}
				}
			}

			err := svc.CreateAuthor(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.errorMessage)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllAuthors(t *testing.T) {
	tests := []struct {
		name         string
		mockAuthors  []*models.Author
		mockError    error
		expectError  bool
		errorMessage string
	}{
		{
			name:         "repository error",
			mockAuthors:  nil,
			mockError:    errors.New("db error"),
			expectError:  true,
			errorMessage: "db error",
		},
		{
			name:         "no authors found",
			mockAuthors:  []*models.Author{},
			mockError:    nil,
			expectError:  true,
			errorMessage: "no authors found in the system",
		},
		{
			name:        "success",
			mockAuthors: []*models.Author{{ID: 1, Name: "John"}},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockAuthorRepository)
			svc := author.NewAuthorService(mockRepo)

			mockRepo.On("GetAllAuthors").Return(tt.mockAuthors, tt.mockError)

			result, err := svc.GetAllAuthors()

			if tt.expectError {
				assert.Nil(t, result)
				assert.EqualError(t, err, tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockAuthors, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetByAuthorID(t *testing.T) {
	tests := []struct {
		name         string
		inputID      int
		mockResult   *models.Author
		mockError    error
		expectError  bool
		errorMessage string
	}{
		{
			name:         "invalid ID",
			inputID:      0,
			expectError:  true,
			errorMessage: "invalid author ID",
		},
		{
			name:         "repository error",
			inputID:      1,
			mockResult:   nil,
			mockError:    errors.New("db failure"),
			expectError:  true,
			errorMessage: "failed to retrieve author: db failure",
		},
		{
			name:         "author not found",
			inputID:      2,
			mockResult:   nil,
			mockError:    nil,
			expectError:  true,
			errorMessage: "author not found",
		},
		{
			name:        "success",
			inputID:     3,
			mockResult:  &models.Author{ID: 3, Name: "Jane"},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockAuthorRepository)
			svc := author.NewAuthorService(mockRepo)

			if tt.inputID > 0 {
				mockRepo.On("GetByAuthorID", tt.inputID).Return(tt.mockResult, tt.mockError)
			}

			result, err := svc.GetByAuthorID(tt.inputID)

			if tt.expectError {
				assert.Nil(t, result)
				assert.EqualError(t, err, tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResult, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteById(t *testing.T) {
	tests := []struct {
		name         string
		inputID      int
		mockResult   *models.Author
		mockError    error
		expectError  bool
		errorMessage string
	}{
		{
			name:         "invalid ID",
			inputID:      0,
			expectError:  true,
			errorMessage: "invalid author ID",
		},
		{
			name:         "repository error",
			inputID:      1,
			mockError:    errors.New("db error"),
			expectError:  true,
			errorMessage: "failed to delete author: db error",
		},
		{
			name:         "author not found",
			inputID:      2,
			mockResult:   nil,
			mockError:    nil,
			expectError:  true,
			errorMessage: "author not found or already deleted",
		},
		{
			name:        "successfully deleted",
			inputID:     3,
			mockResult:  &models.Author{ID: 3, Name: "Deleted Author"},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockAuthorRepository)
			svc := author.NewAuthorService(mockRepo)

			if tt.inputID > 0 {
				mockRepo.On("DeleteById", tt.inputID).Return(tt.mockResult, tt.mockError)
			}

			result, err := svc.DeleteById(tt.inputID)

			if tt.expectError {
				assert.Nil(t, result)
				assert.EqualError(t, err, tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResult, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateById(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name           string
		inputAuthor    *models.Author
		existingAuthor *models.Author
		allAuthors     []*models.Author
		mockGetErr     error
		mockGetAllErr  error
		mockUpdateErr  error
		mockUpdateRes  *models.Author
		expectError    bool
		errorMessage   string
	}{
		{
			name:         "nil author",
			inputAuthor:  nil,
			expectError:  true,
			errorMessage: "author is nil",
		},
		{
			name:         "invalid author ID",
			inputAuthor:  &models.Author{ID: 0},
			expectError:  true,
			errorMessage: "invalid author ID",
		},
		{
			name:         "empty author name",
			inputAuthor:  &models.Author{ID: 1, Name: "   "},
			expectError:  true,
			errorMessage: "author name cannot be empty",
		},
		{
			name:         "get by ID error",
			inputAuthor:  &models.Author{ID: 1, Name: "Valid"},
			mockGetErr:   errors.New("db error"),
			expectError:  true,
			errorMessage: "failed to fetch existing author: db error",
		},
		{
			name:           "author not found",
			inputAuthor:    &models.Author{ID: 1, Name: "Valid"},
			existingAuthor: nil,
			expectError:    true,
			errorMessage:   "author not found",
		},
		{
			name:           "duplicate name found",
			inputAuthor:    &models.Author{ID: 2, Name: "Jane"},
			existingAuthor: &models.Author{ID: 2, Name: "Jane"},
			allAuthors:     []*models.Author{{ID: 3, Name: "jane"}},
			expectError:    true,
			errorMessage:   "another author with the same name already exists",
		},
		{
			name:           "error when getting all authors",
			inputAuthor:    &models.Author{ID: 2, Name: "Jane"},
			existingAuthor: &models.Author{ID: 2, Name: "Jane"},
			mockGetAllErr:  errors.New("getAll error"),
			expectError:    true,
			errorMessage:   "failed to validate author name: getAll error",
		},
		{
			name:           "update error",
			inputAuthor:    &models.Author{ID: 2, Name: "Jane"},
			existingAuthor: &models.Author{ID: 2, Name: "Jane"},
			allAuthors:     []*models.Author{{ID: 2, Name: "Jane"}},
			mockUpdateErr:  errors.New("update failed"),
			expectError:    true,
			errorMessage:   "failed to update author : update failed",
		},
		{
			name:           "successful update",
			inputAuthor:    &models.Author{ID: 2, Name: "Updated", Nationality: "US", CreatedAt: now, UpdatedAt: now},
			existingAuthor: &models.Author{ID: 2, Name: "OldName"},
			allAuthors:     []*models.Author{{ID: 2, Name: "OldName"}},
			mockUpdateRes:  &models.Author{ID: 2, Name: "Updated", Nationality: "US", CreatedAt: now, UpdatedAt: now},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockrepo.MockAuthorRepository)
			svc := author.NewAuthorService(mockRepo)
			if tt.inputAuthor != nil && tt.inputAuthor.ID > 0 && strings.TrimSpace(tt.inputAuthor.Name) != "" {
				mockRepo.On("GetByAuthorID", tt.inputAuthor.ID).Return(tt.existingAuthor, tt.mockGetErr)

				if tt.mockGetErr == nil && tt.existingAuthor != nil {
					mockRepo.On("GetAllAuthors").Return(tt.allAuthors, tt.mockGetAllErr)

					if tt.mockGetAllErr == nil && (tt.mockUpdateRes != nil || tt.mockUpdateErr != nil) {
						mockRepo.On("UpdateById", tt.inputAuthor).Return(tt.mockUpdateRes, tt.mockUpdateErr)
					}
				}
			}

			result, err := svc.UpdateById(tt.inputAuthor)

			if tt.expectError {
				assert.Nil(t, result)
				assert.EqualError(t, err, tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockUpdateRes, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

package database

import (
	"github.com/maithuc2003/Test_GIN_golang/config"
)

func HasAccess(userID int, accessName string) bool {
	var count int
	query := `SELECT COUNT(*)
		FROM users u
		JOIN user_role ur ON u.id = ur.user_id
		JOIN role_access ra ON ur.role_id = ra.role_id
		JOIN access a ON ra.access_id = a.access_id
		WHERE u.id = ? AND a.access_name = ?`

	sqlDB, err := config.DB.DB() // lấy *sql.DB từ GORM
	if err != nil {
		return false
	}

	err = sqlDB.QueryRow(query, userID, accessName).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

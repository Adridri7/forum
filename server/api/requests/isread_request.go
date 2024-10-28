package requests

import (
	"database/sql"
	"fmt"
	"forum/server"
	"net/http"
	"time"
)

type AdminRequest struct {
	RequestUUID string    `json:"request_uuid"`
	UserUUID    string    `json:"user_uuid"`
	CreatedAt   time.Time `json:"created_at"`
	IsRead      bool      `json:"isRead"`
	Content     string    `json:"content"`
}

func MarkRequestsAsReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := MarkRequestsAsRead(server.Db); err != nil {
		http.Error(w, "Failed to mark requests as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func FetchUnreadRequests(db *sql.DB) ([]AdminRequest, error) {
	query := `
    SELECT * FROM admin_requests
    WHERE is_read = FALSE
    ORDER BY created_at DESC`

	rows, err := server.RunQuery(query)
	if err != nil {
		return nil, fmt.Errorf("could not fetch requests: %v", err)
	}

	var adminRequests []AdminRequest
	for _, row := range rows {
		request := AdminRequest{
			RequestUUID: GetStringFromRow(row["request_uuid"]),
			UserUUID:    GetStringFromRow(row["user_uuid"]),
			CreatedAt:   row["created_at"].(time.Time),
			IsRead:      row["isRead"].(bool),
			Content:     GetStringFromRow(row["content"]),
		}
		adminRequests = append(adminRequests, request)
	}
	return adminRequests, nil
}

func GetStringFromRow(value interface{}) string {
	if value == nil {
		return ""
	}
	return value.(string)
}

func MarkRequestsAsRead(db *sql.DB) error {
	query := `
    UPDATE admin_requests
    SET isRead = TRUE
    WHERE isRead = False`

	_, err := server.RunQuery(query)
	if err != nil {
		return fmt.Errorf("database query failed: %v", err)
	}
	return nil
}

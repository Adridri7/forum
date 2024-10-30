package requests

import (
	"database/sql"
	"encoding/json"
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
	var params map[string]interface{}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	request_uuid, ok := params["request_uuid"].(string)
	if request_uuid == "" || !ok {
		http.Error(w, "missing argument in request payload", http.StatusBadRequest)
		return
	}

	if err := MarkRequestsAsRead(server.Db, request_uuid); err != nil {
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

func MarkRequestsAsRead(db *sql.DB, request_uuid string) error {
	query := `
    UPDATE admin_requests
    SET isRead = TRUE
    WHERE isRead = False AND request_uuid = ?`

	_, err := server.RunQuery(query, request_uuid)
	if err != nil {
		return fmt.Errorf("database query failed: %v", err)
	}
	return nil
}

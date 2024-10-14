package notifications

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/server"
	posts "forum/server/utils"
	"net/http"
	"strings"
	"time"
)

type Notification struct {
	NotificationID string    `json:"notification_id"`
	UserUUID       string    `json:"user_uuid"`
	Action         string    `json:"action"`
	TargetType     string    `json:"target_type"`
	ReferenceID    string    `json:"reference_id"`
	CreatedAt      time.Time `json:"created_at"`
	IsRead         bool      `json:"is_read"`
	Username       string    `json:"username"`
}

func FetchUnreadNotificationsHandler(w http.ResponseWriter, r *http.Request) {

	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user_UUID := params["user_uuid"].(string)

	notifications, err := FetchUnreadNotifications(server.Db, user_UUID)
	if err != nil {
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(notifications); err != nil {
		http.Error(w, "Failed to encode notifications", http.StatusInternalServerError)
		return
	}
}

func MarkNotificationsAsReadHandler(w http.ResponseWriter, r *http.Request) {
	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userUUID := params["user_uuid"].(string)

	notifications, err := FetchUnreadNotifications(server.Db, userUUID)
	if err != nil {
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}

	var notificationIDs []string
	for _, notif := range notifications {
		notificationIDs = append(notificationIDs, notif.NotificationID)
	}

	if err := MarkNotificationsAsRead(server.Db, userUUID, notificationIDs); err != nil {
		http.Error(w, "Failed to mark notifications as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func FetchUnreadNotifications(db *sql.DB, userUUID string) ([]Notification, error) {
	query := `
    SELECT notification_id, user_uuid, action, target_type, reference_id, created_at, is_read, username
    FROM notifications
    WHERE user_uuid = ? AND is_read = FALSE
    ORDER BY created_at DESC`

	rows, err := server.RunQuery(query, userUUID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch notifications: %v", err)
	}

	var notifications []Notification
	for _, row := range rows {
		notification := Notification{
			NotificationID: GetStringFromRow(row["notification_id"]),
			UserUUID:       GetStringFromRow(row["user_uuid"]),
			Action:         GetStringFromRow(row["action"]),
			TargetType:     GetStringFromRow(row["target_type"]),
			ReferenceID:    GetStringFromRow(row["reference_id"]),
			CreatedAt:      row["created_at"].(time.Time),
			IsRead:         row["is_read"].(bool),
			Username:       GetStringFromRow(row["username"]),
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func GetStringFromRow(value interface{}) string {
	if value == nil {
		return "" // ou autre valeur par défaut
	}
	return value.(string)
}

func MarkNotificationsAsRead(db *sql.DB, userUUID string, notificationIDs []string) error {
	// Convertir []string en []interface{}
	interfaceIDs := make([]interface{}, len(notificationIDs))
	for i, v := range notificationIDs {
		interfaceIDs[i] = v
	}

	query := `
    UPDATE notifications
    SET is_read = TRUE
    WHERE user_uuid = ? AND notification_id IN (?` + strings.Repeat(",?", len(interfaceIDs)-1) + `)`

	args := append([]interface{}{userUUID}, interfaceIDs...)

	_, err := server.RunQuery(query, args...)
	if err != nil {
		return fmt.Errorf("database query failed: %v", err)
	}
	return nil
}

func InsertNotification(db *sql.DB, userUUID, action, targetType, referenceID string, username string) error {
	notification_UUID, _ := posts.GenerateUUID()

	// usernameQuery := `
	// SELECT username FROM users
	// WHERE user_uuid = ?`

	// rows, err := server.RunQuery(usernameQuery, userUUID)
	// if err != nil {
	// 	return fmt.Errorf("database query failed: %v", err)
	// }

	// var username string

	// for _, row := range rows {
	// 	username = row["username"].(string)
	// }

	query := `
    INSERT INTO notifications (notification_id, user_uuid, action, target_type, reference_id, created_at, is_read, username)
    VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, FALSE, ?)`

	_, err := server.RunQuery(query, notification_UUID, userUUID, action, targetType, referenceID, username)
	if err != nil {
		return fmt.Errorf("failed to insert notification: %v", err)
	}
	return nil
}

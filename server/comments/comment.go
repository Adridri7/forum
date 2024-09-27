package comments

import "time"

type Comment struct {
	Comment_id int       `json:"comment_id"`
	Post_uuid  string    `json:"post_uuid"`
	User_uuid  string    `json:"user_uuid"`
	Content    string    `json:"content"`
	Created_at time.Time `json:"created_at"`
}

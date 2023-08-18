// models/comment.go
package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	BlogID    int       `json:"blog_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

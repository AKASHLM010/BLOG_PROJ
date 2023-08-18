package models

import "time"

type Blog struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    UserID    int       `json:"user_id"`
    Image     string    `json:"image"` // New field for image URL
    Tags      string    `json:"tags"`  // Change type to string
}

package models

import "time"

type Skill struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	Proficiency int8   `json:"proficiency"` // e.g. "Beginner", "Intermediate", "Advanced"
}

type Experience struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"` // empty = "Present"
	Description string `json:"description"`
	Type        string `json:"type"` // "education" or "experience"
}

type Project struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	LongDesc    string    `json:"long_desc"`
	ImageURL    string    `json:"image_url"`
	RepoURL     string    `json:"repo_url"`
	LiveURL     string    `json:"live_url"`
	Tags        string    `json:"tags"` // comma-separated
	CreatedAt   time.Time `json:"created_at"`
}

type BlogPost struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Excerpt   string    `json:"excerpt"`
	Content   string    `json:"content"`
	Tags      string    `json:"tags"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

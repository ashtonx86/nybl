package schemas

import "time"

type Message struct {
	ID             string
	Content        string
	AttachmentUrls []string
	ChatID         string // ChatID is equivalent to ID of Space
	AuthorID       string

	CreatedAt time.Time 
	UpdatedAt time.Time
}
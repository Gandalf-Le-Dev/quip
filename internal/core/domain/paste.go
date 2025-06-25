package domain

import (
	"time"

	"github.com/go-enry/go-enry/v2"
)

type Paste struct {
	ID        string
	Content   string
	Language  string
	Title     string
	Views     int
	MaxViews  int
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewPaste(content, language, title string, ttl time.Duration) *Paste {
	if language == "" {
		language = detectLanguage(content)
	}

	return &Paste{
		ID:        generateID(),
		Content:   content,
		Language:  language,
		Title:     title,
		Views:     0,
		MaxViews:  -1, // unlimited
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (p *Paste) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

func (p *Paste) CanView() bool {
	if p.IsExpired() {
		return false
	}
	if p.MaxViews > 0 && p.Views >= p.MaxViews {
		return false
	}
	return true
}

func (p *Paste) IncrementViews() {
	p.Views++
}

func detectLanguage(content string) string {
	return enry.GetLanguage("", []byte(content))
}

// Package models provides a low-level implementation of the basic structures present in the document system.
package models

import (
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Model is a wrapper around gorm.Model allowing methods to be attached to it.
type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// prepareString properly escapes all characters in a string.
// Used to prevent injection attacks.
func prepareString(s string) string {
	return html.EscapeString(strings.TrimSpace(s))
}

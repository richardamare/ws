// Package session manages curated Claude Code session bookmarks for a project:
// a short, named list stored in the project config so good-context sessions can
// be reopened with `ws resume` instead of scrolling Claude's full history.
package session

import (
	"fmt"

	"github.com/richardamare/ws/internal/config"
)

// Find returns the bookmark with the given label.
func Find(bookmarks []config.Bookmark, label string) (config.Bookmark, bool) {
	for _, b := range bookmarks {
		if b.Label == label {
			return b, true
		}
	}
	return config.Bookmark{}, false
}

// Upsert adds or updates a bookmark by label, returning the new slice.
func Upsert(bookmarks []config.Bookmark, b config.Bookmark) []config.Bookmark {
	for i := range bookmarks {
		if bookmarks[i].Label == b.Label {
			bookmarks[i] = b
			return bookmarks
		}
	}
	return append(bookmarks, b)
}

// ResumeArgs builds the `claude --resume <id>` arguments for a bookmark.
func ResumeArgs(b config.Bookmark) ([]string, error) {
	if b.ID == "" {
		return nil, fmt.Errorf("bookmark %q has no session id", b.Label)
	}
	return []string{"--resume", b.ID}, nil
}

package session

import (
	"testing"

	"github.com/richardamare/ws/internal/config"
)

func TestUpsertAddsAndUpdates(t *testing.T) {
	var bm []config.Bookmark
	bm = Upsert(bm, config.Bookmark{Label: "a", ID: "1"})
	bm = Upsert(bm, config.Bookmark{Label: "b", ID: "2"})
	if len(bm) != 2 {
		t.Fatalf("expected 2, got %d", len(bm))
	}
	bm = Upsert(bm, config.Bookmark{Label: "a", ID: "99", Note: "updated"})
	if len(bm) != 2 {
		t.Fatalf("upsert should not grow on update, got %d", len(bm))
	}
	got, ok := Find(bm, "a")
	if !ok || got.ID != "99" || got.Note != "updated" {
		t.Errorf("update not applied: %+v", got)
	}
}

func TestFindMissing(t *testing.T) {
	if _, ok := Find(nil, "x"); ok {
		t.Error("expected not found")
	}
}

func TestResumeArgs(t *testing.T) {
	args, err := ResumeArgs(config.Bookmark{Label: "a", ID: "abc"})
	if err != nil {
		t.Fatal(err)
	}
	if len(args) != 2 || args[0] != "--resume" || args[1] != "abc" {
		t.Errorf("got %v", args)
	}
	if _, err := ResumeArgs(config.Bookmark{Label: "a"}); err == nil {
		t.Error("expected error for empty id")
	}
}

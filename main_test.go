package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func makeFile(t *testing.T, dir, name string, modified time.Time) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, modified, modified); err != nil {
		t.Fatal(err)
	}
}

func TestOrganiseMovesByDate(t *testing.T) {
	dir := t.TempDir()
	modified := time.Date(2025, 3, 4, 12, 0, 0, 0, time.Local)
	makeFile(t, dir, "receipt.pdf", modified)
	s, err := organise(dir, options{yes: true}, strings.NewReader(""), &bytes.Buffer{}, modified.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if s.moved != 1 {
		t.Fatalf("moved = %d, want 1", s.moved)
	}
	if _, err := os.Stat(filepath.Join(dir, "2025-03-04", "receipt.pdf")); err != nil {
		t.Fatal(err)
	}
}

func TestDryRunDoesNotCreateDirectories(t *testing.T) {
	dir := t.TempDir()
	modified := time.Date(2025, 3, 4, 12, 0, 0, 0, time.Local)
	makeFile(t, dir, "receipt.pdf", modified)
	if _, err := organise(dir, options{dryRun: true}, strings.NewReader(""), &bytes.Buffer{}, modified); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "2025-03-04")); !os.IsNotExist(err) {
		t.Fatalf("dry run created directory: %v", err)
	}
}

func TestConflictDoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	modified := time.Date(2025, 3, 4, 12, 0, 0, 0, time.Local)
	makeFile(t, dir, "receipt.pdf", modified)
	dest := filepath.Join(dir, "2025-03-04")
	if err := os.Mkdir(dest, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dest, "receipt.pdf"), []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}
	s, err := organise(dir, options{yes: true}, strings.NewReader(""), &bytes.Buffer{}, modified)
	if err != nil {
		t.Fatal(err)
	}
	if s.conflicts != 1 {
		t.Fatalf("conflicts = %d, want 1", s.conflicts)
	}
	contents, _ := os.ReadFile(filepath.Join(dest, "receipt.pdf"))
	if string(contents) != "original" {
		t.Fatal("destination was overwritten")
	}
}

func TestOlderThanSkipsRecentFiles(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	makeFile(t, dir, "still-downloading.zip", now.Add(-time.Minute))
	s, err := organise(dir, options{yes: true, minAge: 10 * time.Minute}, strings.NewReader(""), &bytes.Buffer{}, now)
	if err != nil {
		t.Fatal(err)
	}
	if s.skipped != 1 || s.moved != 0 {
		t.Fatalf("unexpected stats: %+v", s)
	}
}

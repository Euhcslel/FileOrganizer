package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetFileCategory(t *testing.T) {
	cases := []struct {
		ext string
		want string
	}{
		{"application/pdf", "pdf"},
		{"image/png", "image"},
		{"video/mp4", "video"},
		{"application/octet-stream", "binary"},
		{"unknown/thing", "others"},
	}

	for _, c := range cases {
		got := getFileCategory(c.ext)
		if got != c.want {
			t.Errorf("getFileCategory(%q) = %q, want %q", c.ext, got, c.want)
		}
	}
}

func TestGetFileExtension(t *testing.T) {
	tmp := t.TempDir()

	// Создаём файл PDF
	path := filepath.Join(tmp, "test.pdf")
	if err := os.WriteFile(path, []byte("%PDF-1.4\n..."), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := getFileExtension(path)
	if err != nil {
		t.Fatal(err)
	}

	want := "application/pdf"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

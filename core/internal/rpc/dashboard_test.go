package rpc

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallDashboardReplacesExistingFiles(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "dashboard")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "old.html"), []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}

	archive := filepath.Join(root, "dashboard.zip")
	file, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	entry, err := writer.Create("dashboard-main/index.html")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write([]byte("dashboard")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	if err := installDashboard(archive, target); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(target, "index.html"))
	if err != nil || string(content) != "dashboard" {
		t.Fatalf("dashboard content = %q, error = %v", content, err)
	}
	if _, err := os.Stat(filepath.Join(target, "old.html")); !os.IsNotExist(err) {
		t.Fatalf("old dashboard file still exists or could not be checked: %v", err)
	}
}

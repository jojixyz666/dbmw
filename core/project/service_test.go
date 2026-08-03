package project_test

import (
	"os"
	"path/filepath"
	"testing"
	"dbmw/core/project"
)

func TestProjectService(t *testing.T) {
	svc := project.NewService()
	tempDir := t.TempDir()

	t.Run("Detect Laravel Project", func(t *testing.T) {
		laravelDir := filepath.Join(tempDir, "laravel_app")
		if err := os.MkdirAll(laravelDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(laravelDir, "artisan"), []byte("<?php"), 0644); err != nil {
			t.Fatal(err)
		}
		envContent := "DB_CONNECTION=mysql\nDB_HOST=127.0.0.1\nDB_PORT=3306\nDB_DATABASE=blog_db\nDB_USERNAME=root\n"
		if err := os.WriteFile(filepath.Join(laravelDir, ".env"), []byte(envContent), 0644); err != nil {
			t.Fatal(err)
		}

		info, err := svc.DetectFramework(laravelDir)
		if err != nil {
			t.Fatalf("detection failed: %v", err)
		}
		if info.Name != "Laravel" {
			t.Errorf("expected Laravel, got %s", info.Name)
		}
		if info.SuggestedDriver != "mysql" || info.SuggestedDB != "blog_db" {
			t.Errorf("expected mysql and blog_db, got %s and %s", info.SuggestedDriver, info.SuggestedDB)
		}
	})

	t.Run("Generate dbmw.yml", func(t *testing.T) {
		appDir := filepath.Join(tempDir, "sample_app")
		if err := os.MkdirAll(appDir, 0755); err != nil {
			t.Fatal(err)
		}
		cfg := project.DBMWConfig{
			Version:           "1",
			ProjectName:       "Sample App",
			DefaultConnection: "local_pg",
			Connections: []project.ProjectConnInfo{
				{Name: "local_pg", Driver: "postgres", Host: "localhost", Port: 5432, Database: "sample"},
			},
		}

		createdPath, err := svc.GenerateConfig(appDir, cfg)
		if err != nil {
			t.Fatalf("config gen failed: %v", err)
		}
		if _, err := os.Stat(createdPath); os.IsNotExist(err) {
			t.Fatalf("generated file does not exist at %s", createdPath)
		}
	})
}

package project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"gopkg.in/yaml.v3"
)

// Service provides framework detection and dbmw.yml generation.
type Service struct{}

// NewService instantiates ProjectService.
func NewService() *Service {
	return &Service{}
}

// DetectFramework inspects a directory and returns framework and database hints.
func (s *Service) DetectFramework(dirPath string) (*FrameworkInfo, error) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, err
	}

	info := &FrameworkInfo{
		Name:          "Unknown",
		Language:      "Generic",
		DetectedFiles: []string{},
	}

	envVars := s.parseEnvFile(filepath.Join(absPath, ".env"))
	if len(envVars) == 0 {
		envVars = s.parseEnvFile(filepath.Join(absPath, ".env.local"))
	}

	// 1. Laravel
	if s.fileExists(filepath.Join(absPath, "artisan")) {
		info.Name = "Laravel"
		info.Language = "PHP"
		info.DetectedFiles = append(info.DetectedFiles, "artisan")
		info.SuggestedDriver = envVars["DB_CONNECTION"]
		info.SuggestedHost = envVars["DB_HOST"]
		info.SuggestedUser = envVars["DB_USERNAME"]
		info.SuggestedDB = envVars["DB_DATABASE"]
		if p, err := strconv.Atoi(envVars["DB_PORT"]); err == nil {
			info.SuggestedPort = p
		}
		if info.SuggestedDriver == "" {
			info.SuggestedDriver = "mysql"
		}
		return info, nil
	}

	// 2. Prisma / Node
	if s.fileExists(filepath.Join(absPath, "prisma", "schema.prisma")) {
		info.Name = "Prisma"
		info.Language = "TypeScript/Node"
		info.DetectedFiles = append(info.DetectedFiles, "prisma/schema.prisma")
		driver, host, port, user, db := s.parseDatabaseURL(envVars["DATABASE_URL"])
		info.SuggestedDriver = driver
		info.SuggestedHost = host
		info.SuggestedPort = port
		info.SuggestedUser = user
		info.SuggestedDB = db
		return info, nil
	}

	// 3. Rails
	if s.fileExists(filepath.Join(absPath, "config", "database.yml")) {
		info.Name = "Ruby on Rails"
		info.Language = "Ruby"
		info.DetectedFiles = append(info.DetectedFiles, "config/database.yml")
		info.SuggestedDriver = "postgres"
		return info, nil
	}

	// 4. Django
	if s.fileExists(filepath.Join(absPath, "manage.py")) {
		info.Name = "Django"
		info.Language = "Python"
		info.DetectedFiles = append(info.DetectedFiles, "manage.py")
		info.SuggestedDriver = "postgres"
		return info, nil
	}

	// 5. Go Project
	if s.fileExists(filepath.Join(absPath, "go.mod")) {
		info.Name = "Go Application"
		info.Language = "Go"
		info.DetectedFiles = append(info.DetectedFiles, "go.mod")
		if dbURL := envVars["DATABASE_URL"]; dbURL != "" {
			driver, host, port, user, db := s.parseDatabaseURL(dbURL)
			info.SuggestedDriver = driver
			info.SuggestedHost = host
			info.SuggestedPort = port
			info.SuggestedUser = user
			info.SuggestedDB = db
		}
		return info, nil
	}

	// 6. Generic Node with .env
	if s.fileExists(filepath.Join(absPath, "package.json")) {
		info.Name = "Node.js Application"
		info.Language = "JavaScript"
		info.DetectedFiles = append(info.DetectedFiles, "package.json")
		if dbURL := envVars["DATABASE_URL"]; dbURL != "" {
			driver, host, port, user, db := s.parseDatabaseURL(dbURL)
			info.SuggestedDriver = driver
			info.SuggestedHost = host
			info.SuggestedPort = port
			info.SuggestedUser = user
			info.SuggestedDB = db
		}
		return info, nil
	}

	return info, nil
}

// GenerateConfig creates a dbmw.yml in the target folder.
func (s *Service) GenerateConfig(dirPath string, cfg DBMWConfig) (string, error) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return "", err
	}

	if cfg.Version == "" {
		cfg.Version = "1"
	}
	if cfg.ProjectName == "" {
		cfg.ProjectName = filepath.Base(absPath)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to encode yaml: %w", err)
	}

	targetPath := filepath.Join(absPath, "dbmw.yml")
	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write dbmw.yml: %w", err)
	}

	return targetPath, nil
}

func (s *Service) fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (s *Service) parseEnvFile(filePath string) map[string]string {
	res := make(map[string]string)
	f, err := os.Open(filePath)
	if err != nil {
		return res
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			res[k] = v
		}
	}
	return res
}

func (s *Service) parseDatabaseURL(rawURL string) (driver, host string, port int, user, database string) {
	if rawURL == "" {
		return "", "", 0, "", ""
	}
	// e.g. postgresql://user:pass@localhost:5432/dbname
	// e.g. mysql://root:pass@tcp(127.0.0.1:3306)/dbname
	// e.g. file:./dev.db (sqlite)
	if strings.HasPrefix(rawURL, "file:") || strings.HasSuffix(rawURL, ".db") {
		return "sqlite", "", 0, "", rawURL
	}
	if strings.HasPrefix(rawURL, "postgres://") || strings.HasPrefix(rawURL, "postgresql://") {
		driver = "postgres"
		port = 5432
	} else if strings.HasPrefix(rawURL, "mysql://") {
		driver = "mysql"
		port = 3306
	}

	return driver, "localhost", port, "", ""
}

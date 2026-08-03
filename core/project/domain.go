package project

// FrameworkInfo holds auto-detected framework information.
type FrameworkInfo struct {
	Name            string   `json:"name"`
	Language        string   `json:"language"`
	DetectedFiles   []string `json:"detectedFiles"`
	SuggestedDriver string   `json:"suggestedDriver"`
	SuggestedHost   string   `json:"suggestedHost"`
	SuggestedPort   int      `json:"suggestedPort"`
	SuggestedUser   string   `json:"suggestedUser"`
	SuggestedDB     string   `json:"suggestedDb"`
	ConfigPath      string   `json:"configPath,omitempty"`
}

// DBMWConfig represents the project-level dbmw.yml configuration.
type DBMWConfig struct {
	Version           string            `json:"version" yaml:"version"`
	ProjectName       string            `json:"projectName" yaml:"project_name"`
	DefaultConnection string            `json:"defaultConnection,omitempty" yaml:"default_connection,omitempty"`
	Connections       []ProjectConnInfo `json:"connections,omitempty" yaml:"connections,omitempty"`
}

type ProjectConnInfo struct {
	Name     string `json:"name" yaml:"name"`
	Driver   string `json:"driver" yaml:"driver"`
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
	User     string `json:"user,omitempty" yaml:"user,omitempty"`
	Database string `json:"database,omitempty" yaml:"database,omitempty"`
	FilePath string `json:"filePath,omitempty" yaml:"file_path,omitempty"`
}

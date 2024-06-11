package cmd

// Config represents the structure of the configuration file
type Config struct {
	AppName        string             `yaml:"app_name"`
	Port           int                `yaml:"port"`
	Version        string             `yaml:"version"`
	Database       DatabaseConfig     `yaml:"database"`
	Features       []FeatureConfig    `yaml:"features"`
	NestedObject   NestedObjectConfig `yaml:"nested_object"`
	ArrayOfStrings []string           `yaml:"array_of_strings"`
	ArrayOfObjects []ObjectConfig     `yaml:"array_of_objects"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type FeatureConfig struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
}

type NestedObjectConfig struct {
	Level1 Level1Config `yaml:"level1"`
}

type Level1Config struct {
	Level2 Level2Config `yaml:"level2"`
}

type Level2Config struct {
	Key string `yaml:"key"`
}

type ObjectConfig struct {
	ID   int    `yaml:"id"`
	Name string `yaml:"name"`
}

package data

type ServerConfig struct {
	Type string `yaml:"type"`
	Port int    `yaml:"port"`
}

type DbConfig struct {
	Type     string         `yaml:"type"`
	Postgres PostgresConfig `yaml:"postgres,omitempty"`
	SQLite   SQLiteConfig   `yaml:"sqlite,omitempty"`
	MySQL    MySQLConfig    `yaml:"mysql,omitempty"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type SQLiteConfig struct {
	File string `yaml:"file"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

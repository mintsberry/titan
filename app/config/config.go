package config

type Config struct {
	BaseDir  string   `json:"base_dir" env:"DIR" yaml:"baseDir"`
	Database Database `json:"database"`
}

type Database struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	User     string `json:"user" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASS"`
	Name     string `json:"name" env:"DB_NAME"`
	DBFile   string `json:"dbFile" env:"DB_FILE"`
}

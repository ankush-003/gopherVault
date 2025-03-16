package db

type DBConfig struct {
	Path string
}

func NewDBConfig(path string) *DBConfig {
	return &DBConfig{Path: path}
}

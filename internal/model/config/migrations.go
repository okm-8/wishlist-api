package config

type MigrationsConfig struct {
	dirPath string
}

func NewMigrationsConfig(dirPath string) *MigrationsConfig {
	return &MigrationsConfig{
		dirPath: dirPath,
	}
}

func (config *MigrationsConfig) DirPath() string {
	return config.dirPath
}

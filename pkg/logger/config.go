package logger

type LoggerConfig struct {
	File    LoggerFileConfig    `mapstructure:"file"`
	Console LoggerConsoleConfig `mapstructure:"console"`
}

type LoggerFileConfig struct {
	Level      string `mapstructure:"level"`
	Enabled    bool   `mapstructure:"enabled"`
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"`
	Compress   bool   `mapstructure:"compress"`
}

type LoggerConsoleConfig struct {
	Level   string `mapstructure:"level"`
	Enabled bool   `mapstructure:"enabled"`
}

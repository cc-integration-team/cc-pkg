package logger

type LoggerConfig struct {
	CallerDebug bool                `mapstructure:"callerDebug"`
	File        LoggerFileConfig    `mapstructure:"file"`
	Console     LoggerConsoleConfig `mapstructure:"console"`
}

type LoggerFileConfig struct {
	Level      string `mapstructure:"level"`
	Enabled    bool   `mapstructure:"enabled"`
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"`
	Compress   bool   `mapstructure:"compress"`
	Pretty     bool   `mapstructure:"pretty"`
}

type LoggerConsoleConfig struct {
	Level   string `mapstructure:"level"`
	Enabled bool   `mapstructure:"enabled"`
	Pretty  bool   `mapstructure:"pretty"`
}

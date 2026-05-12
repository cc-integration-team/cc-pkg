package logger

type LoggerConfig struct {
	Service string              `mapstructure:"service"`
	Caller  bool                `mapstructure:"caller"`
	Masking MaskingConfig       `mapstructure:"masking"`
	File    LoggerFileConfig    `mapstructure:"file"`
	Console LoggerConsoleConfig `mapstructure:"console"`
}

// MaskingConfig controls which log fields are masked before writing to output.
// Only applies to log levels >= Info; Debug logs are never masked.
type MaskingConfig struct {
	Enabled      bool                `mapstructure:"enabled"`
	Fields       []string            `mapstructure:"fields"`       // top-level JSON field names
	NestedFields map[string][]string `mapstructure:"nestedFields"` // parent key → child field names
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

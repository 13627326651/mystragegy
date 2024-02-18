package logger

import (
	"errors"
	"time"

	"go.uber.org/zap/zapcore"
)

type Format string

const (
	FormatPlain Format = "plain"
	FormatJSON  Format = "json"
)

type Config struct {
	EnableCaller    bool          `yaml:"EnableCaller"`
	ExtraCallerSkip int           `yaml:"ExtraCallerSkip"`
	Console         ConsoleConfig `yaml:"Console"`
	File            FileConfig    `yaml:"File"`
}

type ConsoleConfig struct {
	Base BaseConfig `yaml:"Base"`
}

type FileConfig struct {
	Base BaseConfig `yaml:"Base"`

	Format       Format       `yaml:"Format"`
	Name         string       `yaml:"Name"`
	MaxSize      int          `yaml:"MaxSize"`
	MaxBackups   int          `yaml:"MaxBackups"`
	MaxAge       int          `yaml:"MaxAge"`
	UseLocalTime bool         `yaml:"UseLocalTime"`
	Buffer       BufferConfig `yaml:"Buffer"`
}

type BaseConfig struct {
	Enable       bool          `yaml:"Enable"`
	WithFunc     bool          `yaml:"WithFunc"`
	BaseNameOnly bool          `yaml:"BaseNameOnly"`
	Coloring     bool          `yaml:"Coloring"`
	Level        zapcore.Level `yaml:"Level"`
	TimeLayout   string        `yaml:"TimeLayout"`
}

type BufferConfig struct {
	Enable        bool          `yaml:"Enable"`
	Size          int           `yaml:"Size"`
	FlushInterval time.Duration `yaml:"FlushInterval"`
}

func DefaultConfig() Config {
	return Config{
		EnableCaller:    true,
		ExtraCallerSkip: 0,
		Console: ConsoleConfig{
			Base: BaseConfig{
				Enable:       true,
				WithFunc:     true,
				BaseNameOnly: false,
				Coloring:     true,
				Level:        zapcore.DebugLevel,
				TimeLayout:   "2006-01-02 15:04:05.000000",
			},
		},
		File: FileConfig{
			Base: BaseConfig{
				Enable:       false,
				WithFunc:     true,
				BaseNameOnly: false,
				Coloring:     false,
				Level:        zapcore.DebugLevel,
				TimeLayout:   "2006-01-02 15:04:05.000000",
			},
			Format:       FormatPlain,
			Name:         "./log.log",
			MaxSize:      10,
			MaxBackups:   20,
			MaxAge:       30,
			UseLocalTime: true,
			Buffer: BufferConfig{
				Enable:        false,
				Size:          16 * 1024,
				FlushInterval: time.Second,
			},
		},
	}
}

func (config *Config) Check() error {
	if config.ExtraCallerSkip < 0 {
		return errors.New("config: invalid ExtraCallerSkip")
	}

	consoleConfig := &config.Console
	if consoleConfig.Base.TimeLayout == "" {
		return errors.New("config.Console: empty TimeLayout")
	}

	fileConfig := &config.File
	if fileConfig.Base.TimeLayout == "" {
		return errors.New("config.File: empty TimeLayout")
	}
	if fileConfig.Format != FormatPlain && fileConfig.Format != FormatJSON {
		return errors.New("config.File: invalid Format")
	}
	if fileConfig.Name == "" {
		return errors.New("config.File: empty Name")
	}
	if fileConfig.MaxSize < 0 {
		return errors.New("config.File: invalid MaxSize")
	}
	if fileConfig.MaxBackups < 0 {
		return errors.New("config.File: invalid MaxBackups")
	}
	if fileConfig.MaxAge < 0 {
		return errors.New("config.File: invalid MaxAge")
	}
	if fileConfig.Buffer.Size < 0 {
		return errors.New("config.File: invalid Buffer Size")
	}
	if fileConfig.Buffer.FlushInterval < 0 {
		return errors.New("config.File: invalid Buffer FlushInterval")
	}

	return nil
}

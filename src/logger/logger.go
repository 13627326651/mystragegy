package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"tinyquant/src/util"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger() {
	var consolLevel zapcore.Level
	switch util.ConsoleLevel {
	case "debug":
		consolLevel = zapcore.DebugLevel
	case "info":
		consolLevel = zapcore.InfoLevel
	case "warn":
		consolLevel = zapcore.WarnLevel
	case "error":
		consolLevel = zapcore.ErrorLevel
	default:
		consolLevel = zapcore.Level(zapcore.DebugLevel)
	}

	var fileLevel zapcore.Level
	switch util.ConsoleLevel {
	case "debug":
		fileLevel = zapcore.DebugLevel
	case "info":
		fileLevel = zapcore.InfoLevel
	case "warn":
		fileLevel = zapcore.WarnLevel
	case "error":
		fileLevel = zapcore.ErrorLevel
	default:
		fileLevel = zapcore.Level(zapcore.DebugLevel)
	}

	config := Config{
		EnableCaller:    true,
		ExtraCallerSkip: 0,
		Console: ConsoleConfig{
			Base: BaseConfig{
				Enable:       util.Console,
				WithFunc:     true,
				BaseNameOnly: false,
				Coloring:     true,
				Level:        consolLevel,
				TimeLayout:   "2006-01-02 15:04:05.000000",
			},
		},
		File: FileConfig{
			Base: BaseConfig{
				Enable:       util.File,
				WithFunc:     true,
				BaseNameOnly: false,
				Coloring:     false,
				Level:        fileLevel,
				TimeLayout:   "2006-01-02 15:04:05.000000",
			},
			Format:       FormatPlain,
			Name:         util.Path,
			MaxSize:      200,
			MaxBackups:   200,
			MaxAge:       15,
			UseLocalTime: true,
			Buffer: BufferConfig{
				Enable:        true,
				Size:          64 * 1024,
				FlushInterval: 5 * time.Second,
			},
		},
	}
	logger, err := NewLogger(config)
	if err != nil {
		panic(err)
	}
	Logger = logger
}

func NewLogger(config Config) (*zap.Logger, error) {
	if err := config.Check(); err != nil {
		return nil, err
	}

	cores := make([]zapcore.Core, 0, 2)

	if config.Console.Base.Enable {
		cores = append(cores, newConsoleCore(&config.Console))
	}
	if config.File.Base.Enable {
		cores = append(cores, newFileCore(&config.File))
	}

	return zap.New(zapcore.NewTee(cores...), zap.WithCaller(config.EnableCaller),
		zap.AddCallerSkip(config.ExtraCallerSkip)), nil
}

func newConsoleCore(config *ConsoleConfig) zapcore.Core {
	baseCfg := &config.Base
	encoder := newEncoder(FormatPlain, baseCfg.TimeLayout,
		baseCfg.WithFunc, baseCfg.BaseNameOnly, baseCfg.Coloring)
	writeSyncer := zapcore.AddSync(os.Stderr)
	level := zap.NewAtomicLevelAt(config.Base.Level)

	return zapcore.NewCore(encoder, writeSyncer, level)
}

func newFileCore(config *FileConfig) zapcore.Core {
	baseCfg := &config.Base
	encoder := newEncoder(config.Format, baseCfg.TimeLayout,
		baseCfg.WithFunc, baseCfg.BaseNameOnly, baseCfg.Coloring)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.Name,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		LocalTime:  config.UseLocalTime,
	}

	writeSyncer := zapcore.AddSync(lumberJackLogger)
	if config.Buffer.Enable {
		writeSyncer = &zapcore.BufferedWriteSyncer{
			WS:            writeSyncer,
			Size:          config.Buffer.Size,
			FlushInterval: config.Buffer.FlushInterval,
		}
	}

	level := zap.NewAtomicLevelAt(config.Base.Level)

	return zapcore.NewCore(encoder, writeSyncer, level)
}

func newEncoder(format Format, timeLayout string, withFunc, baseNameOnly, coloring bool) zapcore.Encoder {
	levelEncoder := zapcore.CapitalLevelEncoder
	if coloring {
		levelEncoder = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   levelEncoder,
		EncodeTime: func(tm time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(tm.Format(timeLayout))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			path := caller.TrimmedPath()
			if baseNameOnly {
				path = filepath.Base(path)
			}
			if withFunc {
				enc.AppendString(fmt.Sprintf("%s:%s", path, extractFuncName(caller.Function)))
			} else {
				enc.AppendString(path)
			}
		},
		ConsoleSeparator: " ",
	}

	if format == FormatJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	} else {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func extractFuncName(full string) string {
	lastSlash := strings.LastIndexByte(full, '/')
	nextDot := strings.IndexByte(full[lastSlash+1:], '.')
	if nextDot < 0 {
		return "???"
	} else {
		return full[lastSlash+1+nextDot+1:]
	}
}

package logger

import "strings"

type Level int8

const (
	LevelDebug Level = iota - 1
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func GetLevelByString(lv string) Level {
	upLv := strings.ToUpper(lv)
	if upLv == LevelDebug.String() {
		return LevelDebug
	} else if upLv == LevelInfo.String() {
		return LevelInfo
	} else if upLv == LevelWarn.String() {
		return LevelWarn
	} else if upLv == LevelError.String() {
		return LevelError
	} else if upLv == LevelFatal.String() {
		return LevelFatal
	} else {
		return LevelError
	}
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

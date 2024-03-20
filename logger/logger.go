package logger

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/khan-lau/kutils/container/kstrings"
	"github.com/khan-lau/kutils/datetime"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

type AppLogFunc func(lvl Level, f string, args ...interface{})

func (l Level) String() string {
	switch l {
	case -1:
		return "DEBUG"
	case 0:
		return "INFO"
	case 1:
		return "WARNING"
	case 2:
		return "ERROR"
	case 3:
		return "DPANIC"
	case 4:
		return "PANIC"
	case 5:
		return "FATAL"
	}
	panic("invalid LogLevel")
}

type Logger struct {
	log *zap.Logger
}

// var log *zap.Logger
func LoggerInstanceWithoutConsole(filename string, logLevel int8) *Logger {
	return LoggerInstance(filename, logLevel, false, true)
}

func LoggerInstanceOnlyConsole(logLevel int8) *Logger {
	return LoggerInstance("", logLevel, true, false)
}

func LoggerInstance(filename string, logLevel int8, needConsole bool, needTerminalColor bool) *Logger {
	if zapcore.Level(logLevel) < zapcore.DebugLevel || zapcore.Level(logLevel) > zapcore.FatalLevel {
		logLevel = int8(zapcore.InfoLevel)
	}

	filename = strings.TrimSpace(filename)

	cfg := zap.NewProductionEncoderConfig()
	//cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg.EncodeTime = zapcore.TimeEncoderOfLayout(datetime.DATETIME_FORMATTER_Mill)
	if needTerminalColor {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	encoder := zapcore.NewConsoleEncoder(cfg)

	var writeSyncer zapcore.WriteSyncer
	syncers := make([]zapcore.WriteSyncer, 0, 10)

	if len(filename) > 0 {
		// fullname := path.Base(filename)                          // 获取不包含目录的文件名
		file_suffix := path.Ext(filename)                         // 获取文件扩展名
		filen_prefix := strings.TrimSuffix(filename, file_suffix) // 获取文件名称和路径, 不包含扩展名

		logFile, _ := rotatelogs.New(filen_prefix+".%Y%m%d%H%M"+file_suffix,
			rotatelogs.WithMaxAge(30*24*time.Hour),    // 最长保存30天
			rotatelogs.WithRotationTime(time.Hour*24)) // 24小时切割一次

		syncers = append(syncers, zapcore.AddSync(logFile))
	}

	if needConsole {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	if len(syncers) == 0 {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	writeSyncer = zapcore.NewMultiWriteSyncer(syncers...)

	core := zapcore.NewCore(encoder, //NewJSONEncoder
		writeSyncer,
		zapcore.Level(logLevel))

	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // AddCaller() 显示文件名与行号; zap.AddCallerSkip(1)打印的文件名与行号在调用栈往外跳一层

	return &Logger{log: log}
}

///////////////////////////////////////////////////////////////

func (logger *Logger) Debug(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().Debugf(template, args...)
	}
}

func (logger *Logger) Info(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().Infof(template, args...)
	}
}

func (logger *Logger) Warrn(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().Warnf(template, args...)
	}
}

func (logger *Logger) Error(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().Errorf(template, args...)
	}
}

func (logger *Logger) DPanic(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().DPanicf(template, args...)
	}
}
func (logger *Logger) Fatal(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().Errorf(template, args...)
	}
}

///////////////////////////////////////////////////////////////

func (logger *Logger) D(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().Debugf(kstrings.FormatString(template, args...))
	}
}

func (logger *Logger) I(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().Infof(kstrings.FormatString(template, args...))
	}
}

func (logger *Logger) W(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().Warnf(kstrings.FormatString(template, args...))
	}
}

func (logger *Logger) E(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().Errorf(kstrings.FormatString(template, args...))
	}
}

func (logger *Logger) DP(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().DPanicf(kstrings.FormatString(template, args...))
	}
}
func (logger *Logger) F(template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().Errorf(kstrings.FormatString(template, args...))
	}
}

///////////////////////////////////////////////////////////////

func (logger *Logger) KDebug(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Debugf(template, args...)
	}
}

func (logger *Logger) KInfo(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Infof(template, args...)
	}
}

func (logger *Logger) KWarrn(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Warnf(template, args...)
	}
}

func (logger *Logger) KError(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Errorf(template, args...)
	}
}

func (logger *Logger) KDPanic(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).DPanicf(template, args...)
	}
}

func (logger *Logger) KFatal(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Errorf(template, args...)
	}
}

///////////////////////////////////////////////////////////////

func (logger *Logger) KD(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Debugf(kstrings.FormatString(template, args...))
	}
}

func (logger *Logger) KI(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Infof(kstrings.FormatString(template, args...))
	}
}

func (logger *Logger) KW(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Warnf(kstrings.FormatString(template, args...))
	}
}

func (logger *Logger) KE(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Errorf(kstrings.FormatString(template, args...))
	}
}

func (logger *Logger) KDP(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).DPanicf(kstrings.FormatString(template, args...))
	}
}

func (logger *Logger) KF(skip int, template string, args ...interface{}) {
	if logger != nil {
		logger.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Errorf(kstrings.FormatString(template, args...))
	}
}

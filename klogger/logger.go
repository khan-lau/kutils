package klogger

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
type AppLogFuncWithTag func(lvl Level, tag string, f string, args ...interface{})

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
	return LoggerInstance(filename, logLevel, false, false)
}

func LoggerInstanceOnlyConsole(logLevel int8) *Logger {
	return LoggerInstance("", logLevel, true, true)
}

func LoggerInstance(filename string, logLevel int8, needConsole bool, needTerminalColor bool) *Logger {
	conf := NewConfigure().SetLogFile(filename).SetLevel(Level(logLevel)).ShowConsole(needConsole).IsColorful(needTerminalColor)
	return GetLoggerWithConfig(conf)
}

func GetLoggerWithConfig(conf *LoggerConfigure) *Logger {
	if Level(conf.Level) < DebugLevel || Level(conf.Level) > FatalLevel {
		conf.Level = InfoLevel
	}

	filename := strings.TrimSpace(conf.LogFile)

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout(datetime.DATETIME_FORMATTER_Mill)
	if conf.Colorful {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	encoder := zapcore.NewConsoleEncoder(cfg)

	// var multiSyncer zapcore.WriteSyncer
	syncers := make([]zapcore.WriteSyncer, 0, 10)

	if len(filename) > 0 {
		file_suffix := path.Ext(filename)                         // 获取文件扩展名
		filen_prefix := strings.TrimSuffix(filename, file_suffix) // 获取文件名称和路径, 不包含扩展名

		if conf.MaxAge == 0 {
			conf.MaxAge = 3 * 24 // 默认最长保存3天
		}
		if conf.RotationTime == 0 {
			conf.RotationTime = 24 // 默认24小时滚动一次
		}

		logFile, _ := rotatelogs.New(filen_prefix+".%Y%m%d%H%M"+file_suffix,
			rotatelogs.WithMaxAge(time.Duration(conf.MaxAge)*time.Hour),             // 最长保存30天
			rotatelogs.WithRotationTime(time.Duration(conf.RotationTime)*time.Hour)) // 24小时切割一次

		// logFile := &lumberjack.Logger{
		// 	Filename:   filen_prefix + file_suffix,
		// 	MaxSize:    1024,        // 最大保存10MB日志文件
		// 	MaxBackups: 50,          // 最多保存10个备份
		// 	MaxAge:     conf.MaxAge, // 最长保存30天
		// 	LocalTime:  true,        // 本地时间
		// 	Compress:   true,        // 是否压缩
		// }

		syncers = append(syncers, zapcore.AddSync(logFile))
	}

	if conf.ToConsole {
		// os.Stdout.Fd() == syscall.Stdin
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	if len(syncers) == 0 {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	// 1. 先合并所有的写入端
	multiSyncer := zapcore.NewMultiWriteSyncer(syncers...)

	// // 2. 使用官方推荐的 BufferedWriteSyncer 实现异步批量写入
	// // 这在 v1.26.0 中是稳定且标准的方式
	// bufferedWriter := &zapcore.BufferedWriteSyncer{
	// 	WS:            multiSyncer,
	// 	Size:          4 * 1024,        // 4KB 缓冲区
	// 	FlushInterval: 1 * time.Second, // 每秒强制刷盘
	// }

	core := zapcore.NewCore(encoder, //NewJSONEncoder
		// bufferedWriter, //
		multiSyncer,
		zapcore.Level(zapcore.Level(conf.Level)))

	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // AddCaller() 显示文件名与行号; zap.AddCallerSkip(1)打印的文件名与行号在调用栈往外跳一层

	return &Logger{log: log}
}

func (that *Logger) Sync() {
	if that != nil {
		that.log.Sync() // 清空缓冲区
	}
}

///////////////////////////////////////////////////////////////

func (that *Logger) Debug(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().Debugf(template, args...)
	}
}

func (that *Logger) Info(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().Infof(template, args...)
	}
}

func (that *Logger) Warrn(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().Warnf(template, args...)
	}
}

func (that *Logger) Error(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().Errorf(template, args...)
	}
}

func (that *Logger) DPanic(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().DPanicf(template, args...)
	}
}
func (that *Logger) Fatal(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().Errorf(template, args...)
	}
}

///////////////////////////////////////////////////////////////

func (that *Logger) D(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().Debugf(kstrings.FormatString(template, args...))
	}
}

func (that *Logger) I(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().Infof(kstrings.FormatString(template, args...))
	}
}

func (that *Logger) W(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().Warnf(kstrings.FormatString(template, args...))
	}
}

func (that *Logger) E(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().Errorf(kstrings.FormatString(template, args...))
	}
}

func (that *Logger) DP(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().DPanicf(kstrings.FormatString(template, args...))
	}
}
func (that *Logger) F(template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().Errorf(kstrings.FormatString(template, args...))
	}
}

///////////////////////////////////////////////////////////////

func (that *Logger) KDebug(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Debugf(template, args...)
	}
}

func (that *Logger) KInfo(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Infof(template, args...)
	}
}

func (that *Logger) KWarrn(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Warnf(template, args...)
	}
}

func (that *Logger) KError(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Errorf(template, args...)
	}
}

func (that *Logger) KDPanic(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).DPanicf(template, args...)
	}
}

func (that *Logger) KFatal(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Errorf(template, args...)
	}
}

///////////////////////////////////////////////////////////////

func (that *Logger) KD(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Debugf(kstrings.FormatString(template, args...))
	}
}

func (that *Logger) KI(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Infof(kstrings.FormatString(template, args...))
	}
}

func (that *Logger) KW(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Warnf(kstrings.FormatString(template, args...))
	}
}

func (that *Logger) KE(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Errorf(kstrings.FormatString(template, args...))
	}
}

func (that *Logger) KDP(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).DPanicf(kstrings.FormatString(template, args...))
	}
}

func (that *Logger) KF(skip int, template string, args ...interface{}) {
	if that != nil {
		that.log.Sugar().WithOptions(zap.AddCallerSkip(skip)).Errorf(kstrings.FormatString(template, args...))
	}
}

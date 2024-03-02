package logger

import (
	"os"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log *zap.Logger
}

// var log *zap.Logger
func LoggerInstanceWithoutConsole(filename string, logLevel int8) *Logger {
	return LoggerInstance(filename, logLevel, false)
}

func LoggerInstance(filename string, logLevel int8, needConsole bool) *Logger {
	if zapcore.Level(logLevel) < zapcore.DebugLevel || zapcore.Level(logLevel) > zapcore.FatalLevel {
		logLevel = int8(zapcore.InfoLevel)
	}

	// fullname := path.Base(filename)                          // 获取不包含目录的文件名
	file_suffix := path.Ext(filename)                         // 获取文件扩展名
	filen_prefix := strings.TrimSuffix(filename, file_suffix) // 获取文件名称和路径, 不包含扩展名

	logFile, _ := rotatelogs.New(filen_prefix+".%Y%m%d%H%M"+file_suffix,
		rotatelogs.WithMaxAge(30*24*time.Hour),    // 最长保存30天
		rotatelogs.WithRotationTime(time.Hour*24)) // 24小时切割一次

	cfg := zap.NewProductionEncoderConfig()
	//cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(cfg)

	var writeSyncer zapcore.WriteSyncer
	if needConsole {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(logFile), zapcore.AddSync(os.Stdout))
	} else {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(logFile))
	}

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

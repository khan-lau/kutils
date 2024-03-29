package logger

// `
// // 日志等级,
// // DebugLevel Level = -1
// // // InfoLevel is the default logging priority.
// // InfoLevel = 0
// // // WarnLevel logs are more important than Info, but don't need individual
// // // human review.
// // WarnLevel = 1
// // // ErrorLevel logs are high-priority. If an application is running smoothly,
// // // it shouldn't generate any error-level logs.
// // ErrorLevel = 2
// // // DPanicLevel logs are particularly important errors. In development the
// // // logger panics after writing the message.
// // DPanicLevel = 3
// // // PanicLevel logs a message, then panics.
// // PanicLevel = 4
// // // FatalLevel logs a message, then calls os.Exit(1).
// // FatalLevel = 5
// logLevel: 0,

// // 是否需要彩色
// color : false,
// maxAge : 720,       // 日志最长保留时间, 单位 小时
// rotationTime : 24,  // 日志滚动周期, 单位 小时, 24小时滚动一个文件

// console: false,
// // 日志存放目录, 日志按天滚动, 最长保留30天, 默认值为 'logs'
// logFile: 'logs',
// `

type LoggerConfigure struct {
	Level        Level  `json:"logLevel"`     // 日志等级
	Colorful     bool   `json:"colorful"`     // 是否需要彩色
	MaxAge       int    `json:"maxAge"`       // 日志最长保留时间, 单位 小时
	RotationTime int    `json:"rotationTime"` // 日志滚动周期, 单位 小时, 24小时滚动一个文件
	ToConsole    bool   `json:"console"`      // 是否输出到控制台
	LogFile      string `json:"logFile"`      // 输出到文件, 如果文件名为空, 则不输出到文件
}

func NewConfigure() *LoggerConfigure {
	return &LoggerConfigure{
		Level:        0,
		Colorful:     false,
		MaxAge:       720,
		RotationTime: 24,
		ToConsole:    false,
		LogFile:      "",
	}
}

func (that *LoggerConfigure) SetLevel(level Level) *LoggerConfigure {
	that.Level = level
	return that
}

func (that *LoggerConfigure) IsColorful(need bool) *LoggerConfigure {
	that.Colorful = need
	return that
}

func (that *LoggerConfigure) SetMaxAge(age int) *LoggerConfigure {
	that.MaxAge = age
	return that
}

func (that *LoggerConfigure) SetRotationTime(time int) *LoggerConfigure {
	that.RotationTime = time
	return that
}

func (that *LoggerConfigure) ShowConsole(flag bool) *LoggerConfigure {
	that.ToConsole = flag
	return that
}

func (that *LoggerConfigure) SetLogFile(file string) *LoggerConfigure {
	that.LogFile = file
	return that
}

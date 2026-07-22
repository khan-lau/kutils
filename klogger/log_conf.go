package klogger

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
	Level         Level  `json:"logLevel"` // 日志等级
	Colorful      bool   `json:"colorful"` // 是否需要彩色
	Async         bool   `json:"async"`    // 是否异步输出日志, 默认同步输出
	flushInterval int64  // 强制刷盘周期, 单位 毫秒
	bufferSize    int64  // 缓冲区大小, 单位 Byte
	MaxAge        int    `json:"maxAge"`       // 日志最长保留时间, 单位 小时, MaxAge 与 MaxCount 必须只能设置一个
	MaxSize       int64  `json:"maxSize"`      // 单文件最大滚动大小, 单位 byte, 超过后强制滚动
	MaxCount      uint   `json:"maxCount"`     // 最多保留的备份文件数量, 默认50个, MaxAge 与 MaxCount 必须只能设置一个
	RotationTime  int    `json:"rotationTime"` // 日志滚动周期, 单位 小时, 24小时滚动一个文件
	ToConsole     bool   `json:"console"`      // 是否输出到控制台
	LogFile       string `json:"logFile"`      // 输出到文件, 如果文件名为空, 则不输出到文件
}

func NewConfigure() *LoggerConfigure {
	return &LoggerConfigure{
		Level:        0,
		Colorful:     false,
		Async:        false,
		MaxAge:       720,                     // 日志最长保留时间, 单位 小时
		MaxCount:     0,                       // 最多保留的备份文件数量, 默认50个
		MaxSize:      10 * 1024 * 1024 * 1024, // 默认最大文件大小 10G
		RotationTime: 24,                      // 日志滚动周期, 单位 小时, 24小时滚动一个文件
		ToConsole:    false,                   // 是否输出到控制台
		LogFile:      "",
	}
}

func (that *LoggerConfigure) SetLevel(level Level) *LoggerConfigure {
	that.Level = level
	return that
}

// 设置最多保留的备份文件数量, 默认50个, MaxAge 与 MaxCount 必须只能设置一个
func (that *LoggerConfigure) SetMaxCount(count uint) *LoggerConfigure {
	that.MaxCount = count
	return that
}

func (that *LoggerConfigure) IsColorful(need bool) *LoggerConfigure {
	that.Colorful = need
	return that
}

func (that *LoggerConfigure) SetAsync(async bool, flushInterval, bufferSize int64) *LoggerConfigure {
	that.Async = async
	that.flushInterval = flushInterval
	that.bufferSize = bufferSize
	return that
}

// 设置日志最长保留时间, 单位 小时, MaxAge 与 MaxCount 必须只能设置一个
func (that *LoggerConfigure) SetMaxAge(age int) *LoggerConfigure {
	that.MaxAge = age
	return that
}

// 设置单文件最大滚动大小, 单位 byte, 超过后强制滚动, 默认10G
func (that *LoggerConfigure) SetMaxSize(size int64) *LoggerConfigure {
	that.MaxSize = size
	return that
}

// 设置日志滚动时间, 单位 小时, 默认720小时(30天)
func (that *LoggerConfigure) SetRotationTime(time int) *LoggerConfigure {
	that.RotationTime = time
	return that
}

// 设置是否输出到控制台
func (that *LoggerConfigure) ShowConsole(flag bool) *LoggerConfigure {
	that.ToConsole = flag
	return that
}

// 设置日志输出文件, 如果文件名为空, 则不输出到文件
func (that *LoggerConfigure) SetLogFile(file string) *LoggerConfigure {
	that.LogFile = file
	return that
}

func (that *LoggerConfigure) FlushInterval() int64 {
	return that.flushInterval
}

func (that *LoggerConfigure) BufferSize() int64 {
	return that.bufferSize
}

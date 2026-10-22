package log

import "sync"

var logInstance Logger
var onceLog sync.Once

func GetLogger() Logger {
	onceLog.Do(func() {
		conf := NewDefaultOption()
		conf.Level = "debug"
		logInstance, _ = NewLogger(conf)
	})
	return logInstance
}

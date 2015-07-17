package suffuse

import "log"

type LogLevel uint8
const (
  LogError LogLevel = iota
  LogWarn
  LogInfo
  LogDebug
  LogTrace
)

var sfsLogFlags = 0 // log.LstdFlags would mean date and time
var sfsLogLevel = LogWarn
var sfsLogger   = log.New(OsStderr(), "", sfsLogFlags)

func logAt(level LogLevel, msg string, ctx ...interface{}) {
  if sfsLogLevel >= level {
    sfsLogger.Printf(msg + "\n", ctx...)
  }
}
func info(msg string, ctx ...interface{} ) { logAt(LogInfo, msg, ctx...)  }
func trace(msg string, ctx ...interface{} ) { logAt(LogTrace, msg, ctx...) }

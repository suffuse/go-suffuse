package suffuse

import (
  "os"
  lg "gopkg.in/inconshreveable/log15.v2"
)

const defaultLogLevel = lg.LvlWarn
var sfsLogger = level(defaultLogLevel)

func logI(msg string, ctx ...interface{}) { sfsLogger.Info(msg, ctx...) }
func logD(msg string, ctx ...interface{}) { sfsLogger.Debug(msg, ctx...) }
func logW(msg string, ctx ...interface{}) { sfsLogger.Warn(msg, ctx...) }
func logE(msg string, ctx ...interface{}) { sfsLogger.Error(msg, ctx...) }
func logC(msg string, ctx ...interface{}) { sfsLogger.Crit(msg, ctx...) }

func level(max lg.Lvl) lg.Logger {
  l := lg.New()
  handler := lg.StreamHandler(os.Stdout, lg.LogfmtFormat())
  l.SetHandler(lg.LvlFilterHandler(max, handler))
  return l
}

func SetLogLevel(max lg.Lvl) { sfsLogger = level(max) }

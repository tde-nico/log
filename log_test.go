package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	LOG.SetLogLevel("debug")

	LOG.Debug("Debug message")
	LOG.Info("Info message")
	LOG.Notice("notice message")
	LOG.Warn("Warn message")
	LOG.Error("Error message")
	LOG.Critical("Critical message")

	sub := LOG.WithPrefix("sub")
	sub.SetLevel(InfoLevel)
	sub.SetTimeFormat("2006-01-02 15:04:05")

	sub.Debug("Debug message")
	sub.Info("Info message")
	sub.Notice("notice message")
	sub.Warn("Warn message")
	sub.Error("Error message")
	sub.Critical("Critical message")
}

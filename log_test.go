package log

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestLog(t *testing.T) {
	LOG.SetLogLevel("debug")

	LOG.SetLogFile("test.log")
	defer LOG.CloseLogFile()

	styles := LOG.DefaultStyles()
	styles.Keys["test"] = lipgloss.NewStyle().Foreground(lipgloss.Color("114"))
	styles.Values["test"] = lipgloss.NewStyle().Bold(true)
	LOG.SetStyles(styles)

	LOG.Debug("Debug message")
	LOG.Info("Info message", "test", 500)
	LOG.Notice("notice message")
	LOG.Warn("Warn message")
	LOG.Error("Error message")
	LOG.Critical("Critical message")

	sub := LOG.WithPrefix("sub")
	sub.SetLevel(InfoLevel)
	sub.SetTimeFormat("2006-01-02 15:04:05")

	sub.Debug("Debug message")
	sub.Info("Info message", "test", 500)
	sub.Notice("notice message")
	sub.Warn("Warn message")
	sub.Error("Error message")
	sub.Critical("Critical message")
}

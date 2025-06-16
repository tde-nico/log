package log

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestLog(t *testing.T) {
	SetLogLevel("debug")

	styles := DefaultStyles()
	styles.Keys["test"] = lipgloss.NewStyle().Foreground(lipgloss.Color("114"))
	styles.Values["test"] = lipgloss.NewStyle().Bold(true)
	SetStyles(styles)

	SetLogFile("test.log")
	defer CloseLogFile()

	Debug("Debug message")
	Info("Info message", "test", 500)
	Notice("notice message")
	Warn("Warn message")
	Error("Error message")
	Critical("Critical message")

	sub := WithPrefix("sub")
	sub.SetLevel(InfoLevel)
	sub.SetTimeFormat("2006-01-02 15:04:05")

	sub.Debug("Debug message")
	sub.Info("Info message", "test", 500)
	sub.Notice("notice message")
	sub.Warn("Warn message")
	sub.Error("Error message")
	sub.Critical("Critical message")
	sub.Fatal("Death message")
}

package log

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const NOTICE = log.Level(2)
const CRITICAL = log.Level(10)

const FatalLevel = log.FatalLevel
const CriticalLevel = CRITICAL
const ErrorLevel = log.ErrorLevel
const WarnLevel = log.WarnLevel
const NoticeLevel = NOTICE
const InfoLevel = log.InfoLevel
const DebugLevel = log.DebugLevel

type Logger struct {
	*log.Logger
}

var LOG *Logger = &Logger{log.Default()}

var Fatal = LOG.Fatal
var Fatalf = LOG.Fatalf
var Critical = LOG.Critical
var Criticalf = LOG.Criticalf
var Error = LOG.Error
var Errorf = LOG.Errorf
var Warn = LOG.Warn
var Warnf = LOG.Warnf
var Notice = LOG.Notice
var Noticef = LOG.Noticef
var Info = LOG.Info
var Infof = LOG.Infof
var Debug = LOG.Debug
var Debugf = LOG.Debugf
var Print = LOG.Print
var Printf = LOG.Printf
var Log = LOG.Log
var Logf = LOG.Logf

var Helper = LOG.Helper
var GetLevel = LOG.GetLevel
var With = LOG.With
var WithPrefix = LOG.WithPrefix
var GetPrefix = LOG.GetPrefix
var SetPrefix = LOG.SetPrefix
var SetLevel = LOG.SetLevel
var SetLogLevel = LOG.SetLogLevel
var SetTimeFormat = LOG.SetTimeFormat
var SetStyles = LOG.SetStyles
var DefaultStyles = LOG.DefaultStyles
var CloseLogFile = LOG.CloseLogFile
var SetLogFile = LOG.SetLogFile

var defaultStyles *log.Styles
var file *os.File

func init() {
	os.Setenv("TERM", "xterm-256color")
	os.Setenv("CLICOLOR_FORCE", "1")

	defaultStyles = log.DefaultStyles()
	defaultStyles.Levels[CRITICAL] = lipgloss.NewStyle().
		SetString("CRITICAL").
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.Color("201"))
	defaultStyles.Levels[NOTICE] = lipgloss.NewStyle().
		SetString("NOTICE").
		Bold(true).
		MaxWidth(4).
		Foreground(lipgloss.Color("40"))

	defaultStyles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	defaultStyles.Values["err"] = lipgloss.NewStyle().Bold(true)

	LOG.SetStyles(defaultStyles)
	LOG.SetTimeFormat("15:04:05")
}

func (l *Logger) SetLogFile(fname string) {
	if fname == "" {
		return
	}

	var err error
	file, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		l.Errorf("Failed to open log file '%s': %v", fname, err)
	}

	l.SetOutput(io.MultiWriter(file, os.Stdout))
}

func (l *Logger) CloseLogFile() {
	if file != nil {
		if err := file.Close(); err != nil {
			l.Errorf("Failed to close log file: %v", err)
		}
		file = nil
	}
}

func (l *Logger) DefaultStyles() *log.Styles {
	return defaultStyles
}

func (l *Logger) SetStyles(styles *log.Styles) {
	l.Logger.SetStyles(styles)
}

func (l *Logger) SetTimeFormat(format string) {
	l.Logger.SetTimeFormat(format)
}

func (l *Logger) SetLevel(lvl log.Level) {
	l.Logger.SetLevel(lvl)
	l.SetReportCaller(lvl <= DebugLevel)
}

func (l *Logger) SetLogLevel(level string) {
	var lvl log.Level
	var err error

	level = strings.ToLower(level)
	switch level {
	case "notice":
		lvl = NOTICE
	case "critical":
		lvl = CRITICAL
	default:
		lvl, err = log.ParseLevel(level)
		if err != nil {
			l.Errorf("Invalid log level '%s': %v", level, err)
			return
		}
	}
	l.SetLevel(lvl)
}

func (l *Logger) With(keyvals ...interface{}) *Logger {
	return &Logger{l.Logger.With(keyvals...)}
}
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{l.Logger.WithPrefix(prefix)}
}

func (l *Logger) Critical(msg interface{}, keyvals ...interface{}) {
	l.Helper()
	l.Log(CRITICAL, msg, keyvals...)
}

func (l *Logger) Criticalf(format string, keyvals ...interface{}) {
	l.Helper()
	l.Logf(CRITICAL, format, keyvals...)
}

func (l *Logger) Notice(msg interface{}, keyvals ...interface{}) {
	l.Helper()
	l.Log(NOTICE, msg, keyvals...)
}

func (l *Logger) Noticef(format string, keyvals ...interface{}) {
	l.Helper()
	l.Logf(NOTICE, format, keyvals...)
}

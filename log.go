package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Logger struct {
	logger     *log.Logger
	fileLogger *log.Logger
	file       *os.File
	children   []*Logger
}

const NOTICE = log.Level(2)
const CRITICAL = log.Level(10)

const FatalLevel = log.FatalLevel
const CriticalLevel = CRITICAL
const ErrorLevel = log.ErrorLevel
const WarnLevel = log.WarnLevel
const NoticeLevel = NOTICE
const InfoLevel = log.InfoLevel
const DebugLevel = log.DebugLevel

var LOG *Logger = &Logger{log.Default(), nil, nil, []*Logger{}}

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

var With = LOG.With
var WithPrefix = LOG.WithPrefix
var GetPrefix = LOG.GetPrefix
var SetPrefix = LOG.SetPrefix
var GetLevel = LOG.GetLevel
var SetLevel = LOG.SetLevel
var SetLogLevel = LOG.SetLogLevel
var SetTimeFormat = LOG.SetTimeFormat
var SetStyles = LOG.SetStyles
var DefaultStyles = LOG.DefaultStyles
var CloseLogFile = LOG.CloseLogFile
var SetLogFile = LOG.SetLogFile

var defaultStyles *log.Styles

func init() {
	os.Setenv("TERM", "xterm-256color")
	// os.Setenv("CLICOLOR_FORCE", "1")

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
	l.CloseLogFile()

	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		l.Errorf("Failed to open log file '%s': %v", fname, err)
		return
	}

	l.file = file
	l.fileLogger = log.With()
	l.fileLogger.SetOutput(file)
}

func (l *Logger) CloseLogFile() {
	if l.file != nil {
		l.fileLogger = nil
		if err := l.file.Close(); err != nil {
			l.logger.Errorf("Failed to close log file: %v", err)
		}
		l.file = nil
	}
}

func (l *Logger) DefaultStyles() *log.Styles {
	return defaultStyles
}

func (l *Logger) SetStyles(styles *log.Styles) {
	l.logger.SetStyles(styles)
	if l.fileLogger != nil {
		l.fileLogger.SetStyles(styles)
	}
}

func (l *Logger) SetTimeFormat(format string) {
	l.logger.SetTimeFormat(format)
	if l.fileLogger != nil {
		l.fileLogger.SetTimeFormat(format)
	}
}

func (l *Logger) GetLevel() log.Level {
	return l.logger.GetLevel()
}

func (l *Logger) SetLevel(lvl log.Level) {
	l.logger.SetLevel(lvl)
	l.logger.SetReportCaller(lvl <= DebugLevel)
	if l.fileLogger != nil {
		l.fileLogger.SetLevel(lvl)
		l.fileLogger.SetReportCaller(lvl <= DebugLevel)
	}
	for _, child := range l.children {
		child.SetLevel(lvl)
	}
}

func (l *Logger) SetLogLevel(level string) {
	var lvl log.Level
	var err error

	level = strings.ToLower(level)
	switch level {
	case "notice":
		lvl = NoticeLevel
	case "critical":
		lvl = CriticalLevel
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
	var child *Logger
	if l.fileLogger != nil {
		child = &Logger{l.logger.With(keyvals...), l.fileLogger.With(keyvals...), nil, []*Logger{}}
	} else {
		child = &Logger{l.logger.With(keyvals...), nil, nil, []*Logger{}}
	}
	l.children = append(l.children, child)
	return child
}

func (l *Logger) WithPrefix(prefix string) *Logger {
	var child *Logger
	if l.fileLogger != nil {
		child = &Logger{l.logger.WithPrefix(prefix), l.fileLogger.WithPrefix(prefix), nil, []*Logger{}}
	} else {
		child = &Logger{l.logger.WithPrefix(prefix), nil, nil, []*Logger{}}
	}
	l.children = append(l.children, child)
	return child
}

func (l *Logger) GetPrefix() string {
	return l.logger.GetPrefix()
}

func (l *Logger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
	if l.fileLogger != nil {
		l.fileLogger.SetPrefix(prefix)
	}
}

func (l *Logger) Debug(msg interface{}, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Log(DebugLevel, msg, keyvals...)
	}
	l.logger.Helper()
	l.logger.Log(DebugLevel, msg, keyvals...)
}

func (l *Logger) Debugf(format string, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Logf(DebugLevel, format, keyvals...)
	}
	l.logger.Helper()
	l.logger.Logf(DebugLevel, format, keyvals...)
}

func (l *Logger) Info(msg interface{}, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Log(InfoLevel, msg, keyvals...)
	}
	l.logger.Helper()
	l.logger.Log(InfoLevel, msg, keyvals...)
}

func (l *Logger) Infof(format string, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Logf(InfoLevel, format, keyvals...)
	}
	l.logger.Helper()
	l.logger.Logf(InfoLevel, format, keyvals...)
}

func (l *Logger) Notice(msg interface{}, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Log(NoticeLevel, msg, keyvals...)
	}
	l.logger.Helper()
	l.logger.Log(NoticeLevel, msg, keyvals...)
}

func (l *Logger) Noticef(format string, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Logf(NoticeLevel, format, keyvals...)
	}
	l.logger.Helper()
	l.logger.Logf(NoticeLevel, format, keyvals...)
}

func (l *Logger) Warn(msg interface{}, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Log(WarnLevel, msg, keyvals...)
	}
	l.logger.Helper()
	l.logger.Log(WarnLevel, msg, keyvals...)
}

func (l *Logger) Warnf(format string, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Logf(WarnLevel, format, keyvals...)
	}
	l.logger.Helper()
	l.logger.Logf(WarnLevel, format, keyvals...)
}

func (l *Logger) Error(msg interface{}, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Log(ErrorLevel, msg, keyvals...)
	}
	l.logger.Helper()
	l.logger.Log(ErrorLevel, msg, keyvals...)
}

func (l *Logger) Errorf(format string, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Logf(ErrorLevel, format, keyvals...)
	}
	l.logger.Helper()
	l.logger.Logf(ErrorLevel, format, keyvals...)
}

func (l *Logger) Critical(msg interface{}, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Log(CriticalLevel, msg, keyvals...)
	}
	l.logger.Helper()
	l.logger.Log(CriticalLevel, msg, keyvals...)
}

func (l *Logger) Criticalf(format string, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Logf(CriticalLevel, format, keyvals...)
	}
	l.logger.Helper()
	l.logger.Logf(CriticalLevel, format, keyvals...)
}

func (l *Logger) Fatal(msg interface{}, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Log(CriticalLevel, "FATAL: "+fmt.Sprint(msg), keyvals...)
	}
	l.logger.Helper()
	l.logger.Log(FatalLevel, msg, keyvals...)
}

func (l *Logger) Fatalf(format string, keyvals ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Helper()
		l.fileLogger.Logf(CriticalLevel, "FATAL: "+format, keyvals...)
	}
	l.logger.Helper()
	l.logger.Logf(FatalLevel, format, keyvals...)
}

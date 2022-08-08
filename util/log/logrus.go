package log

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(DebugLevel)
}

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// Entry is the final or intermediate Logrus logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Debug, Info,
// Warn, Error, Fatal or Panic is called on it. These objects can be reused and
// passed around as much as you wish to avoid field duplication.
type Entry = logrus.Entry

// Fields type, used to pass to `WithFields`.
type Fields = logrus.Fields

// FieldMap allows customization of the key names for default fields.
type FieldMap = logrus.FieldMap

// Level type
type Level = logrus.Level

// Logger type
type Logger = logrus.Logger

// JSONFormatter formats logs into parsable json
type JSONFormatter struct {
	logrus.JSONFormatter
}

// TextFormatter formats logs into text
type TextFormatter struct {
	logrus.TextFormatter
}

// Formatter
// The Formatter interface is used to implement a custom Formatter. It takes an
// `Entry`. It exposes all the fields, including the default ones:
//
// * `entry.Data["msg"]`. The message passed from Info, Warn, Error ..
// * `entry.Data["time"]`. The timestamp.
// * `entry.Data["level"]. The level the entry was logged at.
//
// Any additional fields added with `WithField` or `WithFields` are also in
// `entry.Data`. Format is expected to return an array of bytes which are then
// logged to `logger.Out`.

// SetLevel ...
func SetLevel(level Level) {
	logrus.SetLevel(level)
}

// NewWithFields returns a logrus Entry with fields
func NewWithFields(fields Fields) *Entry {
	return logrus.WithFields(fields)
}

// NewEntry return an entry is the final or intermediate Logrus logging entry
func NewEntry(logger *Logger) *Entry {
	return logrus.NewEntry(logger)
}

// Exported from logrus
var (
	// Creates a new logger. Configuration should be set by changing `Formatter`,
	// `Out` and `Hooks` directly on the default logger instance. You can also just
	// instantiate your own:
	//
	//    var log = &Logger{
	//      Out: os.Stderr,
	//      Formatter: new(JSONFormatter),
	//      Level: logrus.DebugLevel,
	//    }
	//
	// It's recommended to make this a global instance called `log`.
	New = logrus.New
	// StandardLogger default logger
	StandardLogger = logrus.StandardLogger
	// SetOutput sets the standard logger output.
	SetOutput = logrus.SetOutput
	// SetFormatter sets the standard logger formatter.
	SetFormatter = logrus.SetFormatter
	// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
	WithError = logrus.WithError
	// WithField creates an entry from the standard logger and adds a field to
	// it. If you want multiple fields, use `WithFields`.
	//
	// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
	// or Panic on the Entry it returns.
	WithField = logrus.WithField
	// WithFields creates an entry from the standard logger and adds multiple
	// fields to it. This is simply a helper for `WithField`, invoking it
	// once for each field.
	//
	// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
	// or Panic on the Entry it returns.
	WithFields = logrus.WithFields

	// Debug logs a message at level Debug on the standard logger.
	Debug = logrus.Debug
	// Print logs a message at level Info on the standard logger.
	Print = logrus.Print
	// Info logs a message at level Info on the standard logger.
	Info = logrus.Info
	// Warn logs a message at level Warn on the standard logger.
	Warn = logrus.Warn
	// Warning logs a message at level Warn on the standard logger.
	Warning = logrus.Warning
	// Error logs a message at level Error on the standard logger.
	Error = logrus.Error
	// Panic logs a message at level Panic on the standard logger.
	Panic = logrus.Panic
	// Fatal logs a message at level Fatal on the standard logger.
	Fatal = logrus.Fatal

	// Debugf logs a message at level Debug on the standard logger.
	Debugf = logrus.Debugf
	// Printf logs a message at level Info on the standard logger.
	Printf = logrus.Printf
	// Infof logs a message at level Info on the standard logger.
	Infof = logrus.Infof
	// Warnf logs a message at level Warn on the standard logger.
	Warnf = logrus.Warnf
	// Warningf logs a message at level Warn on the standard logger.
	Warningf = logrus.Warningf
	// Errorf logs a message at level Error on the standard logger.
	Errorf = logrus.Errorf
	// Panicf logs a message at level Panic on the standard logger.
	Panicf = logrus.Panicf
	// Fatalf logs a message at level Fatal on the standard logger.
	Fatalf = logrus.Fatalf
	// Debugln logs a message at level Debug on the standard logger.
	Debugln = logrus.Debugln
	// Println logs a message at level Info on the standard logger.
	Println = logrus.Println
	// Infoln logs a message at level Info on the standard logger.
	Infoln = logrus.Infoln
	// Warnln logs a message at level Warn on the standard logger.
	Warnln = logrus.Warnln
	// Warningln logs a message at level Warn on the standard logger.
	Warningln = logrus.Warningln
	// Errorln logs a message at level Error on the standard logger.
	Errorln = logrus.Errorln
	// Panicln logs a message at level Panic on the standard logger.
	Panicln = logrus.Panicln
	// Fatalln logs a message at level Fatal on the standard logger.
	Fatalln = logrus.Fatalln
)

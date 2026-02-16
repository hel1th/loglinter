package loggers

var logMethods = map[string]struct{}{
	// slog methods
	"Debug": {},
	"Info":  {},
	"Warn":  {},
	"Error": {},

	// zap sugared methods
	"Debugf": {},
	"Infof":  {},
	"Warnf":  {},
	"Errorf": {},
	"Debugw": {},
	"Infow":  {},
	"Warnw":  {},
	"Errorw": {},

	// standard log
	"Print":   {},
	"Printf":  {},
	"Println": {},
	"Fatal":   {},
	"Fatalf":  {},
	"Fatalln": {},
	"Panic":   {},
	"Panicf":  {},
	"Panicln": {},
}

type LoggerType string

const (
	ZapLogger     LoggerType = "zap"
	LogLogger     LoggerType = "log"
	SlogLogger    LoggerType = "slog"
	UnknownLogger LoggerType = "unknown"
)

// a slice of logger types associated with their pkg links
var loggerChecks = []struct {
	pattern string
	logType LoggerType
}{
	{"log/slog", SlogLogger},
	{"go.uber.org/zap", ZapLogger},
	{"log.Logger", LogLogger},
}

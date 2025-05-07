package logger

import (
	"bytes"
	"context"
	"errors"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	slowThreshold = 200 * time.Millisecond
)

// Option -
func Option(log zerolog.Logger, config *logger.Config) *gorm.Config {
	return &gorm.Config{Logger: New(log, config)}
}

// New -
func New(log zerolog.Logger, config *logger.Config) logger.Interface {
	if config == nil {
		config = &logger.Config{
			SlowThreshold:        slowThreshold,
			LogLevel:             logger.Info,
			Colorful:             false,
			ParameterizedQueries: true,
		}
	}

	return &Logger{
		log:    log,
		config: config,
		isDev:  strings.EqualFold(os.Getenv("REGION_ID"), "dev"),
	}
}

// Logger -
type Logger struct {
	log    zerolog.Logger
	config *logger.Config
	isDev  bool
}

// LogMode -
func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.log = newLogger.log.Level(LogLevelMap[level])
	return &newLogger
}

// Error -
func (l *Logger) Error(ctx context.Context, msg string, opts ...any) {
	l.log.Error().Ctx(ctx).Msgf(msg, opts...)
}

// Warn -
func (l *Logger) Warn(ctx context.Context, msg string, opts ...any) {
	l.log.Warn().Ctx(ctx).Msgf(msg, opts...)
}

// Info -
func (l *Logger) Info(ctx context.Context, msg string, opts ...any) {
	l.log.Debug().Ctx(ctx).Msgf(msg, opts...)
}

// Trace -
func (l *Logger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	span := time.Since(begin)

	var stack []byte
	var slow bool
	var logFunc func(log zerolog.Logger) *zerolog.Event

	switch {
	case isError(err):
		logFunc = logFuncMethod(zerolog.ErrorLevel)
		stack = getStack(5)
	case l.config.SlowThreshold > 0 && span > l.config.SlowThreshold:
		logFunc = logFuncMethod(zerolog.WarnLevel)
		stack = getStack(5)
		slow = true
	case l.isDev:
		logFunc = logFuncMethod(l.log.GetLevel())
	}

	if logFunc == nil {
		return
	}

	event := logFunc(l.log)
	sql, rows := f()

	if l.isDev {
		event.Str("sql_query", sql)
		event.Int64("sql_rows", rows)
	}

	if len(stack) > 0 {
		event.Str("sql_stack", string(stack))
	}

	event.Dur("sql_duration", span)
	event.Str("sql_location", utils.FileWithLineNum())

	msg := "SQL execute - time_spent=%s, slow=%t"
	vals := []any{span, slow}

	if isError(err) {
		msg += ", repository.Error=%s"
		vals = append(vals, err)
	}

	event.Msgf(msg, vals...)
}

var numericPlaceholder = regexp.MustCompile(`\$(\d+)`)

// ParamsFilter -
func (l *Logger) ParamsFilter(ctx context.Context, sql string, params ...any) (string, []any) {
	if !l.isDev && l.config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

func isError(err error) bool {
	return err != nil && !errors.Is(err, gorm.ErrRecordNotFound)
}

func getStack(depth int) []byte {
	buf := debug.Stack()

	head := bytes.Index(buf, []byte("\ngithub.com/jcfug8/daylear/server/adapters/gorm/dialer/logger."))
	if head < 0 {
		return nil
	}

	buf = buf[head+1:]

	var tail int

	for i := 0; i < depth+1; i++ {
		n := bytes.Index(buf[tail:], []byte("\n\t"))
		if n < 0 {
			return nil
		}

		tail += n + 2

		n = bytes.Index(buf[tail:], []byte("\n"))
		if n < 0 {
			return nil
		}

		tail += n + 1

		if i == 0 {
			head = tail
		}

		if tail == len(buf) {
			break
		}
	}

	if head == tail || tail == 0 {
		return nil
	}

	return buf[head : tail-1]
}

// LevelMap -
var LogLevelMap = map[logger.LogLevel]zerolog.Level{
	logger.Silent: zerolog.Disabled,
	logger.Error:  zerolog.ErrorLevel,
	logger.Warn:   zerolog.WarnLevel,
	logger.Info:   zerolog.InfoLevel,
}

var logFuncMethod = func() func(zerolog.Level) func(zerolog.Logger) *zerolog.Event {
	var (
		Error = func(log zerolog.Logger) *zerolog.Event { return log.Error() }
		Warn  = func(log zerolog.Logger) *zerolog.Event { return log.Warn() }
		Debug = func(log zerolog.Logger) *zerolog.Event { return log.Debug() }
		Trace = func(log zerolog.Logger) *zerolog.Event { return log.Trace() }
	)

	Methods := map[zerolog.Level]func(zerolog.Logger) *zerolog.Event{
		zerolog.ErrorLevel: Error,
		zerolog.WarnLevel:  Warn,
		zerolog.DebugLevel: Debug,
		zerolog.TraceLevel: Trace,
	}

	return func(level zerolog.Level) func(log zerolog.Logger) *zerolog.Event {
		if Method, ok := Methods[level]; ok {
			return Method
		}

		return Debug
	}
}()

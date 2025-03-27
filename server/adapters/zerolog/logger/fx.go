package logger

import (
	"strings"

	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

// NewFxLogger -
func NewFxLogger(log zerolog.Logger) fxevent.Logger {
	return &FxLogger{
		log: log,
	}
}

// FxLogger -
type FxLogger struct {
	log zerolog.Logger
}

// LogEvent -
func (l *FxLogger) LogEvent(evt fxevent.Event) {
	switch e := evt.(type) {
	case *fxevent.OnStartExecuting:
		l.log.Debug().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msgf("OnStart hook executing function %s", e.FunctionName)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.log.Error().
				Err(e.Err).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msgf("OnStart hook failed for function %s", e.FunctionName)
		} else {
			l.log.Debug().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msgf("OnStart hook executed for function %s", e.FunctionName)
		}
	case *fxevent.OnStopExecuting:
		l.log.Debug().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msgf("OnStop hook executing function %s", e.FunctionName)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.log.Error().
				Err(e.Err).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msgf("OnStop hook failed for function %s", e.FunctionName)
		} else {
			l.log.Debug().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msgf("OnStop hook executed for function %s", e.FunctionName)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.log.Error().
				Err(e.Err).
				Str("type", e.TypeName).
				Strs("stacktrace", e.StackTrace).
				Strs("moduletrace", e.ModuleTrace).
				Str("module", e.ModuleName).
				Msgf("error encountered while applying options for module %s", e.ModuleName)
		} else {
			l.log.Debug().
				Str("type", e.TypeName).
				Strs("stacktrace", e.StackTrace).
				Strs("moduletrace", e.ModuleTrace).
				Str("module", e.ModuleName).
				Msgf("supplied type %s for module %s", e.TypeName, e.ModuleName)
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.log.Debug().
				Str("constructor", e.ConstructorName).
				Strs("stacktrace", e.StackTrace).
				Strs("moduletrace", e.ModuleTrace).
				Str("module", e.ModuleName).
				Str("type", rtype).
				Bool("private", e.Private).
				Msgf("provided type %s for module %s", rtype, e.ModuleName)
		}
		if e.Err != nil {
			l.log.Error().
				Err(e.Err).
				Strs("stacktrace", e.StackTrace).
				Strs("moduletrace", e.ModuleTrace).
				Str("module", e.ModuleName).
				Msgf("error encountered while providing module %s", e.ModuleName)
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			l.log.Debug().
				Strs("stacktrace", e.StackTrace).
				Strs("moduletrace", e.ModuleTrace).
				Str("module", e.ModuleName).
				Str("type", rtype).
				Msgf("replaced type %s for module %s", rtype, e.ModuleName)
		}
		if e.Err != nil {
			l.log.Error().
				Err(e.Err).
				Strs("stacktrace", e.StackTrace).
				Strs("moduletrace", e.ModuleTrace).
				Str("module", e.ModuleName).
				Msgf("error encountered while replacing module %s", e.ModuleName)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.log.Debug().
				Str("decorator", e.DecoratorName).
				Strs("stacktrace", e.StackTrace).
				Strs("moduletrace", e.ModuleTrace).
				Str("module", e.ModuleName).
				Str("type", rtype).
				Msgf("decorated type %s for module %s", rtype, e.ModuleName)
		}
		if e.Err != nil {
			l.log.Error().
				Err(e.Err).
				Strs("stacktrace", e.StackTrace).
				Strs("moduletrace", e.ModuleTrace).
				Str("module", e.ModuleName).
				Msgf("error encountered while decorating module %s", e.ModuleName)
		}
	case *fxevent.Run:
		if e.Err != nil {
			l.log.Error().
				Err(e.Err).
				Str("name", e.Name).
				Str("kind", e.Kind).
				Str("module", e.ModuleName).
				Msgf("error returned from run %s in module %s", e.Name, e.ModuleName)
		} else {
			l.log.Debug().
				Str("name", e.Name).
				Str("kind", e.Kind).
				Str("module", e.ModuleName).
				Msgf("run %s executed in module %s", e.Name, e.ModuleName)
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		l.log.Debug().
			Str("function", e.FunctionName).
			Str("module", e.ModuleName).
			Msgf("invoking function %s in module %s", e.FunctionName, e.ModuleName)
	case *fxevent.Invoked:
		if e.Err != nil {
			l.log.Error().
				Err(e.Err).
				Str("stack", e.Trace).
				Str("function", e.FunctionName).
				Str("module", e.ModuleName).
				Msgf("invoke failed for function %s in module %s", e.FunctionName, e.ModuleName)
		}
	case *fxevent.Stopping:
		l.log.Debug().
			Str("signal", strings.ToUpper(e.Signal.String())).
			Msgf("received signal %s", e.Signal.String())
	case *fxevent.Stopped:
		if e.Err != nil {
			l.log.Error().Err(e.Err).Msg("stop failed")
		}
	case *fxevent.RollingBack:
		l.log.Error().Err(e.StartErr).Msg("start failed")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.log.Error().Err(e.Err).Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.log.Error().Err(e.Err).Msg("start failed")
		} else {
			l.log.Debug().Msg("started application")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.log.Error().Err(e.Err).Msg("custom logger initialization failed")
		} else {
			l.log.Debug().Str("function", e.ConstructorName).Msg("initialized custom fxevent.Logger")
		}
	}
}

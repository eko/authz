package fx

import (
	"strings"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"golang.org/x/exp/slog"
)

var (
	Logger = fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
		return &SlogLogger{Logger: logger}
	})
)

type SlogLogger struct {
	Logger *slog.Logger
}

// LogEvent logs the given event to the provided structured-log logger.
func (l *SlogLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Logger.Debug("OnStart hook executing",
			slog.String("callee", e.FunctionName),
			slog.String("caller", e.CallerName),
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Logger.Error("OnStart hook failed", e.Err,
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
			)
		} else {
			l.Logger.Debug("OnStart hook executed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slog.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.OnStopExecuting:
		l.Logger.Debug("OnStop hook executing",
			slog.String("callee", e.FunctionName),
			slog.String("caller", e.CallerName),
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Logger.Debug("OnStop hook failed", e.Err,
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
			)
		} else {
			l.Logger.Debug("OnStop hook executed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slog.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.Supplied:
		l.Logger.Debug("Supplied", e.Err,
			slog.String("type", e.TypeName),
			moduleField(e.ModuleName),
		)
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.Debug("Provided",
				slog.String("constructor", e.ConstructorName),
				moduleField(e.ModuleName),
				slog.String("type", rtype),
			)
		}
		if e.Err != nil {
			l.Logger.Error("Error encountered while applying options", e.Err,
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.Debug("Replaced",
				moduleField(e.ModuleName),
				slog.String("type", rtype),
			)
		}
		if e.Err != nil {
			l.Logger.Error("Error encountered while replacing", e.Err,
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.Debug("Decorated",
				slog.String("decorator", e.DecoratorName),
				moduleField(e.ModuleName),
				slog.String("type", rtype),
			)
		}
		if e.Err != nil {
			l.Logger.Error("Error encountered while applying options", e.Err,
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		l.Logger.Debug("Invoking",
			slog.String("function", e.FunctionName),
			moduleField(e.ModuleName),
		)
	case *fxevent.Invoked:
		if e.Err != nil {
			l.Logger.Error("Invoke failed", e.Err,
				slog.String("stack", e.Trace),
				slog.String("function", e.FunctionName),
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Stopping:
		l.Logger.Debug("Received signal",
			slog.String("signal", strings.ToUpper(e.Signal.String())))
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.Error("Stop failed", e.Err)
		} else {
			l.Logger.Info("Stopped")
		}
	case *fxevent.RollingBack:
		l.Logger.Error("Start failed, rolling back", e.StartErr)
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Logger.Error("Rollback failed", e.Err)
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.Error("Start failed", e.Err)
		} else {
			l.Logger.Info("Started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Logger.Error("Custom logger initialization failed", e.Err)
		} else {
			l.Logger.Debug("Initialized custom fxevent.Logger", slog.String("function", e.ConstructorName))
		}
	}
}

func (SlogLogger) String() string { return "SlogLogger" }

func moduleField(name string) slog.Attr {
	if len(name) == 0 {
		return slog.Attr{}
	}
	return slog.String("module", name)
}

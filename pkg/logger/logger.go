package logger

import (
	"context"
	"io"

	appctx "mattermost/internal/context"

	"github.com/TheZeroSlave/zapsentry"

	"github.com/opentracing-contrib/go-zap/utils"
	"github.com/opentracing/opentracing-go"
	otext "github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	jconf "github.com/uber/jaeger-client-go/config"
	otzap "github.com/uber/jaeger-client-go/log/zap"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is alias for zap.Logger.
type Logger struct {
	*zap.Logger                    // underlying logger
	level       zapcore.Level      // level, for opentracing logs
	tracer      opentracing.Tracer // tracer or nil
	closer      io.Closer          // closer or nil
	name        string             // use name of a named log as service name
	fields      []zap.Field        // keep additional fields for tracing
}

// NewNop raps zap.NewNop logger.
func NewNop() *Logger {
	return &Logger{Logger: zap.NewNop()}
}

// Config of the logger.
type Config struct {
	Zap    ZapConfig `json:"zap" yaml:"zap" toml:"zap" mapstructure:"zap"`
	Sentry struct {
		DSN        string `json:"dsn" yaml:"dsn" toml:"dsn" mapstructure:"dsn"`
		Enviroment string `json:"environment" yaml:"environment" toml:"environment" mapstructure:"environment"` //nolint
		Release    string `json:"release" yaml:"release" toml:"release" mapstructure:"release"`                 //nolint
		Debug      bool   `json:"debug" yaml:"debug" toml:"debug" mapstructure:"debug"`                         //nolint
		Component  string `json:"component" yaml:"component" toml:"component" mapstructure:"component"`         //nolint
	} `json:"sentry" yaml:"sentry" toml:"sentry" mapstructure:"sentry"`
	OpenTracing JeagerConfig `json:"opentracing" yaml:"opentracing" toml:"opentracing" mapstructure:"opentracing"` //nolint
}

// NewConfig with defaults.
func NewConfig() (conf *Config) {
	conf = new(Config)
	conf.Zap.Encoding = "console"
	conf.Zap.OutputPaths = []string{"stdout"}
	conf.Zap.ErrorOutputPaths = []string{"stdout"}
	conf.Zap.Level = "debug"
	conf.Zap.EncoderConfig = "development"
	conf.Zap.DisableStacktrace = true

	// - no default sentry configurations
	// - no default open tracing configurations
	return
}

// New logger by given Config. The New can panic on some cases.
func New(conf *Config) (log *Logger) {
	var zapconf = conf.Zap.toZapConfig()
	var zaplog, err = zapconf.Build()
	if err != nil {
		panic(err)
	}
	if conf.Sentry.DSN != "" {
		var core zapcore.Core
		core, err = zapsentry.NewCore(zapsentry.Configuration{
			Level: zapcore.ErrorLevel, // top level to send messages to sentry
			Tags: map[string]string{
				"component": conf.Sentry.Component,
			},
		}, zapsentry.NewSentryClientFromDSN(conf.Sentry.DSN))
		if err != nil {
			zaplog.Panic("failed to initialize zap sentry core", zap.Error(err))
		}

		zaplog = zapsentry.AttachCoreToLogger(core, zaplog)
	}

	log = new(Logger)
	log.Logger = zaplog
	log.level = zapconf.Level.Level()

	var ctx = context.Background()
	log.Info(ctx, "logger created")

	if !conf.OpenTracing.Disabled {
		var otconf = conf.OpenTracing.toConfigurations()
		otconf.Sampler = &jconf.SamplerConfig{
			Type:  "const",
			Param: 1,
		}
		log.tracer, log.closer, err = otconf.NewTracer(
			jconf.Logger(otzap.NewLogger(zaplog.Named("opentracing"))),
		)
		if err != nil {
			zaplog.Panic("failed to initialize opentracing", zap.Error(err))
		}
		opentracing.SetGlobalTracer(log.tracer)
		log.Info(ctx, "open tracing configured")
	}

	if conf.OpenTracing.ServiceName == "" {
		log.Panic(ctx, "no service name configured at"+
			" log.opentracing.service_name")
	}

	// return already named service with equal name for
	// logs and for the opentracing
	log.Logger = log.Logger.Named(conf.OpenTracing.ServiceName)
	// keep name to use as prefix for child named logs
	log.name = conf.OpenTracing.ServiceName
	return
}

//
// open tracing
//

// Tracer, if any.
func (log *Logger) Tracer() opentracing.Tracer {
	return log.tracer
}

// Close, if any.
func (log *Logger) Close() (err error) {
	if log.closer != nil {
		err = log.closer.Close()
	}
	return
}

//
// overwrite Zap log methods adding optional context
//

func (log *Logger) span(ctx context.Context, level zapcore.Level,
	msg string, fields ...zapcore.Field) {

	if log.tracer == nil || level > log.level {
		return
	}

	var span = opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}

	const extend = 2
	var otFields = make([]otlog.Field, 0, len(log.fields)+len(fields)+extend)
	otFields = append(otFields, otlog.String("level", level.String()))

	if msg != "" {
		otFields = append(otFields, otlog.String("event", msg))
	}

	// add the 'With' fields
	fields = append(fields, log.fields...)

	if len(fields) > 0 {
		otFields = append(otFields, utils.ZapFieldsToOpentracing(fields...)...)
	}

	span.LogFields(otFields...)
}

// Span reference.
type Span = opentracing.Span

// Start an operation tracing. Returns updated context. Uses existing
// span as parent span if context already have a span. Finish
// the span in defer.
//
//      // DoSomething quickly and correctly.
//      func (obj *Object) DoSomething(ctx context.Context) (err error) {
//          var span log.Span
//          span, ctx = obj.log.Start(ctx, "do_something")
//          defer span.Finish()
//
//          // stuff
//
//          return
//      }
//
func (log *Logger) Start(ctx context.Context, operationName string,
	opts ...opentracing.StartSpanOption) (sp Span, sctx context.Context) {

	sp, sctx = opentracing.StartSpanFromContext(ctx, operationName, opts...)
	// use name of a named log as service name for traces of the span
	if log.tracer != nil && log.closer == nil && log.name != "" {
		otext.PeerService.Set(sp, log.name)
	}
	return
}

// DPanic of zap. The ctx can be nil. If not, the logger gives
// request_id and trace_id from it.
func (log *Logger) DPanic(ctx context.Context, msg string,
	fields ...zap.Field) {

	if ctx != nil {
		fields = appctx.AddContextFields(ctx, fields...)
	}
	log.span(ctx, zapcore.DPanicLevel, msg, fields...)
	log.Logger.DPanic(msg, fields...)
}

// Debug of zap. The ctx can be nil. If not, the logger gives
// request_id and trace_id from it.
func (log *Logger) Debug(ctx context.Context, msg string,
	fields ...zap.Field) {

	if ctx != nil {
		fields = appctx.AddContextFields(ctx, fields...)
	}
	log.span(ctx, zapcore.DebugLevel, msg, fields...)
	log.Logger.Debug(msg, fields...)
}

// Error of zap. The ctx can be nil. If not, the logger gives
// request_id and trace_id from it.
func (log *Logger) Error(ctx context.Context, msg string,
	fields ...zap.Field) {

	if ctx != nil {
		fields = appctx.AddContextFields(ctx, fields...)
	}
	log.span(ctx, zapcore.ErrorLevel, msg, fields...)
	log.Logger.Error(msg, fields...)
}

// Fatal of zap. The ctx can be nil. If not, the logger gives
// request_id and trace_id from it.
func (log *Logger) Fatal(ctx context.Context, msg string,
	fields ...zap.Field) {

	if ctx != nil {
		fields = appctx.AddContextFields(ctx, fields...)
	}
	log.span(ctx, zapcore.FatalLevel, msg, fields...)
	log.Logger.Fatal(msg, fields...)
}

// Info of zap. The ctx can be nil. If not, the logger gives
// request_id and trace_id from it.
func (log *Logger) Info(ctx context.Context, msg string,
	fields ...zap.Field) {

	if ctx != nil {
		fields = appctx.AddContextFields(ctx, fields...)
	}
	log.span(ctx, zapcore.InfoLevel, msg, fields...)
	log.Logger.Info(msg, fields...)
}

// Panic of zap. The ctx can be nil. If not, the logger gives
// request_id and trace_id from it.
func (log *Logger) Panic(ctx context.Context, msg string,
	fields ...zap.Field) {

	if ctx != nil {
		fields = appctx.AddContextFields(ctx, fields...)
	}
	log.span(ctx, zapcore.PanicLevel, msg, fields...)
	log.Logger.Panic(msg, fields...)
}

// Warn of zap. The ctx can be nil. If not, the logger gives
// request_id and trace_id from it.
func (log *Logger) Warn(ctx context.Context, msg string,
	fields ...zap.Field) {

	if ctx != nil {
		fields = appctx.AddContextFields(ctx, fields...)
	}
	log.span(ctx, zapcore.WarnLevel, msg, fields...)
	log.Logger.Warn(msg, fields...)
}

func (log *Logger) Named(name string) (named *Logger) {
	named = new(Logger)
	named.Logger = log.Logger.Named(log.name + "." + name)
	if log.tracer != nil {
		// don't copy closer to avoid closing of inherited tracer
		named.tracer = log.tracer //
		named.name = name         // use as service name in the traces
		named.level = log.level   //
		named.fields = log.fields //
	}
	return
}

func (log *Logger) With(fields ...zap.Field) (with *Logger) {
	with = new(Logger)
	with.Logger = log.Logger.With(fields...)
	if log.tracer != nil {
		// don't copy closer to avoid closing of inherited tracer
		with.tracer = log.tracer                    //
		with.level = log.level                      //
		with.fields = append(log.fields, fields...) //nolint
	}
	return
}

func (log *Logger) WithOptions(opts ...zap.Option) (wo *Logger) {
	wo = new(Logger)
	wo.Logger = log.Logger.WithOptions(opts...)
	return
}

// SetOperationName for Span in context if any.
func SetOperationName(ctx context.Context, name string) {
	var span = opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}
	span.SetOperationName(name)
}

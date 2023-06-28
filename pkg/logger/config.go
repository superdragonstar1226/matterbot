package logger

import (
	"time"

	"go.uber.org/zap"

	"github.com/opentracing/opentracing-go"
	jconf "github.com/uber/jaeger-client-go/config"
)

// ZapConfig for the Viper.
type ZapConfig struct {
	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	Level string `json:"level" yaml:"level" mapstructure:"level"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `json:"development" yaml:"development" mapstructure:"development"` //nolint
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller" mapstructure:"disableCaller"` //nolint
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace" mapstructure:"disableStacktrace"` //nolint
	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	Sampling *zap.SamplingConfig `json:"sampling" yaml:"sampling" mapstructure:"sampling"` //nolint
	// Encoding sets the logger's encoding. Valid values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	Encoding string `json:"encoding" yaml:"encoding" mapstructure:"encoding"` //nolint
	// EncoderConfig -- development or production
	EncoderConfig string `json:"encoderConfig" yaml:"encoderConfig" mapstructure:"encoderConfig"` //nolint
	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths" mapstructure:"outputPaths"` //nolint
	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths" mapstructure:"errorOutputPaths"` //nolint
	// InitialFields is a collection of fields to add to the root logger.
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields" mapstructure:"initialFields"` //nolint
}

func (zc *ZapConfig) toZapConfig() (conf *zap.Config) {
	conf = new(zap.Config)
	if err := conf.Level.UnmarshalText([]byte(zc.Level)); err != nil {
		panic(err)
	}
	conf.Development = zc.Development
	conf.DisableCaller = zc.DisableCaller
	conf.DisableStacktrace = zc.DisableStacktrace
	conf.Sampling = zc.Sampling
	conf.Encoding = zc.Encoding
	switch zc.EncoderConfig {
	case "development":
		conf.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	case "production":
		conf.EncoderConfig = zap.NewProductionEncoderConfig()
	default:
		panic("unknown or missing zap Config: EncoderConfig: " +
			zc.EncoderConfig)
	}
	conf.OutputPaths = zc.OutputPaths
	conf.ErrorOutputPaths = zc.ErrorOutputPaths
	conf.InitialFields = zc.InitialFields
	return
}

// SamplerConfig for the Viper. Add mapstructure tag, and use snake case.
type SamplerConfig struct {
	//
}

// ReporterConfig for the Viper. Add mapstructure tag, and use snake case.
type ReporterConfig struct {
	QueueSize                  int               `yaml:"queue_size" mapstructure:"queue_size"`                                     // nolint
	BufferFlushInterval        time.Duration     `yaml:"buffer_flush_interval" mapstructure:"buffer_flush_interval"`               // nolint
	LogSpans                   bool              `yaml:"log_spans" mapstructure:"log_spans"`                                       // nolint
	LocalAgentHostPort         string            `yaml:"local_agent_host_port" mapstructure:"local_agent_host_port"`               // nolint
	DisableAttemptReconnecting bool              `yaml:"disable_attempt_reconnecting" mapstructure:"disable_attempt_reconnecting"` // nolint
	AttemptReconnectInterval   time.Duration     `yaml:"attempt_reconnect_interval" mapstructure:"attempt_reconnect_interval"`     // nolint
	CollectorEndpoint          string            `yaml:"collector_endpoint" mapstructure:"collector_endpoint"`                     // nolint
	User                       string            `yaml:"user" mapstructure:"user"`
	Password                   string            `yaml:"password" mapstructure:"password"`         // nolint
	HTTPHeaders                map[string]string `yaml:"http_headers" mapstructure:"http_headers"` // nolint
}

func (rc *ReporterConfig) toReporterConfig() (conf *jconf.ReporterConfig) {
	if rc == nil {
		return // nil
	}
	conf = new(jconf.ReporterConfig)
	conf.QueueSize = rc.QueueSize
	conf.BufferFlushInterval = rc.BufferFlushInterval
	conf.LogSpans = rc.LogSpans
	conf.LocalAgentHostPort = rc.LocalAgentHostPort
	conf.DisableAttemptReconnecting = rc.DisableAttemptReconnecting
	conf.AttemptReconnectInterval = rc.AttemptReconnectInterval
	conf.CollectorEndpoint = rc.CollectorEndpoint
	conf.User = rc.User
	conf.Password = rc.Password
	conf.HTTPHeaders = rc.HTTPHeaders
	return
}

// // HeadersConfig for the Viper. Add mapstructure tag, and use snake case.
// type HeadersConfig struct {
// 	//
// }

// // BaggageRestrictionsConfig for the Viper. Add mapstructure tag,
// //  and use snake case.
// type BaggageRestrictionsConfig struct {
// 	//
// }

// // ThrottlerConfig for the Viper. Add mapstructure tag, and use snake case.
// type ThrottlerConfig struct {
// 	//
// }

// JeagerConfig for the Viper (add mapstructure tags) and use snake case.
type JeagerConfig struct {
	ServiceName string            `yaml:"service_name" mapstructure:"service_name"`
	Disabled    bool              `yaml:"disabled" mapstructure:"disabled"`
	RPCMetrics  bool              `yaml:"rpc_metrics" mapstructure:"rpc_metrics"`
	Gen128Bit   bool              `yaml:"traceid_128bit" mapstructure:"traceid_128bit"` //nolint
	Tags        []opentracing.Tag `yaml:"tags" mapstructure:"tags"`
	// Sampler             *SamplerConfig             `yaml:"sampler"`
	Reporter *ReporterConfig `yaml:"reporter" mapstructure:"reporter"`
	// Headers             *jaeger.HeadersConfig      `yaml:"headers"`
	// BaggageRestrictions *BaggageRestrictionsConfig `yaml:"baggage_restrictions"`
	// Throttler           *ThrottlerConfig           `yaml:"throttler"`
}

func (jc *JeagerConfig) toConfigurations() (conf jconf.Configuration) {
	conf.ServiceName = jc.ServiceName
	conf.Disabled = jc.Disabled
	conf.RPCMetrics = jc.RPCMetrics
	conf.Gen128Bit = jc.Gen128Bit
	conf.Tags = jc.Tags
	conf.Reporter = jc.Reporter.toReporterConfig() // nil is ok
	return
}

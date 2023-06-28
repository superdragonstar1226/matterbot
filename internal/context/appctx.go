package appctx

import (
	"context"

	"go.uber.org/zap"
)

var (
	RequestIDContextKey = "rqid" // keep it short
	TraceIDContextKey   = "trid" // keep it short
)

func AddContextFields(ctx context.Context, flds ...zap.Field) (
	all []zap.Field) {

	all = flds

	var reqID, traceID = RequestID(ctx), TraceID(ctx)
	if reqID != "" {
		all = append(all, zap.String("request_id", reqID))
	}
	if traceID != "" {
		all = append(all, zap.String("trace_id", traceID))
	}
	return
}

func TraceID(ctx context.Context) (traceID string) {
	traceID, _ = ctx.Value(TraceIDContextKey).(string)
	return
}

func RequestID(ctx context.Context) (reqID string) {
	reqID, _ = ctx.Value(RequestIDContextKey).(string)
	return
}

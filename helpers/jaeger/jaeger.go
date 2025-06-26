package jaeger

import (
	"context"

	jaegergin "github.com/Arrim/jaeger-gin"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publishing interface {
	GetHeaders() amqp.Table
	SetHeaders(map[string]any)
}

func StartSpanFromContext(
	ctx context.Context,
	operationName string,
	opts ...opentracing.StartSpanOption,
) (opentracing.Span, context.Context) {
	if gctx, ok := ctx.(*gin.Context); ok {
		return StartSpanFromGinContext(gctx, operationName, opts...)
	}

	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
}

func StartSpanFromGinContext(
	gCtx *gin.Context,
	operationName string,
	opts ...opentracing.StartSpanOption,
) (opentracing.Span, *gin.Context) {
	ctx := jaegergin.GetSpanFromContext(gCtx)

	span, ctx := opentracing.StartSpanFromContext(ctx, operationName, opts...)

	jaegergin.InjectSpanInGinContext(ctx, gCtx)

	return span, gCtx
}

func InjectSpanContextToAmqp(span opentracing.Span, msg Publishing) error {
	c := amqpHeadersCarrier(msg.GetHeaders())

	if err := span.Tracer().Inject(span.Context(), opentracing.TextMap, c); err != nil {
		return err
	}

	msg.SetHeaders(c)

	return nil
}

type InterfaceMapCarrier map[string]any

func (c InterfaceMapCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, val := range c {
		v, ok := val.(string)
		if !ok {
			continue
		}

		if err := handler(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (c InterfaceMapCarrier) Set(key, val string) {
	c[key] = val
}

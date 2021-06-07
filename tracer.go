package tracer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zerologger "github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"net/http"
)

func SetJaegerTracer(traceHeader string) (opentracing.Tracer, io.Closer, error) {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		zerologger.Error().Msg(fmt.Sprintf("Could not parse Jaeger env vars: %s", err.Error()))
	}
	customHeaders := &jaeger.HeadersConfig{
		TraceContextHeaderName: traceHeader,
	}
	cfg.Headers = customHeaders
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		zerologger.Error().Msg(fmt.Sprintf("Could not create Jaeger traces: %s", err.Error()))
	}
	return tracer, closer, err
}
func OpenTracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestSpan opentracing.Span
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			requestSpan = opentracing.StartSpan(c.Request.URL.Path)
			defer requestSpan.Finish()
		} else {
			requestSpan = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
			defer requestSpan.Finish()
		}
		c.Set("tracing-context", requestSpan)
		c.Next()
	}
}

func AddTraceToRequest(c *gin.Context, req *http.Request) {
	tracer := opentracing.GlobalTracer()
	if cspan, ok := c.Get("tracing-context"); ok {
		tracer.Inject(cspan.(opentracing.Span).Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	}

}

package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	tracer opentracing.Tracer
	closer io.Closer
	once   = sync.Once{}
)

func initJaeger(service, addr string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 60 * time.Second,
			LocalAgentHostPort:  addr,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("init Jaeger failed,err:%v\n", err))
	}
	return tracer, closer
}

func TraceHttpRoot(service, agentAddr string) gin.HandlerFunc {
	once.Do(func() {
		tracer, closer = initJaeger(service, agentAddr)
		opentracing.SetGlobalTracer(tracer)
	})
	hn, _ := os.Hostname()
	return func(c *gin.Context) {
		name := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		sp := tracer.StartSpan(name)
		defer sp.Finish()
		Inject(sp.Context(), c.Request.Header)
		c.Next()

		sp.SetTag("trace_id", GetTraceID(sp))
		sp.SetTag("common.user_email", c.Request.Header.Get("user_email"))
		sp.SetTag("common.status_code", c.Writer.Status())
		sp.SetTag("hostname", hn)
	}
}

func TraceHttpSpan(service, agentAddr string) gin.HandlerFunc {
	once.Do(func() {
		tracer, closer = initJaeger(service, agentAddr)
		opentracing.SetGlobalTracer(tracer)
	})
	return func(c *gin.Context) {
		name := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		sp, err := ExtractChildSpan(name, c.Request.Header)
		if err != nil {
			return
		}
		Inject(sp.Context(), c.Request.Header)
		c.Next()

		sp.SetTag("trace_id", GetTraceID(sp))
		sp.SetTag("common.user_email", c.Request.Header.Get("user_email"))
		sp.SetTag("common.status_code", c.Writer.Status())
		defer sp.Finish()
	}
}

func ExtractChildSpan(name string, carrier interface{}) (opentracing.Span, error) {
	ctx, err := ExtractContext(carrier)
	if err != nil {
		return nil, err
	}
	return StartChildSpan(name, ctx), nil
}

//提取
func ExtractContext(carrier interface{}) (opentracing.SpanContext, error) {
	switch t := carrier.(type) {
	case http.Header:
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(t))
		return spanCtx, nil
	case context.Context:
		if v := t.Value("JAEGER_CTX"); v != nil {
			if vCtx, ok := v.(opentracing.SpanContext); ok {
				return vCtx, nil
			} else {
				return nil, errors.New("no trace")
			}
		} else {
			return nil, errors.New("no trace")
		}
	default:
		return nil, errors.New("carrier type not supported")
	}
}

func StartChildSpan(name string, parentCtx opentracing.SpanContext) opentracing.Span {
	return tracer.StartSpan(name, ext.RPCServerOption(parentCtx))
}

//注入
func Inject(spCtx opentracing.SpanContext, carrier interface{}) {
	switch t := carrier.(type) {
	case http.Header:
		opentracing.GlobalTracer().Inject(spCtx, opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(t))
		t.Set("trace_id", spanContextToJaegerContext(spCtx).TraceID().String())
	}
}

func GetTraceID(span opentracing.Span) string {
	jaegerSpanContext := spanContextToJaegerContext(span.Context())
	return jaegerSpanContext.TraceID().String()
}

func GetSpanID(span opentracing.Span) string {
	jaegerSpanContext := spanContextToJaegerContext(span.Context())
	return jaegerSpanContext.SpanID().String()
}

func spanContextToJaegerContext(spanContext opentracing.SpanContext) jaeger.SpanContext {
	if sc, ok := spanContext.(jaeger.SpanContext); ok {
		return sc
	} else {
		return jaeger.SpanContext{}
	}
}

//将http header包装到ctx中
func ExtractStdContext(parent context.Context, carrier interface{}) context.Context {
	var parentCtx context.Context
	if parent != nil {
		parentCtx = parent
	} else {
		parentCtx = context.Background()
	}
	spCtx, err := ExtractContext(carrier)
	if err != nil {
		panic(err.Error())
	}
	return context.WithValue(parentCtx, "JAEGER_CTX", spCtx)
}

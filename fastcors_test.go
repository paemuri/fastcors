package fastcors

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestSetAllowedOrigins(t *testing.T) {
	opt := &Options{}
	allowedOrigins := []string{"batata"}
	SetAllowedOrigins(allowedOrigins)(opt)
	if len(opt.AllowedOrigins) != 1 ||
		opt.AllowedOrigins[0] != allowedOrigins[0] {
		t.Errorf("expected %v got %v", allowedOrigins, opt.AllowedOrigins)
	}
}

func TestSetAllowOriginFunc(t *testing.T) {
	opt := &Options{}
	allowOriginFunc := func(origin string) bool { return true }
	SetAllowOriginFunc(allowOriginFunc)(opt)
	if opt.AllowOriginFunc == nil ||
		opt.AllowOriginFunc("") != allowOriginFunc("") {
		t.Errorf("expected %v got %v", allowOriginFunc(""), opt.AllowOriginFunc(""))
	}
}

func TestSetAllowOriginRequestFunc(t *testing.T) {
	opt := &Options{}
	allowOriginRequestFunc := func(ctx *fasthttp.RequestCtx, origin string) bool { return true }
	SetAllowOriginRequestFunc(allowOriginRequestFunc)(opt)
	if opt.AllowOriginRequestFunc == nil ||
		opt.AllowOriginRequestFunc(nil, "") != allowOriginRequestFunc(nil, "") {
		t.Errorf("expected %v got %v", allowOriginRequestFunc(nil, ""), opt.AllowOriginRequestFunc(nil, ""))
	}
}

func TestSetAllowedMethods(t *testing.T) {
	opt := &Options{}
	allowedMethods := []string{"batata"}
	SetAllowedMethods(allowedMethods)(opt)
	if len(opt.AllowedMethods) != 1 ||
		opt.AllowedMethods[0] != allowedMethods[0] {
		t.Errorf("expected %v got %v", allowedMethods, opt.AllowedMethods)
	}
}

func TestSetAllowedHeaders(t *testing.T) {
	opt := &Options{}
	allowedHeaders := []string{"batata"}
	SetAllowedHeaders(allowedHeaders)(opt)
	if len(opt.AllowedHeaders) != 1 ||
		opt.AllowedHeaders[0] != allowedHeaders[0] {
		t.Errorf("expected %v got %v", allowedHeaders, opt.AllowedHeaders)
	}
}

func TestSetExposedHeaders(t *testing.T) {
	opt := &Options{}
	exposedHeaders := []string{"batata"}
	SetExposedHeaders(exposedHeaders)(opt)
	if len(opt.ExposedHeaders) != 1 ||
		opt.ExposedHeaders[0] != exposedHeaders[0] {
		t.Errorf("expected %v got %v", exposedHeaders, opt.ExposedHeaders)
	}
}

func TestSetMaxAge(t *testing.T) {
	opt := &Options{}
	maxAge := 69
	SetMaxAge(maxAge)(opt)
	if opt.MaxAge != maxAge {
		t.Errorf("expected %v got %v", maxAge, opt.MaxAge)
	}
}

func TestSetAllowCredentials(t *testing.T) {
	opt := &Options{}
	allowCredentials := true
	SetAllowCredentials(allowCredentials)(opt)
	if opt.AllowCredentials != allowCredentials {
		t.Errorf("expected %v got %v", allowCredentials, opt.AllowCredentials)
	}
}

func TestSetDebug(t *testing.T) {
	opt := &Options{}
	debug := true
	SetDebug(debug)(opt)
	if opt.Debug != debug {
		t.Errorf("expected %v got %v", debug, opt.Debug)
	}
}

func TestSetLogger(t *testing.T) {
	opt := &Options{}
	logger := &testLogger{}
	SetLogger(logger)(opt)
	if opt.Logger != logger {
		t.Errorf("expected %v got %v", logger, opt.Logger)
	}
}

func TestAllowAll(t *testing.T) {
	opt := &Options{}
	AllowAll()(opt)
	if len(opt.AllowedOrigins) != 1 || opt.AllowedOrigins[0] != "*" {
		t.Errorf("expected all origins got %v", opt.AllowedOrigins)
	}
	if len(opt.AllowedMethods) != 6 ||
		(opt.AllowedMethods[0] != fasthttp.MethodPost &&
			opt.AllowedMethods[1] != fasthttp.MethodGet &&
			opt.AllowedMethods[2] != fasthttp.MethodPost &&
			opt.AllowedMethods[3] != fasthttp.MethodPut &&
			opt.AllowedMethods[4] != fasthttp.MethodPatch &&
			opt.AllowedMethods[5] != fasthttp.MethodDelete) {
		t.Errorf("expected all methods got %v", opt.AllowedMethods)
	}
	if len(opt.AllowedHeaders) != 1 || opt.AllowedHeaders[0] != "*" {
		t.Errorf("expected all headers got %v", opt.AllowedHeaders)
	}
}

func TestNew_WithSetAllowedOrigins(t *testing.T) {
	New(SetAllowedOrigins([]string{"https://localhost"}))
}

func TestNew_WithSetAllowOriginFunc(t *testing.T) {
	New(SetAllowOriginFunc(func(origin string) bool {
		return true
	}))
}

func TestNew_WithSetAllowOriginRequestFunc(t *testing.T) {
	New(SetAllowOriginRequestFunc(func(ctx *fasthttp.RequestCtx, origin string) bool {
		return true
	}))
}

func TestNew_WithSetAllowedMethods(t *testing.T) {
	New(SetAllowedMethods([]string{
		fasthttp.MethodGet,
		fasthttp.MethodPost,
	}))
}

func TestNew_WithSetAllowedHeaders(t *testing.T) {
	New(SetAllowedHeaders([]string{
		fasthttp.HeaderOrigin,
		fasthttp.HeaderContentType,
		fasthttp.HeaderAccept,
	}))
}

func TestNew_WithSetExposedHeaders(t *testing.T) {
	New(SetExposedHeaders([]string{
		fasthttp.HeaderOrigin,
		fasthttp.HeaderContentType,
		fasthttp.HeaderAccept,
	}))
}

func TestNew_WithSetMaxAge(t *testing.T) {
	New(SetMaxAge(60 * 60 * 24 * 30))
}

func TestNew_WithSetAllowCredentials(t *testing.T) {
	New(SetAllowCredentials(true))
}

func TestNew_WithSetDebug(t *testing.T) {
	New(SetDebug(true))
}

func TestNew_WithSetLogger(t *testing.T) {
	New(SetLogger(&testLogger{}))
}

func TestNew_WithAllowAll(t *testing.T) {
	New(AllowAll())
}

func TestNew_WithAllOptions(t *testing.T) {
	New(
		SetAllowedOrigins([]string{"https://localhost"}),
		SetAllowOriginFunc(func(origin string) bool {
			return true
		}),
		SetAllowOriginRequestFunc(func(ctx *fasthttp.RequestCtx, origin string) bool {
			return true
		}),
		SetAllowedMethods([]string{
			fasthttp.MethodGet,
			fasthttp.MethodPost,
		}),
		SetAllowedHeaders([]string{
			fasthttp.HeaderOrigin,
			fasthttp.HeaderContentType,
			fasthttp.HeaderAccept,
		}),
		SetExposedHeaders([]string{
			fasthttp.HeaderOrigin,
			fasthttp.HeaderContentType,
			fasthttp.HeaderAccept,
		}),
		SetMaxAge(60*60*24*30),
		SetAllowCredentials(true),
		SetDebug(true),
		SetLogger(&testLogger{}),
	)
}

type testLogger struct{}

func (*testLogger) Printf(string, ...interface{}) {}

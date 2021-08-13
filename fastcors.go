package fastcors

import (
	"log"
	"net/textproto"
	"os"
	"strconv"
	"strings"

	"github.com/paemuri/fastcors/lib/cors"
	"github.com/valyala/fasthttp"
)

type RequestMiddleware func(fasthttp.RequestHandler) fasthttp.RequestHandler

type Options struct {
	AllowedOrigins         []string
	AllowOriginFunc        func(origin string) bool
	AllowOriginRequestFunc func(ctx *fasthttp.RequestCtx, origin string) bool
	AllowedMethods         []string
	AllowedHeaders         []string
	ExposedHeaders         []string
	MaxAge                 int
	AllowCredentials       bool
	Debug                  bool
	Logger                 Logger
}

func SetAllowedOrigins(allowedOrigins []string) func(*Options) {
	return func(opt *Options) {
		opt.AllowedOrigins = allowedOrigins
	}
}

func SetAllowOriginFunc(allowOriginFunc func(origin string) bool) func(*Options) {
	return func(opt *Options) {
		opt.AllowOriginFunc = allowOriginFunc
	}
}

func SetAllowOriginRequestFunc(allowOriginRequestFunc func(ctx *fasthttp.RequestCtx, origin string) bool) func(*Options) {
	return func(opt *Options) {
		opt.AllowOriginRequestFunc = allowOriginRequestFunc
	}
}

func SetAllowedMethods(allowedMethods []string) func(*Options) {
	return func(opt *Options) {
		opt.AllowedMethods = allowedMethods
	}
}

func SetAllowedHeaders(allowedHeaders []string) func(*Options) {
	return func(opt *Options) {
		opt.AllowedHeaders = allowedHeaders
	}
}

func SetExposedHeaders(exposedHeaders []string) func(*Options) {
	return func(opt *Options) {
		opt.ExposedHeaders = exposedHeaders
	}
}

func SetMaxAge(maxAge int) func(*Options) {
	return func(opt *Options) {
		opt.MaxAge = maxAge
	}
}

func SetAllowCredentials(allowCredentials bool) func(*Options) {
	return func(opt *Options) {
		opt.AllowCredentials = allowCredentials
	}
}

func SetDebug(debug bool) func(*Options) {
	return func(opt *Options) {
		opt.Debug = debug
	}
}

func SetLogger(logger Logger) func(*Options) {
	return func(opt *Options) {
		opt.Logger = logger
	}
}

func AllowAll() func(*Options) {
	return func(opt *Options) {
		opt.AllowedOrigins = []string{"*"}
		opt.AllowedMethods = []string{
			fasthttp.MethodPost,
			fasthttp.MethodGet,
			fasthttp.MethodPost,
			fasthttp.MethodPut,
			fasthttp.MethodPatch,
			fasthttp.MethodDelete,
		}
		opt.AllowedHeaders = []string{"*"}
	}
}

type Logger interface {
	Printf(string, ...interface{})
}

type middleware struct {
	allowedOrigins         map[string]struct{}
	allowedOriginsAll      bool
	allowOriginFunc        func(origin string) bool
	allowOriginRequestFunc func(ctx *fasthttp.RequestCtx, origin string) bool
	allowedMethods         map[string]struct{}
	allowedHeaders         map[string]struct{}
	allowedHeadersAll      bool
	exposedHeaders         string
	maxAge                 int
	allowCredentials       bool
	debug                  bool
	log                    Logger
}

func normalizeHeader(h string) string {
	return textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(h))
}

func New(opts ...func(*Options)) RequestMiddleware {

	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	mw := new(middleware)
	mw.allowCredentials = options.AllowCredentials
	mw.maxAge = options.MaxAge
	for i, header := range options.ExposedHeaders {
		if i != 0 {
			mw.exposedHeaders += ", "
		}
		mw.exposedHeaders += normalizeHeader(header)
	}
	if options.Debug {
		mw.debug = true
		mw.log = options.Logger
		if mw.log == nil {
			mw.log = log.New(os.Stdout, "[cors] ", log.LstdFlags)
		}
	}

	if len(options.AllowedOrigins) == 0 {
		if options.AllowOriginFunc == nil && options.AllowOriginRequestFunc == nil {
			mw.allowedOriginsAll = true
		}
	} else {
		mw.allowedOrigins = make(map[string]struct{}, len(options.AllowedOrigins))
		for _, origin := range options.AllowedOrigins {
			origin = strings.ToLower(origin)
			if origin == "*" {
				mw.allowedOriginsAll = true
				mw.allowedOrigins = nil
				break
			}
			mw.allowedOrigins[origin] = struct{}{}
		}
	}
	if options.AllowOriginFunc != nil {
		mw.allowOriginFunc = options.AllowOriginFunc
		mw.allowedOriginsAll = false
		mw.allowedOrigins = nil
	}
	if options.AllowOriginRequestFunc != nil {
		mw.allowOriginRequestFunc = options.AllowOriginRequestFunc
		mw.allowedOriginsAll = false
		mw.allowedOrigins = nil
		mw.allowOriginFunc = nil
	}

	if len(options.AllowedMethods) == 0 {
		mw.allowedMethods = map[string]struct{}{
			fasthttp.MethodGet:  {},
			fasthttp.MethodPost: {},
			fasthttp.MethodHead: {},
		}
	} else {
		mw.allowedMethods = make(map[string]struct{}, len(options.AllowedMethods))
		for _, method := range options.AllowedMethods {
			mw.allowedMethods[strings.ToUpper(method)] = struct{}{}
		}
	}

	if len(options.AllowedHeaders) == 0 {
		mw.allowedHeaders = map[string]struct{}{
			fasthttp.HeaderOrigin:         {},
			fasthttp.HeaderAccept:         {},
			fasthttp.HeaderContentType:    {},
			fasthttp.HeaderXRequestedWith: {},
		}
	} else {
		mw.allowedHeaders = make(map[string]struct{}, len(options.AllowedHeaders))
		for _, header := range options.AllowedHeaders {
			if header == "*" {
				mw.allowedHeadersAll = true
				mw.allowedHeaders = nil
				break
			}
			mw.allowedHeaders[normalizeHeader(header)] = struct{}{}
		}
	}

	return mw.middleware
}

func (mw *middleware) middleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if ctx.IsOptions() && len(ctx.Request.Header.Peek(fasthttp.HeaderAccessControlRequestMethod)) > 0 {
			mw.handlePreflight(ctx)
			ctx.SetStatusCode(fasthttp.StatusOK) // some legacy browsers choke on 204
		} else {
			mw.handleActual(ctx)
			next(ctx)
		}
	}
}

func (mw *middleware) handlePreflight(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Add("Vary", fasthttp.HeaderOrigin)
	ctx.Response.Header.Add("Vary", fasthttp.HeaderAccessControlRequestMethod)
	ctx.Response.Header.Add("Vary", fasthttp.HeaderAccessControlRequestHeaders)

	origin := string(ctx.Request.Header.Peek("Origin"))
	if len(origin) == 0 {
		mw.logf("preflight request aborted: missing origin")
		return
	}
	if !mw.isOriginAllowed(ctx, origin) {
		mw.logf("preflight request aborted: origin %s not allowed", origin)
		return
	}
	method := string(ctx.Request.Header.Peek(fasthttp.HeaderAccessControlRequestMethod))
	if !mw.isMethodAllowed(method) {
		mw.logf("preflight request aborted: method %s not allowed", method)
		return
	}

	headers := ctx.Request.Header.Peek(fasthttp.HeaderAccessControlRequestHeaders)
	if !mw.areHeadersAllowed(headers) {
		mw.logf("preflight request aborted: headers %v not allowed", headers)
		return
	}

	if mw.allowedOriginsAll {
		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, "*")
	} else {
		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, origin)
	}
	ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowMethods, method)
	if len(headers) > 0 {
		ctx.Response.Header.SetBytesV(fasthttp.HeaderAccessControlAllowHeaders, headers)
	}
	if mw.allowCredentials {
		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowCredentials, "true")
	}
	if mw.maxAge > 0 {
		ctx.Response.Header.Set(fasthttp.HeaderAccessControlMaxAge, strconv.Itoa(mw.maxAge))
	}
}

func (mw *middleware) handleActual(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Add("Vary", fasthttp.HeaderOrigin)

	origin := string(ctx.Request.Header.Peek("Origin"))
	if len(origin) == 0 {
		mw.logf("actual request aborted: missing origin")
		return
	}
	if !mw.isOriginAllowed(ctx, origin) {
		mw.logf("actual request aborted: origin %s not allowed", origin)
		return
	}
	method := string(ctx.Request.Header.Method())
	if !mw.isMethodAllowed(method) {
		mw.logf("actual request aborted: method %s not allowed", method)
		return
	}

	if mw.allowedOriginsAll {
		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, "*")
	} else {
		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, origin)
	}
	if len(mw.exposedHeaders) > 0 {
		ctx.Response.Header.Set(fasthttp.HeaderAccessControlExposeHeaders, mw.exposedHeaders)
	}
	if mw.allowCredentials {
		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowCredentials, "true")
	}
}

func (mw *middleware) logf(format string, a ...interface{}) {
	if mw.debug {
		mw.log.Printf(format, a...)
	}
}

func (mw *middleware) isOriginAllowed(ctx *fasthttp.RequestCtx, origin string) bool {
	if mw.allowedOriginsAll {
		return true
	}
	if mw.allowOriginRequestFunc != nil {
		return mw.allowOriginRequestFunc(ctx, origin)
	}
	if mw.allowOriginFunc != nil {
		return mw.allowOriginFunc(origin)
	}
	_, ok := mw.allowedOrigins[strings.ToLower(origin)]
	return ok
}

func (mw *middleware) isMethodAllowed(method string) bool {
	method = strings.ToUpper(method)
	if method == fasthttp.MethodOptions {
		return true
	}
	_, ok := mw.allowedMethods[method]
	return ok
}

func (mw *middleware) areHeadersAllowed(rawHeaders []byte) bool {
	if mw.allowedHeadersAll || len(rawHeaders) == 0 {
		return true
	}
	for _, header := range cors.ParseHeaderList(rawHeaders) {
		if header == fasthttp.HeaderOrigin {
			continue
		}
		if _, ok := mw.allowedHeaders[header]; !ok {
			return false
		}
	}
	return true
}

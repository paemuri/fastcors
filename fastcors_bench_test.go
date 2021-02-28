package fastcors

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func sampleHandler(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("voilà")
}

func BenchmarkHandle_PreflightWithDefaultMW(b *testing.B) {
	var ctx fasthttp.RequestCtx
	ctx.Request.SetRequestURI("http://aaa.com/fff/ss.html?sdf")
	ctx.Request.Header.SetMethod(fasthttp.MethodOptions)
	ctx.Request.Header.Set(fasthttp.HeaderOrigin, "http://aaa.com")
	ctx.Request.Header.Set(fasthttp.HeaderAccessControlRequestMethod, fasthttp.MethodGet)
	ctx.Request.Header.Add(fasthttp.HeaderAccessControlRequestHeaders, fasthttp.HeaderOrigin)
	ctx.Request.Header.Add(fasthttp.HeaderAccessControlRequestHeaders, fasthttp.HeaderAccept)
	ctx.Request.Header.Add(fasthttp.HeaderAccessControlRequestHeaders, fasthttp.HeaderContentType)
	ctx.Request.Header.Add(fasthttp.HeaderAccessControlRequestHeaders, fasthttp.HeaderXRequestedWith)
	benchmarkHandle(b, &ctx, New())
}

func BenchmarkHandle_PreflightWithAllowAllMW(b *testing.B) {
	var ctx fasthttp.RequestCtx
	ctx.Request.SetRequestURI("http://aaa.com/fff/ss.html?sdf")
	ctx.Request.Header.SetMethod(fasthttp.MethodOptions)
	ctx.Request.Header.Set(fasthttp.HeaderOrigin, "http://aaa.com")
	ctx.Request.Header.Set(fasthttp.HeaderAccessControlRequestMethod, fasthttp.MethodGet)
	ctx.Request.Header.Add(fasthttp.HeaderAccessControlRequestHeaders, fasthttp.HeaderOrigin)
	ctx.Request.Header.Add(fasthttp.HeaderAccessControlRequestHeaders, fasthttp.HeaderAccept)
	ctx.Request.Header.Add(fasthttp.HeaderAccessControlRequestHeaders, fasthttp.HeaderContentType)
	ctx.Request.Header.Add(fasthttp.HeaderAccessControlRequestHeaders, fasthttp.HeaderXRequestedWith)
	benchmarkHandle(b, &ctx, New(AllowAll()))
}

func BenchmarkHandle_PreflightWithoutOrigin(b *testing.B) {
	var ctx fasthttp.RequestCtx
	ctx.Request.SetRequestURI("http://aaa.com/fff/ss.html?sdf")
	ctx.Request.Header.SetMethod(fasthttp.MethodOptions)
	benchmarkHandle(b, &ctx, New())
}

func BenchmarkHandle_ActualWithDefaultMW(b *testing.B) {
	var ctx fasthttp.RequestCtx
	ctx.Request.SetRequestURI("http://aaa.com/fff/ss.html?sdf")
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	ctx.Request.Header.Set(fasthttp.HeaderOrigin, "http://aaa.com")
	benchmarkHandle(b, &ctx, New())
}

func BenchmarkHandle_ActualWithAllowAllMW(b *testing.B) {
	var ctx fasthttp.RequestCtx
	ctx.Request.SetRequestURI("http://aaa.com/fff/ss.html?sdf")
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	ctx.Request.Header.Set(fasthttp.HeaderOrigin, "http://aaa.com")
	benchmarkHandle(b, &ctx, New(AllowAll()))
}

func BenchmarkHandle_ActualWithoutOrigin(b *testing.B) {
	var ctx fasthttp.RequestCtx
	ctx.Request.SetRequestURI("http://aaa.com/fff/ss.html?sdf")
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	benchmarkHandle(b, &ctx, New())
}

func benchmarkHandle(b *testing.B, ctx *fasthttp.RequestCtx, mw RequestMiddleware) {
	handle := mw(sampleHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handle(ctx)
	}
}

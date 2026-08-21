package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gitcpu-io/zgo"
	"github.com/kataras/iris/v12"
)

// TraceGet TraceGet接口四板斧，这仅仅是一个例子
func TraceGet(ctx iris.Context) {
	trace := zgo.Trace.NewTrace("Get " + ctx.Path())

	// 第一：定义错误返回变量，请求上下文，通过defer来最后响应
	var errStr string

	cotx, cancel := context.WithTimeout(context.Background(), 1*time.Second) //you can change this time number
	defer cancel()

	defer func() {
		if errStr != "" {
			err := zgo.Http.JsonpErr(ctx, errStr)
			if err != nil {
				zgo.Log.Error(err)
			}
		}
		trace.Step("Response")
	}()

	// 第二：解析请求参数
	ms := ctx.URLParam("ms")
	if ms == "" {
		errStr = "必须输入query参数ms,默认是毫秒,比如：300ms"
		return
	}

	pd, err := time.ParseDuration(ms)
	if err != nil {
		errStr = "转化参数为时间可能丢失单位，比如ms. " + err.Error()
		return
	}
	defer trace.LogIfLong(pd) //最长超过多少时才打印Trace日志

	//第三：模拟暂停
	time.Sleep(pd)
	trace.Step(fmt.Sprintf("Stop %s Sleep", ms))

	// 第四：使用select来响应处理结果与超时
	select {
	case <-cotx.Done():
		errStr = "默认 1秒超时 timeout"
		zgo.Log.Error(errStr) //通过zgo.Log统计日志
		trace.Step("Select Timeout")
	default:
		err := zgo.Http.JsonpOK(ctx, fmt.Sprintf("OK 后台日志显示，暂停 %s 后显示Trace日志", ms))
		if err != nil {
			zgo.Log.Error(err)
		}
	}

}

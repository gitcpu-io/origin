package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gitcpu-io/zgo"
	"github.com/kataras/iris/v12"
)

// RedisGet 接口四板斧，这仅仅是一个例子
func RedisGet(ctx iris.Context) {
	// 第一：定义错误返回变量，请求上下文，通过defer来最后响应
	var errStr string

	cotx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //you can change this time number
	defer cancel()

	defer func() {
		if errStr != "" {
			err := zgo.Http.JsonpErr(ctx, errStr)
			if err != nil {
				zgo.Log.Error(err)
			}
		}
	}()

	// 第二：解析请求参数
	name := ctx.URLParam("name")
	if name == "" {
		errStr = "必须输入query参数name"
		return
	}

	key := fmt.Sprintf("%s:%s:%s", "zgo", "origin", name)

	//先写入
	_, err := zgo.Redis.Set(cotx, key, fmt.Sprintf("写入时间戳 %v", zgo.Utils.NowUnix()))
	if err != nil {
		zgo.Log.Error(err)
		return
	}

	// 第三：调用zgo engine来处理业务逻辑
	val, err := zgo.Redis.Get(cotx, key)
	if err != nil {
		errStr = err.Error()
		zgo.Log.Error(err)
		return
	}

	result := zgo.Utils.StringToMap(val.(string))

	if result == nil {
		m := make(map[string]interface{})
		m[name] = val
		result = m
	}

	// 第四：使用select来响应处理结果与超时
	select {
	case <-cotx.Done():
		errStr = "call redis get string timeout"
		zgo.Log.Error(errStr) //通过zgo.Log统计日志
	default:
		err := zgo.Http.JsonpOK(ctx, result)
		if err != nil {
			zgo.Log.Error(err)
		}
	}

}

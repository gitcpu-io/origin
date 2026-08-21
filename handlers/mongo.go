package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gitcpu-io/zgo"
	"github.com/kataras/iris/v12"
)

type User struct {
	Id       zgo.MongoObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Username string            `json:"username" bson:"username" `
	Age      int               `json:"age" bson:"age"`
	Address  int               `json:"address" bson:"address"`
}

// MongoGet 接口四板斧，这仅仅是一个例子
func MongoGet(ctx iris.Context) {
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

	// 第三：调用zgo engine来处理业务逻辑
	result, err := FindOne(cotx, name)
	if err != nil {
		errStr = err.Error()
		zgo.Log.Error(err)
		return
	}

	// 第四：使用select来响应处理结果与超时
	select {
	case <-cotx.Done():
		errStr = "call mongo get string timeout"
		zgo.Log.Error(errStr) //通过zgo.Log统计日志
	default:
		err := zgo.Http.JsonpOK(ctx, result)
		if err != nil {
			zgo.Log.Error(err)
		}
	}

}

// MongoList 接口四板斧，这仅仅是一个例子
func MongoList(ctx iris.Context) {
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

	// 第三：调用zgo engine来处理业务逻辑
	result, err := Find(cotx, name)
	if err != nil {
		errStr = err.Error()
		zgo.Log.Error(err)
		return
	}

	// 第四：使用select来响应处理结果与超时
	select {
	case <-cotx.Done():
		errStr = "call mongo list string timeout"
		zgo.Log.Error(errStr) //通过zgo.Log统计日志
	default:
		err := zgo.Http.JsonpOK(ctx, result)
		if err != nil {
			zgo.Log.Error(err)
		}
	}

}

func FindOne(ctx context.Context, username string) (*User, error) {
	var collection = zgo.Mongo.GetCollection("profile", "bj", "mongo_label_bj")

	filter := make(map[string]interface{}) //查询username是且age >= 30的
	filter["username"] = username
	filter["age"] = map[string]interface{}{
		"$gte": 30,
	}

	sort := make(map[string]interface{})
	sort["_id"] = -1 //-1降序，1升序

	//返回错误：Projection cannot have a mix of inclusion and exclusion; 要么是1，要么是0
	fields := make(map[string]interface{})
	fields["age"] = 1 //要么全是1，要么全是0
	fields["address"] = 1
	fields["username"] = 1

	r := &User{}

	//组织args
	args := &zgo.MongoArgs{
		Filter: filter, //查询条件
		Fields: fields, //对查询出的结果项，筛选字段
		Sort:   sort,   //排序
		Result: r,      //传入 &User{} ,结果
	}

	_, err := zgo.Mongo.FindOne(ctx, collection, args)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func Find(ctx context.Context, username string) ([]*User, error) {
	var collection = zgo.Mongo.GetCollection("profile", "bj", "mongo_label_bj")

	filter := make(map[string]interface{}) //查询username是且age >= 30的
	filter["username"] = username
	filter["age"] = map[string]interface{}{
		"$gte": 10,
	}

	sort := make(map[string]interface{})
	sort["_id"] = -1

	//返回错误：Projection cannot have a mix of inclusion and exclusion; 要么是1，要么是0
	fields := make(map[string]interface{})
	fields["age"] = 1
	fields["address"] = 1 //要么全是1，要么全是0
	fields["username"] = 1

	//组织args
	args := &zgo.MongoArgs{
		Filter: filter, //查询条件
		Fields: fields, //对查询出的结果项，筛选字段
		Sort:   sort,   //排序
		Limit:  10,     //查询结果数量
		Skip:   0,      //从哪一条开始跳过 开区间，不包括skip的值
	}

	results, err := zgo.Mongo.Find(ctx, collection, args)
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0)
	for _, v := range results {
		u := User{}
		err := zgo.Utils.BsonUnmarshal(v, &u) //对每一条数据进行 bsonUnmarshal 转为go结构体
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, &u)
	}

	return users, nil
}

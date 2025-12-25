package resultx

import (
	"chat/pkg/xerr"
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zrpcErr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"net/http"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(data any) *Response {
	return &Response{
		Code: 200,
		Msg:  "",
		Data: data,
	}
}

func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}

func OkHandler(_ context.Context, v any) any {
	return Success(v)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		errCode := xerr.SERVER_COMMON_ERROR
		errMsg := xerr.ErrMsg(errCode)

		// 提取根错误，断言是哪种类型
		causeErr := errors.Cause(err)
		var e *zrpcErr.CodeMsg
		if errors.As(causeErr, &e) {
			errCode = e.Code
			errMsg = e.Msg
		} else {
			// 在拦截器中已经将错误设置到了grpc的响应结构中，所以这里从gStatus获取并设置
			if gStatus, ok := status.FromError(causeErr); ok {
				errCode = int(gStatus.Code())
				errMsg = gStatus.Message()
			}
		}

		// 日志记录
		logx.WithContext(ctx).Errorf("【%s】 err %v", name, err)

		return http.StatusBadRequest, Fail(errCode, errMsg)
	}
}

package rpcserver

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LogInterceptor 把go-zero业务产生的自定义错误转换成gRPC的标准错误
func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any,
	err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	logx.WithContext(ctx).Errorf("【RPC SRV ERR】 %v", err)

	// 对错误溯源，拿到原始错误
	causeErr := errors.Cause(err)
	// 判断是不是我们自定义的错误对象
	var e *zerr.CodeMsg
	if errors.As(causeErr, &e) {
		err = status.Error(codes.Code(e.Code), e.Msg)
	}

	return resp, err
}

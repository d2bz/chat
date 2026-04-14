package retryjob

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

var ErrJobTimeout = errors.New("job timeout")

// RetryJetLagFunc 定义重试的时间策略
type RetryJetLagFunc func(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration

// RetryJetLagAlways 默认固定间隔重试方法
func RetryJetLagAlways(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration {
	return DefaultRetryJetLag
}

// IsRetryFunc 判断是否进行重试
type IsRetryFunc func(ctx context.Context, retryCount int, err error) bool

func RetryAlways(ctx context.Context, retryCount int, err error) bool {
	return true
}

// WithRetry 如果handler没用监听ctx.Done()会出现协程泄漏问题
func WithRetry(ctx context.Context, handler func(ctx context.Context) error, opts ...RetryOptions) error {
	opt := newOptions(opts...)

	// 判断程序本身是否设置了超时
	_, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.timeout)
		defer cancel()
	}

	var (
		herr        error
		retryJetLag time.Duration
		ch          = make(chan error, 1)
	)

	for i := 0; i < opt.retryNums; i++ {
		go func() {
			ch <- handler(ctx)
		}()
		select {
		case herr = <-ch:
			if herr == nil {
				return nil
			}
			if !opt.isRetryFunc(ctx, i, herr) {
				return herr
			}
			retryJetLag = opt.retryJetLag(ctx, i, retryJetLag)
			time.Sleep(retryJetLag)
		case <-ctx.Done():
			return ErrJobTimeout
		}
	}
	return herr
}

// SafeWithRetry 牺牲超时后的响应速度，保证协程不会泄漏
func SafeWithRetry(ctx context.Context, handler func(ctx context.Context) error, opts ...RetryOptions) error {
	opt := newOptions(opts...)

	// 判断程序本身是否设置了超时
	_, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.timeout)
		defer cancel()
	}

	var (
		herr        error
		retryJetLag time.Duration
	)

	for i := 0; i < opt.retryNums; i++ {
		if ctx.Err() != nil {
			return ErrJobTimeout
		}

		herr = handler(ctx)

		if herr == nil {
			return nil
		}

		if !opt.isRetryFunc(ctx, i, herr) {
			return herr
		}

		retryJetLag = opt.retryJetLag(ctx, i, retryJetLag)

		select {
		case <-time.After(retryJetLag):
		case <-ctx.Done():
			return ErrJobTimeout
		}
	}
	return herr
}

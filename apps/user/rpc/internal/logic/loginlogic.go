package logic

import (
	"chat/apps/user/models"
	"chat/pkg/ctxdata"
	"chat/pkg/encrypt"
	"chat/pkg/xerr"
	"context"
	"github.com/pkg/errors"
	"time"

	"chat/apps/user/rpc/internal/svc"
	"chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号没有注册")
	ErrUserPwdError     = xerr.New(xerr.SERVER_COMMON_ERROR, "密码错误")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {

	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.WithStack(ErrPhoneNotRegister)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v , req %v", err, in.Phone)
	}

	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, errors.WithStack(ErrUserPwdError)
	}

	// 生成token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire,
		userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "ctxdata get jwt token err %v", err)
	}

	return &user.LoginResp{
		Id:     userEntity.Id,
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}

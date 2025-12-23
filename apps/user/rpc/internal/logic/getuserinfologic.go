package logic

import (
	"chat/apps/user/models"
	"chat/pkg/xerr"
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"chat/apps/user/rpc/internal/svc"
	"chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserNotFound = xerr.New(xerr.SERVER_COMMON_ERROR, "该用户不存在")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {

	userEntity, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.WithStack(ErrUserNotFound)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by id err: %v", err)
	}

	var resp user.UserEntity

	// 使用一个将结构体转化成另一个结构体的组件, 对字段进行映射，
	// Copy(&目标, &源)
	err = copier.Copy(&resp, userEntity)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "copy struct err: %v", err)
	}

	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}

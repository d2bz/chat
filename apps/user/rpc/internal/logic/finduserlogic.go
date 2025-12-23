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

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	var (
		userEntities []*models.Users
		err          error
	)

	if in.Phone != "" {
		userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err == nil {
			userEntities = append(userEntities, userEntity)
		}
	} else if in.Name != "" {
		userEntities, err = l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		userEntities, err = l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
	}

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "FindUser failed. in:%+v, err:%v", in, err)
	}
	var resp []*user.UserEntity
	// copier.Copy处理结构体切片，(&目标切片，&源切片)
	err = copier.Copy(&resp, &userEntities)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "copy struct err: %v", err)
	}

	return &user.FindUserResp{
		User: resp,
	}, nil
}

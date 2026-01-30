package group

import (
	"chat/apps/social/rpc/socialclient"
	"chat/pkg/ctxdata"
	"context"
	"strconv"

	"chat/apps/social/api/internal/svc"
	"chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请进群处理
func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleReq) (resp *types.GroupPutInHandleResp, err error) {
	uid := ctxdata.GetUId(l.ctx)

	// 转换GroupReqId从string到int32
	groupReqId, err := strconv.ParseInt(req.GroupReqId, 10, 32)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.Social.GroupPutInHandle(l.ctx, &socialclient.GroupPutInHandleReq{
		GroupReqId:   int32(groupReqId),
		GroupId:      req.GroupId,
		HandleUid:    uid,
		HandleResult: req.HandleResult,
	})

	return nil, err
}

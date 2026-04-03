package logic

import (
	"chat/apps/im/immodels"
	"chat/pkg/constants"
	"chat/pkg/xerr"
	"context"
	"github.com/pkg/errors"

	"chat/apps/im/rpc/im"
	"chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {
	res := &im.CreateGroupConversationResp{}

	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)
	// 说明已存在当前群聊会话
	if err == nil {
		return res, nil
	}
	if !errors.Is(err, immodels.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.FindOne err %v,conversationId %v", err, in.GroupId)
	}

	err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Insert err %v,conversationId %v", err, in.GroupId)
	}

	// 创建群聊与其创建者的会话
	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	return res, err
}

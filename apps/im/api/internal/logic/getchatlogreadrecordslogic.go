package logic

import (
	"chat/apps/im/rpc/im"
	"chat/apps/social/rpc/socialclient"
	"chat/apps/user/rpc/user"
	"chat/pkg/bitmap"
	"chat/pkg/constants"
	"context"

	"chat/apps/im/api/internal/svc"
	"chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogReadRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatLogReadRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogReadRecordsLogic {
	return &GetChatLogReadRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatLogReadRecordsLogic) GetChatLogReadRecords(req *types.GetChatLogReadRecordsReq) (resp *types.GetChatLogReadRecordsResp, err error) {
	// 获取对应的消息id
	chatLogs, err := l.svcCtx.Im.GetChatLog(l.ctx, &im.GetChatLogReq{
		MsgId: req.MsgId,
	})
	if err != nil || len(chatLogs.List) == 0 {
		return nil, err
	}

	var (
		chatLog = chatLogs.List[0]
		reads   = []string{chatLog.SendId}
		unReads []string
		ids     []string
	)

	// 判断并设置用户的已读未读情况
	switch constants.ChatType(chatLog.ChatType) {
	case constants.SingleChatType:
		if len(chatLog.ReadRecords) == 0 || chatLog.ReadRecords[0] == 0 {
			unReads = []string{chatLog.RecvId}
		} else {
			reads = append(reads, chatLog.RecvId)
		}
		ids = []string{chatLog.RecvId, chatLog.SendId}
	case constants.GroupChatType:
		groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
			GroupId: chatLog.RecvId,
		})
		if err != nil {
			return nil, err
		}

		bitmaps := bitmap.Load(chatLog.ReadRecords)
		for _, members := range groupUsers.List {
			ids = append(ids, members.UserId)

			if members.UserId == chatLog.SendId {
				continue
			}

			if bitmaps.IsSet(members.UserId) {
				reads = append(reads, members.UserId)
			} else {
				unReads = append(unReads, members.UserId)
			}
		}
	}

	userEntities, err := l.svcCtx.User.FindUser(l.ctx, &user.FindUserReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	userEntitySet := make(map[string]*user.UserEntity, len(userEntities.User))
	for i, entity := range userEntities.User {
		userEntitySet[entity.Id] = userEntities.User[i]
	}

	// 设置手机号码
	for i, read := range reads {
		if u := userEntitySet[read]; u != nil {
			reads[i] = u.Phone
		}
	}
	for i, unread := range unReads {
		if u := userEntitySet[unread]; u != nil {
			unReads[i] = u.Phone
		}
	}

	return &types.GetChatLogReadRecordsResp{
		Reads:   reads,
		UnReads: unReads,
	}, nil
}

-- 好友关系表
CREATE TABLE `friends` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户ID',
    `friend_uid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '好友UID',
    `remark` varchar(255) DEFAULT NULL COMMENT '好友备注',
    `add_source` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '添加来源（枚举值）',
    `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='好友关系表';

-- 好友申请记录表
CREATE TABLE `friend_requests` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请人ID',
    `req_uid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '被申请人UID',
    `req_msg` varchar(255) DEFAULT NULL COMMENT '申请留言',
    `req_time` timestamp NOT NULL COMMENT '申请时间',
    `handle_result` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '处理结果（1-通过 2-拒绝 0-未处理）',
    `handle_msg` varchar(255) DEFAULT NULL COMMENT '处理备注',
    `handled_at` timestamp NULL DEFAULT NULL COMMENT '处理时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='好友申请记录表';

-- 群组基础信息表
CREATE TABLE `groups` (
    `id` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群组ID（主键）',
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群组名称',
    `icon` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群组头像',
    `status` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '群组状态（1-正常 2-解散 3-封禁）',
    `creator_uid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '创建人UID',
    `group_type` int(11) NOT NULL COMMENT '群组类型（枚举值）',
    `is_verify` boolean NOT NULL COMMENT '是否需要验证加入',
    `notification` varchar(255) DEFAULT NULL COMMENT '群组公告',
    `notification_uid` varchar(64) DEFAULT NULL COMMENT '公告发布人UID',
    `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='群组基础信息表';

-- 群成员表
CREATE TABLE `group_members` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `group_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群组ID',
    `user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '成员用户ID',
    `role_level` tinyint COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色等级（1-群主 2-管理员 3-普通成员）',
    `join_time` timestamp NULL DEFAULT NULL COMMENT '加入时间',
    `join_source` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '加入来源（1-主动申请 2-邀请 3-扫码）',
    `inviter_uid` varchar(64) DEFAULT NULL COMMENT '邀请人UID',
    `operator_uid` varchar(64) DEFAULT NULL COMMENT '操作人UID（审批/邀请人）',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='群成员表';

-- 加群申请记录表
CREATE TABLE `group_requests` (
     `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
     `req_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请唯一标识',
     `group_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '目标群组ID',
     `req_msg` varchar(255) DEFAULT NULL COMMENT '申请留言',
     `req_time` timestamp NULL DEFAULT NULL COMMENT '申请时间',
     `join_source` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '申请来源（枚举值）',
     `inviter_user_id` varchar(64) DEFAULT NULL COMMENT '邀请人UID（邀请场景）',
     `handle_user_id` varchar(64) DEFAULT NULL COMMENT '处理人UID（群主/管理员）',
     `handle_time` timestamp NULL DEFAULT NULL COMMENT '处理时间',
     `handle_result` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '处理结果（1-通过 2-拒绝 0-未处理）',
     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='加群申请记录表';
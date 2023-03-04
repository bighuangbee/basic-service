package domain

import (
	"context"
)

type IOperationLogRepo interface {
	Add(ctx context.Context, oplog *OperationLog) error
	ListOperationLog(ctx context.Context, query *ListOperationLogRequest) ([]*OperationLog, int32, error)
	ListOperationLogUser(ctx context.Context, userName string) (user []*UserInfo, err error)
}

type OperationLog struct {
	ID int32 `json:"id"`
	// 应用ID
	AppID int32 `json:"appId"`
	// 企业ID
	CorpID int32 `json:"corpId"`
	// 企业名称
	CorpName string `json:"corpName"`
	// 操作人账号ID
	UserID int32 `json:"userId"`
	// 操作人姓名
	UserName string `json:"userName"`
	// 操作名称
	OperationName string `json:"operationName"`
	// 操作的模块
	OperationModule string `json:"operationModule"`
	// 变更前
	Before string `json:"before"`
	// 变更后
	After string `json:"after"`
	// 自定义详细内容
	Detail string `json:"detail"`
	// 操作时间
	//Timestamp *hiMysql.PBTime `json:"timestamp"`
	// 数据库创建时间
	//CreateAt hiMysql.PBTime `json:"creatAt"`
	// 状态
	Status string `json:"status"`
	// 请求失败原因
	Reason string `json:"reason"`
	// sn号
	SnNo string `json:"snNo"`
}

func (c *OperationLog) TableName() string {
	return "operation_log"
}


type ListOperationLogRequest struct {
	//Pagination     hiPagination.IPagination `json:"pagination,omitempty"`
	UserID         int32                    `json:"userId"`
	CorpID         int32                    `json:"corpId"`
	//OperateStartAt *hiMysql.PBTime          `json:"operateStartAt"`
	//OperateEndAt   *hiMysql.PBTime          `json:"operateEndAt"`
	AppID          int32                    `json:"appId"`
	OperationName  string                   `json:"operationName"`
	Status         []string                 `json:"status"`
}

type OperationLogUserResp struct {
	Items UserInfo `json:"items"`
}

type UserInfo struct {
	UserID   int32  `json:"userId"`
	UserName string `json:"userName"`
}

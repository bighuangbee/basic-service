package middleware

import (
	"context"
	"encoding/json"
	pbBasic "github.com/bighuangbee/basic-service/api/basic/v1"
	pb "github.com/bighuangbee/gokit/api/common/v1"
	"github.com/bighuangbee/gokit/model"
	"github.com/bighuangbee/gokit/tools"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"net"
	nhttp "net/http"
	"strconv"
	"strings"
)

//从proto文件解析router api=>title
var uriTitleMap = make(map[string]model.LogUrlInfoWithKey)

type OpLog struct {
	OpLogGrpcAddr 	string
	OpGrpcCli     	pbBasic.OperationLogClient
	ProtoPath		string
	Logger			log.Logger
}

func NewOpLog(protoPath, opLogGrpcAddr string, logger log.Logger) *OpLog {
	LoadOperationLogWithProto(protoPath, logger)

	return &OpLog{
		ProtoPath: protoPath,
		OpLogGrpcAddr: opLogGrpcAddr,
		Logger: logger,
		//OpGrpcCli: pbBasic.NewOperationLogClient(tools.GetGrpcClient(opLogGrpcAddr)),
	}
}

//后置，避免互依赖
func (r *OpLog) AfterStartOpLogGrpcConn()(func (ctx context.Context) error){
	return func(ctx context.Context) error {
		r.OpGrpcCli = pbBasic.NewOperationLogClient(tools.GetGrpcClient(r.OpLogGrpcAddr))
		r.Logger.Log(log.LevelInfo, "AfterStartOpLogGrpcConn", r.OpLogGrpcAddr)
		return nil
	}
}

func GetHostIp(transport interface{}) (ip string) {
	tr := transport.(*http.Transport)
	req := tr.Request()
	strs := strings.Split(req.Host, ":")
	if len(strs) > 0 {
		remoteIP := net.ParseIP(strs[0])
		if remoteIP == nil {
			return ""
		}
		return remoteIP.String()
	}

	return ""
}

func GetPathTemplateMethodIp(transport interface{}) (uri, method, ip string) {
	tr := transport.(*http.Transport)
	req := tr.Request()
	method = strings.ToUpper(req.Method)
	uri = tr.PathTemplate()
	ip = GetIP(req)
	return
}
func GetIP(r *nhttp.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}

	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return ""
	}
	remoteIP := net.ParseIP(ip)
	if remoteIP == nil {
		return ""
	}
	return remoteIP.String()
}

func (r *OpLog) Middleware() middleware.Middleware {

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {

				var (
					corpId   uint64
					userId   uint64
					userName string
					account  string
					bodyStr  string
				)
				uri, method, ip := GetPathTemplateMethodIp(tr)

				bs, err := json.Marshal(req)
				if err != nil {
					r.Logger.Log(log.LevelError, "NewJwtTokenByParse", err)
				} else {
					bodyStr = string(bs)
				}

				if uri == "/api/v1.0/account/login" {
					obj := make(map[string]interface{})
					json.Unmarshal([]byte(bodyStr), &obj)
					if _, ok := obj["password"]; ok {
						obj["password"] = "***"
					} else if _, ok := obj["token"]; ok {
						obj["token"] = "***"
					}
					if _, ok := obj["account"]; ok {
						account = obj["account"].(string)
					}
					// 更新body string
					bs, err := json.Marshal(obj)
					if err != nil {
						bodyStr = ""
					} else {
						bodyStr = string(bs)
					}
				}

				if token := tr.RequestHeader().Get("jwtToken"); token != "" {
					jwtToken, err := jwtToken.NewJwtTokenByParse(token, jwtTokenKey)
					if err != nil {
						log.NewHelper(t.Logger).Errorw("NewJwtTokenByParse err", err)
					} else {
						corpIdInt, _ := strconv.Atoi(jwtToken.Data["corpId"].(string))
						corpId = uint64(corpIdInt)
						userIdF := jwtToken.Data["userId"].(float64)
						userId = uint64(math.Ceil(userIdF))
						userName = jwtToken.Data["userName"].(string)
						account = jwtToken.Data["account"].(string)
					}
				}
				userAgent := tr.RequestHeader().Get("User-Agent")

				if v, ok := uriTitleMap[method+"_"+uri]; ok {
					details, _ := json.Marshal(&model.SysLogRpc{
						CorpId:    corpId,
						UserId:    uint32(userId),
						Name:      userName,
						LoginName: account,
						Type:      v.LogType,
						IP:        ip,
						Url:       uri,
						UrlTitle:  v.Title,
						Params:    bodyStr,
						UserAgent: userAgent,
					})
					operationName := getLogTypeStr(v.LogType)

					now := timestamppb.Now()
					if _, err := r.OpGrpcCli.Add(ctx, &pbBasic.AddRequest{
						Log: &pbBasic.Log{
							AppId:         0,
							CorpId:        int32(corpId),
							UserId:        int32(userId),
							UserName:      userName,
							OperationName: operationName,
							Detail:        string(details),
							Timestamp:     now,
							CreatAt:       now,
						},
					}); err != nil {
						r.Logger.Log(log.LevelError, "创建操作日志失败 OpGrpcCli.Add", err)
					}

					r.Logger.Log(log.LevelInfo, "method", method, "uri", uri, "title", v.LogType, "logType", v.LogType, "logType", operationName, "corpId", corpId, "userId", userId, "userName", userName, "account", account, "ip", ip, "payload", bodyStr)
				}
				return handler(ctx, req)
			}
			return nil, pb.ErrorUnauthenticated("CheckToken failed.")
		}
	}
}

func getLogTypeStr(logType model.LogType) string {
	switch logType {
	case model.LogDefault:
		return "未知操作类型"
	case model.LogCreate:
		return "创建"
	case model.LogView:
		return "查看"
	case model.LogEdit:
		return "编辑"
	case model.LogDelete:
		return "删除"
	case model.LogLogin:
		return "登陆"
	case model.LogLogout:
		return "登出"
	case model.LogExport:
		return "导出"
	case model.LogImport:
		return "导入"
	case model.LogSave:
		return "保存"
	}
	return "未知操作类型"
}

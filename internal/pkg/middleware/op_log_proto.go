package middleware

import (
	"fmt"
	"github.com/bighuangbee/gokit/model"
	"github.com/emicklei/proto"
	"github.com/go-kratos/kratos/v2/log"
	"io/ioutil"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
)

func LoadOperationLogWithProto(dir string, logger log.Logger) {
	//开发环境路径
	if dir == "../../config/proto/" {
		dir = "../../api"
	}
	filePaths := GetFiles(dir)
	for _, item := range filePaths {
		parseProtoFile(item, logger)
	}
	logger.Log(log.LevelError, "操作日志解析proto文件, 数量", len(filePaths))
}

func GetFiles(folder string) (filePaths []string) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		filename := file.Name()
		if file.IsDir() {
			subFilePaths := GetFiles(folder + "/" + filename)
			filePaths = append(filePaths, subFilePaths...)
		} else {
			if path.Ext(filename) == ".proto" {
				filePaths = append(filePaths, folder+"/"+filename)
			}
		}
	}
	return
}

func parseProtoFile(path string, logger log.Logger) {
	reader, err := os.Open(path)
	if err != nil {
		fmt.Println(string(debug.Stack()))
		panic(err)
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		fmt.Println(string(debug.Stack()))
		panic(err)
	}

	proto.Walk(definition,
		proto.WithService(func(s *proto.Service) {

			for _, e := range s.Elements {
				r, ok := e.(*proto.RPC)
				if ok {
					hasOpLog := false
					logUrlInfos := []model.LogUrlInfoWithKey{}
					for _, item := range r.Options {
						logTypeTitle := []string{}
						for _, itemConstant := range item.AggregatedConstants {
							logUrlInfo := model.LogUrlInfoWithKey{}
							if itemConstant.Name == "opLog" {
								// 该项包含操作日志
								if itemConstant.Literal != nil {
									logTypeTitle = strings.Split(itemConstant.Literal.Source, ",")
								}
								hasOpLog = true
							} else if itemConstant.Name == "get" || itemConstant.Name == "post" || itemConstant.Name == "put" || itemConstant.Name == "delete" {
								// get/ post/ put
								logUrlInfo.HttpMethod = convert2Method(itemConstant.Name)
								url := ""
								if itemConstant.Literal != nil {
									url = itemConstant.Literal.Source
								}
								logUrlInfo.Key = fmt.Sprintf("%s_%s", strings.ToUpper(itemConstant.Name), url)
								logUrlInfos = append(logUrlInfos, logUrlInfo)

								logger.Log(log.LevelDebug, "HTTP method", strings.ToUpper(itemConstant.Name), "url", url)
							}
						}

						if hasOpLog {
							for _, item := range logUrlInfos {
								if len(logTypeTitle) == 2 {
									urlLogType, _ := strconv.Atoi(logTypeTitle[0])
									item.LogType = model.LogType(urlLogType)
									item.Title = logTypeTitle[1]
									uriTitleMap[item.Key] = item
								}
							}
						}
					}
				}
			}
		}),
	)
}

func convert2Method(method string) model.HttpMethod {
	switch strings.ToUpper(method) {
	case "GET":
		return model.HttpGet
	case "POST":
		return model.HttpPost
	case "PUT":
		return model.HttpPut
	case "DELETE":
		return model.HttpDelete
	default:
		return model.HttpPost
	}
}

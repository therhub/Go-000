package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"therhub/wk02/src/service"
)

// api层
func Query(w http.ResponseWriter, req *http.Request) {

	// 记录所有日志和详细日志
	_, err := service.GetQueryResult()

	if err != nil {
		log.Print(err)
	}

	// 判断日志的具体类型，并在业务上做转换
	// 比如设置约定的业务级错误码
	// 判断日志详见dal_test
	if errors.Is(err, sql.ErrNoRows) {
	}
}

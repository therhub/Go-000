package service

import (
	"fmt"
	"therhub/wk02/src/dal"
	"therhub/wk02/src/model"
)

func GetQueryResult() ([]*model.GoCamp, error) {

	// 请求数据
	r, err := dal.Handler.Query(nil)

	if err != nil {
		return nil, fmt.Errorf("service.GetQueryResult,err:%w", err)
	}

	return r, err
}

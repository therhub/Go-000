package dal

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"therhub/wk02/src/model"

	_ "github.com/go-sql-driver/mysql"
)

func Test_handler_Query(t *testing.T) {

	tests := []struct {
		name string
		h    *handler
	}{
		// TODO: Add test cases.
		{name: "test", h: &handler{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// 初始化
			tt.h.initFunc()

			// ++++++++++++++++
			// 没有附加任何信息的情况下
			// 检索
			_, err := tt.h.Query(nil)

			// 检测错误第一种方式直接等于
			if err != sql.ErrNoRows {
				t.Error(err)
			}

			// ++++++++++++++++
			_, err = tt.h.Query(selfWrapAddMessage)

			// 经过组装的error不再等于sql.ErrNoRows
			if err == sql.ErrNoRows {
				t.Error(err)
			}

			// 但是组装后的数据仍然可以通过
			// errors.Is或者As来判断，这样就将最底层的err穿透到最表层
			// 并且error是一直在message中附加其他的信息
			if errors.Is(err, sql.ErrNoRows) == false {
				t.Error(err)
			}

			// ++++++++++++++++
			_, err = tt.h.Query(selfWrapAddEasyError)

			// 自定义错误
			var errObj *model.EasyError

			if errors.As(err, &errObj) == false {
				t.Error(err)
			}

			// ++++++++++++++++
			// 需要将具体
		})
	}
}

// 一般错误输出
func selfWrapAddMessage(err error) error {
	return fmt.Errorf("dal.err:%w", err)
}

// 自定义错误可以在业务上进行通过as来判断是否是这个错误
func selfWrapAddEasyError(err error) error {
	return &model.EasyError{Err: errors.New(err.Error()), Msg: "db.dal.err:"}
}

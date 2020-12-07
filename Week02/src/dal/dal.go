package dal

import (
	"database/sql"
	"therhub/wk02/src/model"

	_ "github.com/go-sql-driver/mysql"
)

// 构建handler
type handler struct {
	mDbHandler *sql.DB
}

// 初始化缓存handler
var Handler = &handler{}

// 初始化
func init() {
	Handler.initFunc()
}

// 初始func
func (h *handler) initFunc() {
	db, err := sql.Open("mysql", "admin:Admin123@tcp(localhost:3306)/gocamp")
	if err != nil {
		panic(err)
	}

	h.mDbHandler = db
}

// 检索
func (h *handler) Query(f func(error) error) ([]*model.GoCamp, error) {

	// 执行
	rows := h.mDbHandler.QueryRow("select Pid from WK02 where pid=1;")

	// 定义返回结果
	r := make([]*model.GoCamp, 0, 16)

	// 检索结果
	err := rows.Scan(&r)
	if err != nil {
		if f != nil {
			err = f(err)
		}

		return nil, err
	}

	return r, nil
}

package login

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// 验证是否实现了LoginDaler
var _ LoginDaler = &gromLoginDal{}

type gromLoginDal struct {
	handler *gorm.DB
}

func NewgromLoginDal(db *gorm.DB) *gromLoginDal {
	return &gromLoginDal{handler: db}
}

func (l *gromLoginDal) Query(u *loginUser) error {

	s := fmt.Sprintf("select Id from wk04_bak where Id=%v and Password='%v';", u.Id, u.Password)

	// 执行sql
	rows := l.handler.DB().QueryRow(s)

	// 定义返回结果
	r := make([]*loginUser, 0, 1)

	// 遍历检查是否存在用户user
	err := rows.Scan(&r)

	// 查询出错误
	if err != nil {

		// 非未查到
		if errors.Is(err, sql.ErrNoRows) == false {
			return errors.Wrap(errors.WithStack(errors.New(fmt.Sprintf("%w", ErrSystemStatus))), fmt.Sprint("查询发生错误,sql:%s,堆栈信息", s))
		}

		// 未查到
		return errors.New(fmt.Sprintf("%w", NoQueryStaus))
	}

	return nil
}

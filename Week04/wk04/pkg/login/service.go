package login

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"google.golang.org/grpc"
	v1 "therhub.com/wk04/api/v1"
)

var (
	dbHandler *gorm.DB
	// 游戏模型数据库对象
	dbHandler2 *sql.DB
)

func init() {

	// sql应该有config统一提供，
	sqlconn := "admin:Admin123@tcp(localhost:3306)/gocamp"
	gormSqlconn := "admin:Admin123@tcp(localhost:3306)/gocamp_bak"

	gormDb, err := gorm.Open("mysql", gormSqlconn)
	if err != nil {
		panic(fmt.Errorf("初始化gorm数据库:%s失败，错误信息为：%s", gormSqlconn, err))
	}

	db, err := sql.Open("mysql", sqlconn)
	if err != nil {
		panic(fmt.Errorf("初始化sql数据库:%s失败，错误信息为：%s", sqlconn, err))
	}

	dbHandler = gormDb
	dbHandler2 = db
}

// 定义自定义的Mux对象
type LoginService struct {
	v1.UnimplementedLoginGrpcServiceServer
	mgr *loginMgr
}

// 新建1个service
func NewLoginService(m *loginMgr) *LoginService {
	return &LoginService{mgr: m}
}

// 初始化dbhandler
// 放于外层可以方面测试，
func NewDbHandler() *gorm.DB {
	return dbHandler
}

// 初始化dbhandler
// 放于外层可以方面测试，
func NewDbHandler2() *sql.DB {
	return dbHandler2
}

// 转换为业务user
func NewloginUser(u *v1.User) *loginUser {

	// md5加密
	pwd := md5.New().Sum([]byte(u.Password))

	return &loginUser{Id: u.Id, Name: u.Name, Password: fmt.Sprintf("%x", pwd)}
}

// Login实现
func (s *LoginService) Login(ctx context.Context, in *v1.User) (*v1.LoginResponse, error) {

	err := s.mgr.handlerLogin(NewloginUser(in))

	if err != nil {
		if errors.Is(err, ErrSystemErr{}) {
			log.Println(fmt.Printf("%w", err))
			return &v1.LoginResponse{Status: ErrSystemStatus}, ErrSystemErr{}
		} else if errors.Is(err, ErrNoQuery{}) {
			return &v1.LoginResponse{Status: NoQueryStaus}, ErrNoQuery{}
		}
	}

	return &v1.LoginResponse{Status: Success}, nil
}

// 注册grpc
func RegisterFunc(rp *grpc.Server) {
	// 注册grpc
	v1.RegisterLoginGrpcServiceServer(rp, InitializeLoginService())
}

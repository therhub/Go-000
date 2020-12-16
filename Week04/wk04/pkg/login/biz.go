package login

// 登陆管理mgr
type loginMgr struct {
	dal LoginDaler
}

// 新登陆
func NewloginMgr(l *loginDal) *loginMgr {
	return &loginMgr{dal: l}
}

func (l *loginMgr) handlerLogin(user *loginUser) error {
	return l.dal.Query(user)
}

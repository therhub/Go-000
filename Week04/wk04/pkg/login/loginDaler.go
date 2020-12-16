package login

// 解耦合
type LoginDaler interface {
	Query(*loginUser) error
}

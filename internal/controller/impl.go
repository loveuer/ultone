package controller

var (
	// UserController todo: 可以实现自己的 controller
	UserController userController
)

func Init() error {
	UserController = uc{}

	return nil
}

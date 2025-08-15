package runner

type Service struct {
	ctrl *Controller
}

func NewService(ctrl *Controller) Service {
	return Service{ctrl: ctrl}
}

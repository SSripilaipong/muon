package runner

import "github.com/SSripilaipong/muon/server/coordinator"

type Service struct {
	ctrl  *Controller
	coord *coordinator.Controller
}

func NewService(ctrl *Controller) *Service {
	return &Service{ctrl: ctrl}
}

func (s *Service) SetCoordinator(coord *coordinator.Controller) {
	s.coord = coord
}

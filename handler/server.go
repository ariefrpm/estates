package handler

import (
	"github.com/SawitProRecruitment/EstateService/core/interfaces"
)

type Server struct {
	estateUsecase interfaces.EstateUsecase
}

func NewServer(estateUsecase interfaces.EstateUsecase) *Server {
	return &Server{
		estateUsecase: estateUsecase,
	}
}

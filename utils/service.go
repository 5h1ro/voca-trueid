package utils

import (
	"strings"
	"vocatrueid/entity"
)

type Service interface {
	CheckNickname(body entity.User, code string) entity.CheckNicknameResponse
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CheckNickname(body entity.User, code string) entity.CheckNicknameResponse {
	var userId, zoneId string
	splitIndex := strings.Index(body.Target, "|")
	if splitIndex >= 0 {
		userId = strings.Split(body.Target, "|")[0]
		zoneId = strings.Split(body.Target, "|")[1]
	} else {
		userId = body.Target
	}
	return s.repository.GetUsernameByCode(userId, zoneId, code)
}

package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

func Profile2DTO(user *core.User) dto.UserProfile {
	return dto.UserProfile{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Avatar:   user.Image,
		Phone:    user.Phone,
		Location: user.Location,
		BirthDay: user.BirthDay,
	}
}

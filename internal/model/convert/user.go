package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

//func User2Core(user *dto.User) core.User {
//	return core.User{
//		ID:    user.ID,
//		Email: user.Email,
//		Name:  user.Name,
//	}
//}

func User2DTO(user *core.User) dto.User {
	return dto.User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Image: user.Image,
	}
}
func User2author(user dto.User) dto.Author {
	return dto.Author{
		ID:    user.ID,
		Name:  user.Name.Full(),
		Image: user.Image,
		Type:  "User",
	}
}

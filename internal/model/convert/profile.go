package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

func Profile2DTO(user *core.User, friendIDs []string) dto.UserProfile {
	return dto.UserProfile{
		UserInfo: dto.User{ID: user.ID,
			Email: user.Email,
			Name:  user.Name},
		Avatar:    user.Image,
		Phone:     user.Phone,
		Location:  user.Location,
		BirthDay:  user.BirthDay,
		FriendIDs: friendIDs,
		PostIDs:   user.Posts,
	}
}

func EditProfile2Core(newProfile *dto.EditProfile) core.EditInfo {
	return core.EditInfo{
		Name:     newProfile.Name,
		Avatar:   newProfile.Avatar,
		Phone:    newProfile.Phone,
		Location: newProfile.Location,
		BirthDay: newProfile.BirthDay,
	}
}

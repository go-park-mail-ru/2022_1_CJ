package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

func Community2DTOprofile(community *core.Community, admins []dto.User) dto.CommunityProfile {
	return dto.CommunityProfile{
		ID:        community.ID,
		Name:      community.Name,
		Image:     community.Image,
		Info:      community.Info,
		Followers: int64(len(community.FollowerIDs)),
		Admins:    admins,
	}
}

func Community2DTOSmallProfile(community *core.Community) dto.CommunityProfile {
	return dto.CommunityProfile{
		ID:    community.ID,
		Name:  community.Name,
		Image: community.Image,
	}
}

func Community2DTO(community *core.Community) dto.Community {
	return dto.Community{
		ID:    community.ID,
		Name:  community.Name,
		Image: community.Image,
	}
}

func CommunityProfile2Author(community dto.CommunityProfile) dto.Author {
	return dto.Author{
		ID:    community.ID,
		Name:  community.Name,
		Image: community.Image,
		Type:  "Community",
	}
}

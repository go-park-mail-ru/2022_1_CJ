package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

func Like2Core(like *dto.Like) core.Like {
	return core.Like{
		Amount:  like.Amount,
		UserIDs: like.UserIDs,
	}
}

func isMyLike(userIDs []string, userID string) bool {
	for _, id := range userIDs {
		if id == userID {
			return true
		}
	}
	return false
}

func Like2DTO(like *core.Like, userID string) dto.Like {
	return dto.Like{
		Amount:  like.Amount,
		UserIDs: like.UserIDs,
		MyLike:  isMyLike(like.UserIDs, userID),
	}
}

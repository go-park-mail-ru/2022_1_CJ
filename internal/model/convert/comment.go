package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

func Comment2DTO(post *core.Comment, user *core.User) dto.Comment {
	return dto.Comment{
		ID:        post.ID,
		Message:   post.Message,
		Images:    post.Images,
		Author:    User2DTO(user),
		CreatedAt: post.CreatedAt,
	}
}

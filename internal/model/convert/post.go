package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

func Post2Core(post *dto.Post) core.Post {
	return core.Post{
		AuthorID: post.AuthorID,
		Message:  post.Message,
		Images:   post.Images,
	}
}

func Post2DTO(post *core.Post) dto.Post {
	return dto.Post{
		AuthorID: post.AuthorID,
		Message:  post.Message,
		Images:   post.Images,
	}
}

package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
)

//func Post2Core(post *dto.Post) core.Post {
//	return core.Post{
//		ID:       post.ID,
//		AuthorID: post.Author.ID,
//		Message:  post.Message,
//		Images:   post.Images,
//	}
//}

func Post2DTOByUser(post *core.Post, author *core.User) dto.Post {
	return dto.Post{
		ID:            post.ID,
		Author:        User2author(User2DTO(author)),
		Message:       post.Message,
		Images:        post.Images,
		CreatedAt: post.CreatedAt,
    CountComments: int64(len(post.CommentsIDs)),
	}
}

func Post2DTOByCommunity(post *core.Post, community *core.Community) dto.Post {
	return dto.Post{
		ID:            post.ID,
		Author:        CommunityProfile2Author(Community2DTOSmallProfile(community)),
		Message:       post.Message,
		Images:        post.Images,
		CreatedAt: post.CreatedAt,
    CountComments: int64(len(post.CommentsIDs)),
	}
}

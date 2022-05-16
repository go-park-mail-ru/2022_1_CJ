package convert

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommunity2DTOprofile(t *testing.T) {
	admins := []dto.User{{ID: "123"}, {ID: "1235"}}
	community := &core.Community{ID: "12312312", Name: "nameCool", FollowerIDs: []string{"2131", "231"}}
	res := convert.Community2DTOprofile(community, admins)
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, res, dto.CommunityProfile{ID: "12312312", Name: "nameCool", Followers: 2, Admins: admins}) {
			t.Error("got : ", res, " expected :", dto.CommunityProfile{ID: "12312312", Name: "nameCool", Followers: 2, Admins: admins})
		}
	})
}

func TestCommunity2DTOSmallProfile(t *testing.T) {
	community := &core.Community{ID: "12312312", Name: "nameCool", FollowerIDs: []string{"2131", "231"}}
	res := convert.Community2DTOSmallProfile(community)
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, res, dto.CommunityProfile{ID: "12312312", Name: "nameCool"}) {
			t.Error("got : ", res, " expected :", dto.CommunityProfile{ID: "12312312", Name: "nameCool"})
		}
	})
}

func TestCommunity2DTO(t *testing.T) {
	community := &core.Community{ID: "12312312", Name: "nameCool", FollowerIDs: []string{"2131", "231"}}
	res := convert.Community2DTO(community)
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, res, dto.Community{ID: "12312312", Name: "nameCool"}) {
			t.Error("got : ", res, " expected :", dto.Community{ID: "12312312", Name: "nameCool"})
		}
	})
}

func TestCommunityProfile2Author(t *testing.T) {
	community := &dto.CommunityProfile{ID: "12312312", Name: "nameCool"}
	res := convert.CommunityProfile2Author(*community)
	t.Run("Check equals", func(t *testing.T) {
		if !assert.Equal(t, res, dto.Author{ID: "12312312", Name: "nameCool", Type: "Community"}) {
			t.Error("got : ", res, " expected :", dto.Author{ID: "12312312", Name: "nameCool", Type: "Community"})
		}
	})
}

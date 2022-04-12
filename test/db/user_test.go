package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	testdb, err := TestConnectToDB(t)
	if err != nil {
		return
	}

	ctx := context.Background()
	u := TestUser(t)

	assert.NoError(t, testdb.UserRepo.CreateUser(ctx, u))
	assert.NotNil(t, u.ID)
	assert.NotNil(t, u.FriendsID)
	assert.NotNil(t, u.CreatedAt)
}

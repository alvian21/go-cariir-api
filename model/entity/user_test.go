package entity

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserEntity(t *testing.T) {
	id := uuid.New()
	user := User{
		ID:       id,
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
		IsActive: true,
	}

	assert.Equal(t, id, user.ID)
	assert.Equal(t, "Test User", user.FullName)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.Password)

	// Test JSON Marshalling (Password should be hidden)
	data, err := json.Marshal(user)
	assert.NoError(t, err)

	var userMap map[string]interface{}
	err = json.Unmarshal(data, &userMap)
	assert.NoError(t, err)

	assert.Equal(t, "Test User", userMap["full_name"])
	assert.Equal(t, "test@example.com", userMap["email"])
	assert.NotContains(t, userMap, "password")
}

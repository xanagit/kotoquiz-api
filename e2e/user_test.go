package main

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xanagit/kotoquiz-api/models"
	"net/http"
	"strconv"
	"testing"
)

func Test_should_create_user(t *testing.T) {
	t.Parallel()

	user := GenerateUser()
	var insertedUser models.User
	httpResCode := post("/api/v1/tech/users", ToJson(&user), &insertedUser)

	assert.Equal(t, http.StatusCreated, httpResCode)
	assert.NotEqual(t, user.ID, insertedUser.ID)
	assert.Equal(t, "", insertedUser.Password)
	assert.Equal(t, user.Email, insertedUser.Email)
	assert.Equal(t, user.Username, insertedUser.Username)
}

func Test_should_read_user(t *testing.T) {
	t.Parallel()
	var httpResCode int

	user := GenerateUser()
	user.ID = uuid.Nil
	var insertedUser models.User
	httpResCode = post("/api/v1/tech/users", ToJson(&user), &insertedUser)

	assert.Equal(t, http.StatusCreated, httpResCode)
	var fetchedUser models.User
	httpResCode = get("/api/v1/tech/users/"+insertedUser.ID.String(), &fetchedUser)

	assert.Equal(t, http.StatusOK, httpResCode)
	insertedUser.Password = "" // Password is not returned in response
	assert.Equal(t, insertedUser, fetchedUser)
}

func Test_should_update_user(t *testing.T) {
	t.Parallel()

	user := GenerateUser()
	user.ID = uuid.Nil
	var insertedUser models.User
	post("/api/v1/tech/users", ToJson(&user), &insertedUser)

	var updatedUser models.User
	user.ID = uuid.New()
	user.Email = "updated@example.com"
	user.Username = "updateduser"
	httpResCode := put("/api/v1/tech/users/"+insertedUser.ID.String(), ToJson(&user), &updatedUser)

	assert.Equal(t, http.StatusOK, httpResCode)
	assert.NotEqual(t, user.ID, updatedUser.ID)
	assert.Equal(t, insertedUser.ID, updatedUser.ID)
	assert.NotEqual(t, user.ID, insertedUser.ID)
	assert.Equal(t, "", updatedUser.Password)
	assert.Equal(t, user.Email, updatedUser.Email)
	assert.Equal(t, user.Username, updatedUser.Username)
}

func Test_should_delete_user(t *testing.T) {
	t.Parallel()

	user := GenerateUser()
	user.ID = uuid.Nil
	var insertedUser models.User
	post("/api/v1/tech/users", ToJson(&user), &insertedUser)

	httpResCode := del("/api/v1/tech/users/" + insertedUser.ID.String())

	assert.Equal(t, http.StatusNoContent, httpResCode)
}

func Test_should_not_create_user_with_duplicate_email(t *testing.T) {
	t.Parallel()

	users := []models.User{GenerateUser(), GenerateUser()}
	users[1].Email = users[0].Email // Set same email

	var insertedUser models.User
	httpResCode := post("/api/v1/tech/users", ToJson(&users[0]), &insertedUser)
	assert.Equal(t, http.StatusCreated, httpResCode)

	httpResCode = post("/api/v1/tech/users", ToJson(&users[1]), &insertedUser)
	assert.Equal(t, http.StatusInternalServerError, httpResCode)
}

func Test_should_create_multiple_users(t *testing.T) {
	t.Parallel()
	var httpResCode int

	users := []models.User{GenerateUser(), GenerateUser(), GenerateUser()}

	insertedUsers := make([]models.User, 3)
	for idx, user := range users {
		user.Email = "test" + strconv.Itoa(idx) + "@example.com"
		user.Username = "user" + strconv.Itoa(idx)
		httpResCode = post("/api/v1/tech/users", ToJson(&user), &insertedUsers[idx])
		assert.Equal(t, http.StatusCreated, httpResCode)
	}

	// Verify each user was created correctly
	for _, user := range insertedUsers {
		var fetchedUser models.User
		httpResCode = get("/api/v1/tech/users/"+user.ID.String(), &fetchedUser)
		assert.Equal(t, http.StatusOK, httpResCode)
		user.Password = "" // Password is not returned in response
		assert.Equal(t, user, fetchedUser)
	}
}

func GenerateUser() models.User {
	return models.User{
		ID:       uuid.New(),
		Email:    uuid.New().String() + "@example.com",
		Username: "testuser",
		Password: "password123",
	}
}

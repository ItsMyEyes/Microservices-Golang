package repository

import (
	"log"
	"pakawai_service/cmd/auth/model"
	"pakawai_service/configs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUsersRepositorySave(t *testing.T) {
	user := &model.User{
		Id:       primitive.NewObjectID(),
		Name:     "TEST",
		Email:    "test@mail.com",
		Password: "123456789",
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	r := NewUserRepository(*configs.Client)
	err := r.Create(user)
	assert.NoError(t, err)

	found, err := r.GetById(user.Id)
	assert.NoError(t, err)
	assert.NotNil(t, found)
}

func TestUserRepositoryFind(t *testing.T) {
	user := &model.User{
		Name:  "TEST",
		Email: "test@mail.com",
	}

	r := NewUserRepository(*configs.Client)
	found, err := r.GetByEmail(user.Email)
	log.Printf("found: %v", found)
	assert.NoError(t, err)

	assert.Equal(t, user.Name, found.Name)
	assert.Equal(t, user.Email, found.Email)

}

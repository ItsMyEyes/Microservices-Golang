package repository

import (
	"context"
	"pakawai_service/cmd/auth/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const UserCollection = "users"

type UserRepository interface {
	Create(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetById(id primitive.ObjectID) (*model.User, error)
	Update(user *model.User) error
	Delete(id primitive.ObjectID) error
	GetAll() ([]*model.User, error)
	ValidUser(id primitive.ObjectID) bool
}

type usersRepository struct {
	c *mongo.Collection
}

func NewUserRepository(conn mongo.Client) UserRepository {
	return &usersRepository{
		c: conn.Database("pakawai").Collection(UserCollection),
	}
}

func (r *usersRepository) Create(user *model.User) error {
	_, err := r.c.InsertOne(context.Background(), user)
	return err
}

func (r *usersRepository) GetById(id primitive.ObjectID) (user *model.User, err error) {
	r.c.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	return user, err
}

func (r *usersRepository) ValidUser(id primitive.ObjectID) bool {
	var user model.User
	r.c.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	return user.Id != primitive.NilObjectID
}

func (r *usersRepository) GetByEmail(email string) (user *model.User, err error) {
	r.c.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	return user, err
}

func (r *usersRepository) GetAll() (users []*model.User, err error) {
	var user model.User
	cursor, err := r.c.Find(context.Background(), bson.D{})
	if err != nil {
		defer cursor.Close(context.Background())
		return users, err
	}

	for cursor.Next(context.Background()) {
		err := cursor.Decode(&user)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *usersRepository) Update(user *model.User) error {
	_, err := r.c.UpdateByID(context.Background(), user.Id, user)
	return err
}

func (r *usersRepository) Delete(id primitive.ObjectID) error {
	_, err := r.c.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *usersRepository) DeleteAll() error {
	return r.c.Drop(context.Background())
}

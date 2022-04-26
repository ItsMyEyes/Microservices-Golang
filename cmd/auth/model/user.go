package model

import (
	pb "pakawai_service/common/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	Created  time.Time          `bson:"created"`
	Updated  time.Time          `bson:"updated"`
}

func (u *User) ToProtoBuffer() *pb.User {
	return &pb.User{
		Id:       u.Id.Hex(),
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (u *User) FromProtoBuffer(user *pb.User) {
	oid, _ := primitive.ObjectIDFromHex(user.GetId())
	u.Id = oid
	u.Name = user.GetName()
	u.Email = user.GetEmail()
	u.Password = user.GetPassword()
}

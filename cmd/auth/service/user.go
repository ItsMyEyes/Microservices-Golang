package service

import (
	"context"
	"log"
	"pakawai_service/cmd/auth/model"
	"pakawai_service/cmd/auth/repository"
	"pakawai_service/cmd/auth/validator"
	"strings"
	"time"

	"pakawai_service/cmd/auth/security"
	pb "pakawai_service/common/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/peer"
)

type authService struct {
	usersRepository repository.UserRepository
}

func NewAuthService(usersRepository repository.UserRepository) pb.AuthServiceServer {
	return &authService{usersRepository: usersRepository}
}

func (s *authService) SignUp(ctx context.Context, req *pb.User) (*pb.User, error) {
	err := validator.ValidateSignUp(req)
	if err != nil {
		return nil, err
	}

	req.Password, err = security.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = validator.NormalizeEmail(req.Email)

	found, err := s.usersRepository.GetByEmail(req.Email)
	log.Printf("found: %v", found)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if found == nil {
		user := new(model.User)
		user.FromProtoBuffer(req)
		err := s.usersRepository.Create(user)
		if err != nil {
			return nil, err
		}
		return user.ToProtoBuffer(), nil
	}

	p, _ := peer.FromContext(ctx)
	log.Printf("from: %v someone hit signup service", p.Addr.String())
	return nil, validator.ErrEmailAlreadyExists
}

func (s *authService) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	req.Email = validator.NormalizeEmail(req.Email)
	user, _ := s.usersRepository.GetByEmail(req.Email)
	log.Printf("user: %v", user)
	if user == nil {
		log.Println("signin failed:")
		return nil, validator.ErrSignInFailed
	}

	auth := security.VerifyPassword(req.Password, user.Password)
	if !auth {
		log.Println("signin failed:", validator.ErrSignInFailed.Error())
		return nil, validator.ErrSignInFailed
	}

	token, err := security.MakeJWT(user)
	if err != nil {
		log.Println("signin failed:", err.Error())
		return nil, validator.ErrSignInFailed
	}

	p, _ := peer.FromContext(ctx)
	log.Printf("from: %v someone hit signin service", p.Addr.String())
	return &pb.SignInResponse{User: user.ToProtoBuffer(), Token: token}, nil
}

func (s *authService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	oid, _ := primitive.ObjectIDFromHex(req.Id)
	if !primitive.IsValidObjectID(req.Id) {
		return nil, validator.ErrInvalidUserId
	}
	found, err := s.usersRepository.GetById(oid)
	if err != nil {
		return nil, err
	}
	return found.ToProtoBuffer(), nil
}

func (s *authService) ListUsers(req *pb.ListUsersRequest, stream pb.AuthService_ListUsersServer) error {
	users, err := s.usersRepository.GetAll()
	if err != nil {
		return err
	}
	for _, user := range users {
		err := stream.Send(user.ToProtoBuffer())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *authService) UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	oid, _ := primitive.ObjectIDFromHex(req.Id)
	if !primitive.IsValidObjectID(req.Id) {
		return nil, validator.ErrInvalidUserId
	}
	user, err := s.usersRepository.GetById(oid)
	if err != nil {
		return nil, err
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, validator.ErrEmptyName
	}
	if req.Name == user.Name {
		return user.ToProtoBuffer(), nil
	}

	user.Name = req.Name
	user.Updated = time.Now()
	err = s.usersRepository.Update(user)
	return user.ToProtoBuffer(), err
}

func (s *authService) DeleteUser(ctx context.Context, req *pb.GetUserRequest) (*pb.DeleteUserResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(req.Id)
	if !primitive.IsValidObjectID(req.Id) {
		return nil, validator.ErrInvalidUserId
	}
	err := s.usersRepository.Delete(oid)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{Id: req.Id}, nil
}

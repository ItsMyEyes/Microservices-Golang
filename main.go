package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"pakawai_service/cmd/auth/repository"
	"pakawai_service/cmd/auth/service"
	"pakawai_service/configs"
	_ "pakawai_service/configs"
	"strconv"

	pb "pakawai_service/common/model"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

var (
	local bool
	port  int
)

func init() {
	port = 3000
	if os.Getenv("PORT_AUTH") != "" {
		port, _ = strconv.Atoi(os.Getenv("PORT_AUTH"))
	}
	flag.IntVar(&port, "port", port, "authentication service port")
	flag.BoolVar(&local, "local", true, "run authentication service local")
	flag.Parse()
}

func main() {
	if local {
		err := godotenv.Load()
		if err != nil {
			log.Panicln(err)
		}
	}

	userRepository := repository.NewUserRepository(*configs.Client)
	authService := service.NewAuthService(userRepository)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authService)

	log.Printf("Authentication service running on 0.0.0.0:%d\n", port)

	grpcServer.Serve(lis)
}

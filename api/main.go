package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"pakawai_service/api/resthandlers"
	"pakawai_service/api/routes"
	pb "pakawai_service/common/model"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

var (
	port     int
	authAddr string
)

func init() {
	port := 8080
	port_service := os.Getenv("PORT_AUTH")
	if os.Getenv("PORT_CLIENT") != "" {
		port, _ = strconv.Atoi(os.Getenv("PORT_CLIENT"))
	}
	if port_service != "" {
		port_service = strconv.Itoa(3000)
	}
	flag.IntVar(&port, "port", port, "api service port")
	flag.StringVar(&authAddr, "auth_addr", port_service, "authenticaton service address")
	flag.Parse()
}

func main() {

	conn, err := grpc.Dial(authAddr, grpc.WithInsecure())
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	authSvcClient := pb.NewAuthServiceClient(conn)
	authHandlers := resthandlers.NewAuthHandlers(authSvcClient)
	authRoutes := routes.NewAuthRoutes(authHandlers)

	router := mux.NewRouter().StrictSlash(true)
	routes.Install(router, authRoutes)

	log.Printf("API service running on [::]:%d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), routes.WithCORS(router)))
}

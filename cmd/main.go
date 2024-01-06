package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/akshay0074700747/orders-service/db"
	"github.com/akshay0074700747/orders-service/initializer"
	"github.com/akshay0074700747/orders-service/service"
	"github.com/akshay0074700747/proto-files-for-microservices/pb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err.Error())
	}

	addr := os.Getenv("DATABASE_ADDR")

	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	servicee := initializer.Initialize(DB)

	server := grpc.NewServer()

	pb.RegisterOrderServiceServer(server, servicee)

	listener, err := net.Listen("tcp", ":50003")
	if err != nil {
		log.Fatal(err.Error())
	}

	userConn, err := grpc.Dial("user-service:50002", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err.Error())
	}

	productConn, err := grpc.Dial("product-service:50004", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err.Error())
	}

	userRes := pb.NewUserServiceClient(userConn)
	productRes := pb.NewProductServiceClient(productConn)

	service.Initializer(userRes, productRes)

	log.Println("Order Server is running...")

	if err := server.Serve(listener); err != nil {
		log.Fatal(err.Error())
	}
}

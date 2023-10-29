package main

import (
	"context"
	"log"

	pb "github.com/Sergio10RPR/grpcClient/proto" // Importa el paquete generado desde el archivo .proto

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

var ctx = context.Background()

func insertData(c *fiber.Ctx) error {

	var data map[string]string
	e := c.BodyParser(&data)
	if e != nil {
		return e
	}

	rank := &pb.User{
		Carnet:   data["carnet"],
		Nombre:   data["nombre"],
		Curso:    data["curso"],
		Nota:     data["nota"],
		Semestre: data["semestre"],
		Year:     data["year"],
	}

	err := sendGRPCRequest(rank)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Datos insertados en el servidor gRPC.",
	})
}

func sendGRPCRequest(user *pb.User) error {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	_, err = client.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	app := fiber.New()

	app.Post("/insertar", insertData)

	err := app.Listen(":3000")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "studentgrpc/proto"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr           = flag.String("addr", "go-server-service-swimming:50051", "the address to connect to")
	addr_athletism = flag.String("addr_athletism", "go-server-service-athletics:50052", "the address to connect to")
	addr_boxing    = flag.String("addr_boxing", "go-server-service-boxing:50053", "the address to connect to")
)

type Student struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Faculty    string `json:"faculty"`
	Discipline int    `json:"discipline"`
}

func sendData(fiberCtx *fiber.Ctx) error {
	var body Student
	if err := fiberCtx.BodyParser(&body); err != nil {
		return fiberCtx.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Seleccionar la dirección del servidor basado en la disciplina
	var serverAddr string
	switch body.Discipline {
	case 1:
		serverAddr = *addr // Natacion
	case 2:
		serverAddr = *addr_athletism // Atletismo
	case 3:
		serverAddr = *addr_boxing // Boxeo
	default:
		serverAddr = *addr // Valor por defecto
	}

	// Establecer conexión con el servidor gRPC adecuado
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return fiberCtx.Status(500).JSON(fiber.Map{
			"error": "did not connect to the server",
		})
	}
	defer conn.Close()
	c := pb.NewStudentClient(conn)

	// Crear un canal para recibir la respuesta y el error
	responseChan := make(chan *pb.StudentResponse)
	errorChan := make(chan error)
	go func() {
		// Contactar al servidor gRPC
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		r, err := c.GetStudent(ctx, &pb.StudentRequest{
			Name:       body.Name,
			Age:        int32(body.Age),
			Faculty:    body.Faculty,
			Discipline: pb.Discipline(body.Discipline),
		})

		if err != nil {
			errorChan <- err
			return
		}

		responseChan <- r
	}()

	// Esperar la respuesta o error
	select {
	case response := <-responseChan:
		return fiberCtx.JSON(fiber.Map{
			"message": response.GetSuccess(),
		})
	case err := <-errorChan:
		return fiberCtx.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	case <-time.After(5 * time.Second):
		return fiberCtx.Status(500).JSON(fiber.Map{
			"error": "timeout",
		})
	}
}

func main() {
	flag.Parse() // Necesario para usar flag en el main

	app := fiber.New()
	app.Post("/agronomia", sendData)

	err := app.Listen(":8080")
	if err != nil {
		log.Println(err)
		return
	}
}

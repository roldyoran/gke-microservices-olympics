package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	pb "studentgrpc/proto"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

var (
	port           = flag.Int("port", 50051, "The server port")
	kafkaBootstrap = "my-cluster-kafka-bootstrap:9092" // Dirección del clúster de Kafka
)

// Server is used to implement the gRPC server in the proto library
type server struct {
	pb.UnimplementedStudentServer
	writer *kafka.Writer
}

// Implement the GetStudent method
func (s *server) GetStudent(_ context.Context, in *pb.StudentRequest) (*pb.StudentResponse, error) {
	log.Printf("\nReceived: %v", in)
	log.Printf("Student name: %s", in.GetName())
	log.Printf("Student faculty: %s", in.GetFaculty())
	log.Printf("Student age: %d", in.GetAge())
	log.Printf("Student discipline: %d", in.GetDiscipline())

	// Determina si ganó o perdió
	resultado := rand.Intn(2) // Genera un número aleatorio entre 0 y 1

	var estado string
	var topic string

	// Definir el tópico basado en si el estudiante ganó o perdió
	if resultado == 1 {
		estado = "ganó"
		topic = "student-winners" // Tópico para ganadores
	} else {
		estado = "perdió"
		topic = "student-losers" // Tópico para perdedores
	}

	log.Printf("El estudiante: %s, %s", in.GetName(), estado)

	// Crear el mensaje a enviar a Kafka
	studentInfo := fmt.Sprintf(`{"name": "%s", "status": "%s", "faculty": "%s", "discipline": %d}`, in.GetName(), estado, in.GetFaculty(), in.GetDiscipline())

	// Generar una clave aleatoria
	randomKey := strconv.Itoa(rand.Intn(1000000)) // Generar un número aleatorio como clave

	// Enviar el mensaje a Kafka
	err := s.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(randomKey), // Usar la clave aleatoria
			Value: []byte(studentInfo),
			Topic: topic,
		},
	)

	if err != nil {
		log.Printf("Error al enviar mensaje a Kafka: %v", err)
		return nil, err
	}

	log.Printf("Mensaje enviado a Kafka (tópico: %s): %s", topic, studentInfo)

	// Devolver la respuesta gRPC
	return &pb.StudentResponse{
		Success: true,
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Crear un escritor de Kafka
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaBootstrap},
		Balancer: &kafka.LeastBytes{}, // Puedes usar otros balanceadores si lo prefieres
	})

	log.Println("Kafka writer created successfully")

	defer writer.Close()

	// Inicializar el servidor gRPC con el writer de Kafka
	s := grpc.NewServer()
	pb.RegisterStudentServer(s, &server{writer: writer})

	log.Printf("Server swimming started on port %d", *port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

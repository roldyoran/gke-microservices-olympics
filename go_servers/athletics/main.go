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
	port           = flag.Int("port", 50052, "The server port")
	kafkaBootstrap = "my-cluster-kafka-bootstrap:9092" // Kafka cluster address
)

// Server is used to implement the gRPC server defined in the proto library
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

	// Determine if the student won or lost
	result := rand.Intn(2) // Generate a random number between 0 and 1

	var status string
	var topic string

	// Define the topic based on whether the student won or lost
	if result == 1 {
		status = "won"
		topic = "student-winners" // Topic for winners
	} else {
		status = "lost"
		topic = "student-losers" // Topic for losers
	}

	log.Printf("Student: %s, %s", in.GetName(), status)

	// Create the message to send to Kafka
	studentInfo := fmt.Sprintf(`{"name": "%s", "status": "%s", "faculty": "%s", "discipline": %d}`, in.GetName(), status, in.GetFaculty(), in.GetDiscipline())

	// Generate a random key
	randomKey := strconv.Itoa(rand.Intn(1000000)) // Generate a random number as the key

	// Send the message to Kafka
	err := s.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(randomKey), // Use the random key
			Value: []byte(studentInfo),
			Topic: topic,
		},
	)

	if err != nil {
		log.Printf("Error sending message to Kafka: %v", err)
		return nil, err
	}

	log.Printf("Message sent to Kafka (topic: %s): %s", topic, studentInfo)

	// Return the gRPC response
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

	// Create a Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaBootstrap},
		Balancer: &kafka.LeastBytes{}, // You can use other balancers if preferred
	})

	log.Println("Kafka writer created successfully")

	defer writer.Close()

	// Initialize the gRPC server with the Kafka writer
	s := grpc.NewServer()
	pb.RegisterStudentServer(s, &server{writer: writer})

	log.Printf("Athletics server started on port %d", *port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

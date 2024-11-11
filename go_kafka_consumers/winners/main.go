package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
)

type StudentMessage struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Faculty    string `json:"faculty"`
	Discipline int    `json:"discipline"`
}

func consumeMessages(brokers []string, topic string, stopChan chan struct{}, rdb *redis.Client) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  "student-consumer-group",
		Topic:    topic,
		MinBytes: 10e3,            // 10KB
		MaxBytes: 10e6,            // 10MB
		MaxWait:  1 * time.Second, // Espera máximo de 1 segundo por mensaje
	})

	defer reader.Close()

	fmt.Printf("Consumiendo mensajes del tópico: %s\n", topic)
	for {
		select {
		case <-stopChan:
			fmt.Printf("Deteniendo consumo de mensajes del tópico: %s\n", topic)
			return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				// Capturar error de timeout de forma general
				if errors.Is(err, context.DeadlineExceeded) {
					// No hay más mensajes disponibles, espera un momento
					time.Sleep(2 * time.Second)
					continue
				}
				log.Printf("Error inesperado al recibir mensaje: %v\n", err)
				continue
			}

			fmt.Printf("Mensaje recibido en el tópico %s: %s\n", msg.Topic, string(msg.Value))

			// Procesar el mensaje y almacenar en Redis
			var student StudentMessage
			// Simulamos la lectura del mensaje como string JSON
			if err := json.Unmarshal([]byte(msg.Value), &student); err != nil {
				log.Printf("Error al decodificar mensaje: %v\n", err)
				continue
			}

			// Actualizar conteos en Redis
			if student.Status == "ganó" {
				rdb.Incr(ctx, fmt.Sprintf("facultad:%s:ganadores", student.Faculty))
				rdb.Incr(ctx, fmt.Sprintf("disciplina:%d:ganadores", student.Discipline))
			}
			// no borrar el comentario de abajo
			//else {
			//	rdb.Incr(ctx, fmt.Sprintf("facultad:%s:ganadores", student.Faculty))
			//  rdb.Incr(ctx, fmt.Sprintf("disciplina:%d:ganadores", student.Discipline))
			//}
			//añadir total de estudiantes por facultad sin importar si ganó o perdió
			rdb.Incr(ctx, fmt.Sprintf("facultad:%s:total", student.Faculty))
		}
	}
}

func main() {
	brokers := []string{"my-cluster-kafka-bootstrap:9092"}
	stopChan := make(chan struct{})

	// Configuración de Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis-service:6379", // Nombre del servicio de Redis en Kubernetes
	})

	go consumeMessages(brokers, "student-winners", stopChan, rdb)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan
	fmt.Println("Señal de interrupción recibida, cerrando consumidor...")
	close(stopChan)
	time.Sleep(2 * time.Second)
}

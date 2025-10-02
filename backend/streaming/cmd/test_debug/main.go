package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"tchat.dev/streaming/models"
	"tchat.dev/streaming/repository"
)

func main() {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "tchat"
	cluster.Consistency = gocql.LocalQuorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 5 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer session.Close()

	repo, err := repository.NewChatMessageRepository(session)
	if err != nil {
		fmt.Printf("Failed to create repository: %v\n", err)
		return
	}

	ctx := context.Background()
	message := &models.ChatMessage{
		StreamID:          uuid.New(),
		MessageID:         uuid.New(),
		SenderID:          uuid.New(),
		Timestamp:         time.Now(),
		SenderDisplayName: "TestUser",
		MessageText:       "Test message",
		ModerationStatus:  models.ModerationStatusVisible,
		MessageType:       models.MessageTypeText,
	}

	err = repo.Create(ctx, message)
	if err != nil {
		fmt.Printf("Error creating message: %v\n", err)
		return
	}

	fmt.Println("Success!")
}

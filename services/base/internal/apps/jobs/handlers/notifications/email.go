package notifications

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func EmailJobHandler(rdb *redis.Client) {
	ctx := context.Background()
	for {
		// BRPop blocks until a message is available or timeout occurs
		// 0 timeout means block indefinitely
		result, err := rdb.BRPop(ctx, 0*time.Second, "emails").Result()
		if err != nil {
			time.Sleep(1 * time.Second) // wait before retrying
			continue
		}

		// result[0] is the key (queue name), result[1] is the value
		message := result[1]
		fmt.Printf("Received message: %s\n", message)

		//// Process the message
		//if err := processMessage(message); err != nil {
		//	log.Printf("Failed to process message %s: %v", message, err)
		//	// Optionally push to dead letter queue or retry
		//}
	}
}

package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/google/uuid"

	crmPb "github.com/usegro/proto/crm"
)

// InboundHandler processes inbound WhatsApp messages and status updates
type InboundHandler struct {
	dynamo    *dynamodb.Client
	crmClient crmPb.CRMServiceClient
}

func NewInboundHandler(dynamo *dynamodb.Client, crmClient crmPb.CRMServiceClient) *InboundHandler {
	return &InboundHandler{dynamo: dynamo, crmClient: crmClient}
}

// HandleInbound processes an inbound WhatsApp message:
// 1. Upsert contact in CRM
// 2. Find or create conversation
// 3. Store message in DynamoDB
// 4. Publish real-time event to frontend (via Redis)
func (h *InboundHandler) HandleInbound(ctx context.Context, event *InboundMessageEvent) error {
	log.Printf("[inbound] %s from %s (waba: %s)", event.Type, event.From, event.WabaID)

	// 1. Upsert contact
	contact, err := h.crmClient.UpsertContactByPhone(ctx, &crmPb.UpsertContactByPhoneRequest{
		TenantId: event.WabaID, // wabaID maps to tenantID
		Phone:    event.From,
		Name:     event.Contact.Name,
	})
	if err != nil {
		log.Printf("failed to upsert contact: %v", err)
		// Don't fail — still store the message
	}

	// 2. Find or create conversation in DynamoDB
	conversationID, err := h.findOrCreateConversation(ctx, event, contact)
	if err != nil {
		return err
	}

	// 3. Store message in DynamoDB
	msg := MessageRecord{
		PK:             fmt.Sprintf("CONV#%s", conversationID),
		SK:             fmt.Sprintf("MSG#%s", time.Now().Format(time.RFC3339Nano)),
		MessageID:      uuid.New().String(),
		WAMessageID:    event.MessageID,
		ConversationID: conversationID,
		TenantID:       event.WabaID,
		Direction:      "inbound",
		Type:           event.Type,
		Content:        string(event.Content),
		Status:         "received",
		FromPhone:      event.From,
		CreatedAt:      time.Now().Unix(),
	}

	item, err := attributevalue.MarshalMap(msg)
	if err != nil {
		return err
	}

	_, err = h.dynamo.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("usegro_messages"),
		Item:      item,
	})

	return err
}

func (h *InboundHandler) HandleStatusUpdate(ctx context.Context, event *StatusUpdateEvent) error {
	log.Printf("[status] message %s → %s", event.MessageID, event.Status)

	// Update message status in DynamoDB by WA message ID
	// Uses a GSI on wa_message_id for efficient lookup
	_, err := h.dynamo.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("usegro_messages"),
		Key: map[string]types.AttributeValue{
			"wa_message_id": &types.AttributeValueMemberS{Value: event.MessageID},
		},
		UpdateExpression: aws.String("SET #s = :status, updated_at = :ts"),
		ExpressionAttributeNames: map[string]string{
			"#s": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: event.Status},
			":ts":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().Unix())},
		},
	})

	return err
}

func (h *InboundHandler) findOrCreateConversation(ctx context.Context, event *InboundMessageEvent, contact *crmPb.Contact) (string, error) {
	// Try to find open conversation for this phone + tenant
	result, err := h.dynamo.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("usegro_conversations"),
		IndexName:              aws.String("phone-tenant-index"),
		KeyConditionExpression: aws.String("tenant_id = :tid AND contact_phone = :phone"),
		FilterExpression:       aws.String("#s = :open"),
		ExpressionAttributeNames: map[string]string{"#s": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":tid":   &types.AttributeValueMemberS{Value: event.WabaID},
			":phone": &types.AttributeValueMemberS{Value: event.From},
			":open":  &types.AttributeValueMemberS{Value: "open"},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return "", err
	}

	if len(result.Items) > 0 {
		var conv ConversationRecord
		attributevalue.UnmarshalMap(result.Items[0], &conv)
		return conv.ConversationID, nil
	}

	// Create new conversation
	convID := uuid.New().String()
	contactName := event.Contact.Name
	if contact != nil {
		contactName = contact.Name
	}

	conv := ConversationRecord{
		PK:             fmt.Sprintf("TENANT#%s", event.WabaID),
		SK:             fmt.Sprintf("CONV#%s", convID),
		ConversationID: convID,
		TenantID:       event.WabaID,
		ContactPhone:   event.From,
		ContactName:    contactName,
		Status:         "open",
		UnreadCount:    1,
		CreatedAt:      time.Now().Unix(),
		LastMessageAt:  time.Now().Unix(),
	}

	item, _ := attributevalue.MarshalMap(conv)
	h.dynamo.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("usegro_conversations"),
		Item:      item,
	})

	return convID, nil
}

// ─── DynamoDB record types ─────────────────────────────────

type MessageRecord struct {
	PK             string `dynamodbav:"PK"`
	SK             string `dynamodbav:"SK"`
	MessageID      string `dynamodbav:"message_id"`
	WAMessageID    string `dynamodbav:"wa_message_id"`
	ConversationID string `dynamodbav:"conversation_id"`
	TenantID       string `dynamodbav:"tenant_id"`
	Direction      string `dynamodbav:"direction"`
	Type           string `dynamodbav:"type"`
	Content        string `dynamodbav:"content"`
	Status         string `dynamodbav:"status"`
	FromPhone      string `dynamodbav:"from_phone"`
	AgentID        string `dynamodbav:"agent_id,omitempty"`
	CreatedAt      int64  `dynamodbav:"created_at"`
}

type ConversationRecord struct {
	PK             string `dynamodbav:"PK"`
	SK             string `dynamodbav:"SK"`
	ConversationID string `dynamodbav:"conversation_id"`
	TenantID       string `dynamodbav:"tenant_id"`
	ContactPhone   string `dynamodbav:"contact_phone"`
	ContactName    string `dynamodbav:"contact_name"`
	Status         string `dynamodbav:"status"`
	AssignedTo     string `dynamodbav:"assigned_to,omitempty"`
	UnreadCount    int    `dynamodbav:"unread_count"`
	Labels         []string `dynamodbav:"labels,omitempty"`
	CreatedAt      int64  `dynamodbav:"created_at"`
	LastMessageAt  int64  `dynamodbav:"last_message_at"`
}

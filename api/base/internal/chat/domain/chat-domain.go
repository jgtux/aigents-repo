package domain

import (
	"time"
)

type Chat struct {
	ChatUUID  string `json:"chat_uuid"`
	AgentUUID string `json:"agent_uuid"`
	AuthUUID  string `json:"auth_uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type Message struct {
	MessageUUID        string  `json:"message_uuid"`
	SenderUUID         string  `json:"sender_uuid"`
	SenderType         string `json:"sender_type"`
	ReceiverUUID       string  `json:"receiver_uuid"`
	ReceiverType       string `json:"receiver_type"`
	ChatUUID           string  `json:"chat_uuid"`
	MessageContent     MessageContent `json:"message_content"`
	CreatedAt          time.Time  `json:"created_at"`
}

type MessageContent struct {
	MessageContentUUID string  `json:"message_content_uuid"`
	Content            string     `json:"content"`
}

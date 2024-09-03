package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/8thgencore/microservice-chat/internal/model"
	chatv1 "github.com/8thgencore/microservice-chat/pkg/chat/v1"
)

// ToChatFromService converts service layer model to structure of API layer.
func ToChatFromService(chat *model.Chat) *chatv1.Chat {
	return &chatv1.Chat{
		Usernames: chat.Usernames,
	}
}

// ToChatFromDesc converts structure of API layer to service layer model.
func ToChatFromDesc(chat *chatv1.Chat) *model.Chat {
	return &model.Chat{
		Usernames: chat.Usernames,
	}
}

// ToMessageFromDesc converts structure of API layer to service layer model.
func ToMessageFromDesc(message *chatv1.Message) *model.Message {
	return &model.Message{
		From:      message.From,
		Text:      message.Text,
		Timestamp: message.Timestamp.AsTime(),
	}
}

// ToStreamFromDesc converts interface of API layer to service layer interface.
func ToStreamFromDesc(stream chatv1.ChatV1_ConnectServer) model.Stream {
	return stream.(model.Stream)
}

// ToMessageFromService converts service layer model to structure of API layer.
func ToMessageFromService(message *model.Message) *chatv1.Message {
	return &chatv1.Message{
		From:      message.From,
		Text:      message.Text,
		Timestamp: timestamppb.New(message.Timestamp),
	}
}

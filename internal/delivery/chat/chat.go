package chat

import (
	"context"
	"log"

	"github.com/8thgencore/microservice-chat/internal/converter"
	"github.com/golang/protobuf/ptypes/empty"

	chatv1 "github.com/8thgencore/microservice-chat/pkg/chat/v1"
)

// Connect is used for connecting to a chat.
func (i *ChatImplementation) Connect(req *chatv1.ConnectRequest, stream chatv1.ChatV1_ConnectServer) error {
	err := i.chatService.Connect(req.GetChatId(), req.GetUsername(), converter.ToStreamFromDesc(stream))
	log.Println(err)
	return err
}

// Create is used for creating new chat.
func (i *ChatImplementation) Create(ctx context.Context, req *chatv1.CreateRequest) (*chatv1.CreateResponse, error) {
	id, err := i.chatService.Create(ctx, converter.ToChatFromDesc(req.GetChat()))
	if err != nil {
		return nil, err
	}

	return &chatv1.CreateResponse{
		Id: id,
	}, nil
}

// Delete is used for deleting chat.
func (i *ChatImplementation) Delete(ctx context.Context, req *chatv1.DeleteRequest) (*empty.Empty, error) {
	err := i.chatService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// SendMessage is used for sending messages to connected chat.
func (i *ChatImplementation) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest) (*empty.Empty, error) {
	err := i.chatService.SendMessage(ctx, req.GetChatId(), converter.ToMessageFromDesc(req.GetMessage()))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

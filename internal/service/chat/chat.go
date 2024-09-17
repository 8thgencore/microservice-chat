package chat

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/8thgencore/microservice-chat/internal/converter"
	"github.com/8thgencore/microservice-chat/internal/model"
)

// Connect implements service.ChatService.
func (s *chatService) Connect(chatID string, username string, stream model.Stream) error {
	s.mxChannels.RLock()
	chatChan, ok := s.channels[chatID]
	s.mxChannels.RUnlock()

	if !ok {
		return errors.New("chat not found")
	}

	// Init streams for chat, if they don't exist
	s.mxChat.Lock()
	if _, okChat := s.chats[chatID]; !okChat {
		s.chats[chatID] = &chat{
			streams: make(map[string]model.Stream),
		}
	}
	s.mxChat.Unlock()

	// Set stream for user
	s.chats[chatID].m.Lock()
	s.chats[chatID].streams[username] = stream
	s.chats[chatID].m.Unlock()

	if err := s.loadHistory(chatID, stream); err != nil {
		// If history not loaded, there's no problem, you can still send messages
		log.Printf("failed to load history: %v", err)
	}

	for {
		select {
		case msg, okCh := <-chatChan:
			// Check if channel is closed
			if !okCh {
				return nil
			}

			// Send message for everyone in chat
			for _, st := range s.chats[chatID].streams {
				if err := st.Send(converter.ToMessageFromService(msg)); err != nil {
					return err
				}
			}
		case <-stream.Context().Done():
			// Delete stream for user when context is dead
			s.chats[chatID].m.Lock()
			delete(s.chats[chatID].streams, username)
			s.chats[chatID].m.Unlock()
			return nil
		}
	}
}

func (s *chatService) loadHistory(chatID string, stream model.Stream) error {
	messages, err := s.messagesRepository.GetMessages(stream.Context(), chatID)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		if err := stream.Send(converter.ToMessageFromService(msg)); err != nil {
			return err
		}
	}

	return nil
}

// Create implements service.ChatService.
func (s *chatService) Create(ctx context.Context, chat *model.Chat) (string, error) {
	var id string

	if s.txManager == nil {
		return "", errors.New("txManager is not initialized")
	}
	if s.chatRepository == nil {
		return "", errors.New("chatRepository is not initialized")
	}
	if s.logRepository == nil {
		return "", errors.New("logRepository is not initialized")
	}
	if s.channels == nil {
		s.channels = make(map[string]chan *model.Message) // initialize channels map
	}
	if s.messagesRepository == nil {
		return "", errors.New("messagesRepository is not initialized")
	}

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.Create(ctx, chat)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &model.Log{
			Text: fmt.Sprintf("Created chat with id: %v", id),
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		log.Print(err)
		return "", errors.New("failed to create chat")
	}

	// Create buffered channel for new chat
	s.channels[id] = make(chan *model.Message, messagesBuffer)

	return id, nil
}

// Delete implements service.ChatService.
func (s *chatService) Delete(ctx context.Context, id string) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.messagesRepository.DeleteChat(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.chatRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &model.Log{
			Text: fmt.Sprintf("Deleted chat with id: %v", id),
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		log.Print(err)
		return errors.New("failed to delete chat")
	}

	// Delete channel associated with chat
	delete(s.channels, id)

	return nil
}

// InitChannels implements service.ChatService.
func (s *chatService) InitChannels(ctx context.Context) error {
	// Get chats from repository
	ids, err := s.chatRepository.GetChats(ctx)
	if err != nil {
		return errors.New("failed to init existing chats")
	}

	// Fill chats and channels for already existing chats
	for _, id := range ids {
		s.channels[id] = make(chan *model.Message, messagesBuffer)
	}

	return nil
}

// SendMessage implements service.ChatService.
func (s *chatService) SendMessage(ctx context.Context, chatID string, message *model.Message) error {
	s.mxChannels.RLock()
	chatChan, ok := s.channels[chatID]
	s.mxChannels.RUnlock()

	if !ok {
		return errors.New("chat not found")
	}

	// Save message in repository
	err := s.messagesRepository.Create(ctx, chatID, message)
	if err != nil {
		return errors.New("failed to save message")
	}

	chatChan <- message

	return nil
}

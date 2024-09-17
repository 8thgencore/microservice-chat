package converter

import (
	"github.com/8thgencore/microservice-chat/internal/model"
	"github.com/8thgencore/microservice-chat/internal/repository/messages/dao"
)

// ToMessagesFromRepo converts repository layer model to structure of service layer.
func ToMessagesFromRepo(messages []*dao.Message) []*model.Message {
	var res []*model.Message
	for _, m := range messages {
		res = append(res, &model.Message{
			From:      m.From,
			Text:      m.Text,
			Timestamp: m.Timestamp,
		})
	}

	return res
}

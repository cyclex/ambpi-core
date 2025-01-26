package domain

import (
	"context"
	"time"

	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain/model"
)

type ChatUcase interface {
	IncomingMessages(payload api.CproMessage) (trxChatBotID string, err error)
	ReplyMessages(waID string, payload api.CproMessage) (isRedeem bool, outgoing string, err error)
	ChatToUser(waID string, chat []string, types, media string) (res []byte, statusCode int, err error)
	ChatToUserCoster(waID, chat, types, media string) (res []byte, statusCode int, err error)
	DoRedeem(usUcode model.UsersUniqueCode) (chat string, err error)
	CreateQueueReply(data api.PayloadReply) (err error)
	CreateConversationsLog(data model.ConversationsLog) (err error)
	RefreshToken(token string) (res api.RefreshResponse, statusCode int, err error)
	SetToken(ctx context.Context, token string, ttl time.Duration) (err error)
	GetToken(ctx context.Context) (token string, err error)
}

// type ChatUcaseV1 interface {
// 	ReplyMessages(waID, reply string, c model.ConversationsLog, isNew bool) (isRedeem bool, trxChatBotID string, err error)
// 	IncomingMessages(payload api.Message) (trxChatBotID string, err error)
// 	RefreshToken() (res []byte, statusCode int, err error)
// 	FindToken() (data model.Token, err error)
// 	SetToken(updated map[string]interface{}) (err error)
// 	ChatToUser(waID, chat, types, media string) (res []byte, statusCode int, err error)
// 	ValidateCode(usUcode model.UsersUniqueCode) (checkVoucher int)

// 	SetConversations(cond, updated map[string]interface{}) (affected int64, err error)
// }

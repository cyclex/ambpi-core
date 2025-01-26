package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain"
	"github.com/cyclex/ambpi-core/domain/model"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func Redeem(processing *bool, chatUcase domain.ChatUcase, orderUcase domain.OrdersUcase, ctx context.Context, debug bool) {

	*processing = true

	defer func() {
		*processing = false
	}()

	queue, err := orderUcase.GetQueueRedeem(ctx)
	if err != nil {
		return
	}

	for _, row := range queue {
		err = orderUcase.UpdateQueueRedeem(ctx, row.ID)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[Redeem]usecase - fail UpdateQueue. id:%v", row.ID))
			winLog.Error(err)
			return
		}

		uCode := model.UsersUniqueCode{
			WaID:          row.Messages.WaID,
			Name:          row.Messages.Name,
			Profession:    row.Messages.Profession,
			IsFormatValid: row.Messages.IsFormatValid,
			Msisdn:        row.Messages.WaID,
			Raw:           row.Messages.Raw,
			NIK:           row.Messages.NIK,
			County:        row.Messages.County,
		}

		chat, err := chatUcase.DoRedeem(uCode)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[Redeem]usecase - fail DoRedeem. data:%+v", uCode))
			winLog.Error(err)
			return
		}

		msg := api.PayloadReply{
			WaID:        row.Messages.WaID,
			Chat:        chat,
			ScheduledAt: time.Now().Local().Add(time.Second * 5).Unix(),
			Raw:         row.Messages.Raw,
		}
		err = chatUcase.CreateQueueReply(msg)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[Redeem]usecase - fail CreateQueueReply. data:%+v", msg))
			winLog.Error(err)
		}

	}

	*processing = false

	return
}

func DoSendReply(processing *bool, chatUcase domain.ChatUcase, orderUcase domain.OrdersUcase, ctx context.Context) {

	*processing = true
	var err error

	defer func() {
		*processing = false
		if err != nil {
			fmt.Println(err)
		}
	}()

	queue, err := orderUcase.GetQueueReply(ctx)
	if err != nil {
		return
	}

	for _, row := range queue {

		_, statusCode, err := chatUcase.ChatToUserCoster(row.Messages.WaID, row.Messages.Chat, "text", "")
		if err != nil {
			return
		}

		if statusCode != http.StatusOK {
			return
		}

		clog := model.ConversationsLog{
			SessionID: uuid.NewString(),
			Incoming:  row.Raw,
			WAID:      row.Messages.WaID,
			Outgouing: row.Messages.Chat,
		}
		err = chatUcase.CreateConversationsLog(clog)
		if err != nil {
			return
		}

		err = orderUcase.UpdateQueueReply(ctx, row.ID)
		if err != nil {
			return
		}

	}

	*processing = false

	return

}

func RefreshToken(processing *bool, chatUcase domain.ChatUcase) {

	*processing = true

	defer func() {
		*processing = false
	}()

	token, statusCode, err := chatUcase.RefreshToken("a29waWFiY21hbnRhcDpBQkNrb3BpNDhe")
	if err != nil {
		err = errors.Wrap(err, "[RefreshToken]usecase - fail RefreshToken.")
		return
	}

	if statusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("[RefreshToken]usecase - fail RefreshToken. statusCode:%d", statusCode))
		return
	}

	if token.Token != "" {
		expiresAt, err := time.Parse("2006-01-02 15:04:05", token.ExpiresAfter)
		if err != nil {
			err = errors.New(fmt.Sprintf("[RefreshToken]usecase - fail parsing expires_after:%v", token.ExpiresAfter))
			return
		}

		ttl := time.Until(expiresAt)
		err = chatUcase.SetToken(context.Background(), token.Token, ttl)
		if err != nil {
			err = errors.Wrap(err, "[RefreshToken]usecase - fail SetToken.")
			return
		}
	}

	*processing = false

	return
}

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
		err = errors.Wrap(err, "[cron.Redeem] GetQueueRedeem")
		winLog.Error(err)
		return
	}

	for _, row := range queue {
		err = orderUcase.UpdateQueueRedeem(ctx, row.ID)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[cron.Redeem] UpdateQueueRedeem. id:%v", row.ID))
			winLog.Error(err)
			return
		}

		files, err := chatUcase.DownloadMedia(ctx, row.Messages.MediaID)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[cron.Redeem] DownloadMedia. id:%v", row.ID))
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
			Receipt:       files,
		}

		chat, err := chatUcase.DoRedeem(uCode)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("[cron.Redeem] DoRedeem. data:%+v", uCode))
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
			err = errors.Wrap(err, fmt.Sprintf("[cron.Redeem] CreateQueueReply. data:%+v", msg))
			winLog.Error(err)
		}

	}

	*processing = false

}

func DoSendReply(processing *bool, chatUcase domain.ChatUcase, orderUcase domain.OrdersUcase, ctx context.Context) {

	*processing = true
	var err error

	defer func() {
		*processing = false
		if err != nil {
			err = errors.Wrap(err, "[cron.DoSendReply]")
			winLog.Error(err)
		}
	}()

	queue, err := orderUcase.GetQueueReply(ctx)
	if err != nil {
		err = errors.Wrap(err, "[cron.DoSendReply] GetQueueReply")
		winLog.Error(err)
		return
	}

	for _, row := range queue {

		_, statusCode, err := chatUcase.ChatToUserCoster(row.Messages.WaID, row.Messages.Chat, "text", "")
		if err != nil {
			err = errors.Wrap(err, "[cron.DoSendReply] ChatToUserCoster")
			winLog.Error(err)
			return
		}

		if statusCode != http.StatusOK {
			err = errors.Wrapf(err, "[cron.DoSendReply] ChatToUserCoster. statusCode: %d", statusCode)
			winLog.Error(err)
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
			err = errors.Wrap(err, "[cron.DoSendReply] CreateConversationsLog")
			winLog.Error(err)
			return
		}

		err = orderUcase.UpdateQueueReply(ctx, row.ID)
		if err != nil {
			err = errors.Wrap(err, "[cron.DoSendReply] UpdateQueueReply")
			winLog.Error(err)
			return
		}

	}

	*processing = false

}

func RefreshToken(processing *bool, chatUcase domain.ChatUcase) {

	*processing = true

	defer func() {
		*processing = false
	}()

	token, statusCode, err := chatUcase.RefreshToken()
	if err != nil {
		err = errors.Wrap(err, "[cron.RefreshToken] RefreshToken")
		winLog.Error(err)
		return
	}

	if statusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("statusCode:%d", statusCode))
		err = errors.Wrap(err, "[cron.RefreshToken] RefreshToken")
		winLog.Error(err)
		return
	}

	if token.Token != "" {
		expiresAt, err := time.Parse("2006-01-02 15:04:05", token.ExpiresAfter)
		if err != nil {
			err = errors.New(fmt.Sprintf("[cron.RefreshToken] fail parsing expires_after:%v", token.ExpiresAfter))
			winLog.Error(err)
			return
		}

		ttl := time.Until(expiresAt)
		err = chatUcase.SetToken(context.Background(), token.Token, ttl)
		if err != nil {
			err = errors.Wrap(err, "[cron.RefreshToken] SetToken")
			winLog.Error(err)
			return
		}
	}

	*processing = false

}

package main

import (
	"context"
	"time"

	"github.com/cyclex/ambpi-core/domain"
	"github.com/cyclex/ambpi-core/pkg"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var processing, processingAuth, processingRedeem, processingCalculateQuota, processingSendReply, processingRecurringQuota, processingRefresh, processingJob bool

var winLog *logrus.Logger

func InitCron(orderUcase domain.OrdersUcase, chatUcase domain.ChatUcase, cmsUcase domain.CmsUcase, timeout time.Duration, debug bool) {

	c := context.Background()
	_, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	winLog = pkg.New("winners", debug)

	cr := cron.New()

	cr.AddFunc("@every 5s", func() {
		if !processingRedeem {
			Redeem(&processingRedeem, chatUcase, orderUcase, c, debug)
		}
	})

	cr.AddFunc("@every 5s", func() {
		if !processingSendReply {
			DoSendReply(&processingSendReply, chatUcase, orderUcase, c)
		}
	})

	// cr.AddFunc("* * * * *", func() {
	// 	if !processingRefresh {
	// 		RefreshToken(&processingRefresh, chatUcase)
	// 	}
	// })

	cr.AddFunc("* * * * *", func() {
		if !processingJob {
			ProcessJobUpload(&processingJob, orderUcase, cmsUcase, c, debug)
		}
	})

	cr.Start()

}

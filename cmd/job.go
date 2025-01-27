package main

import (
	"context"

	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain"
	"github.com/pkg/errors"
)

func ProcessJobUpload(processing *bool, orderUcase domain.OrdersUcase, cmsUcase domain.CmsUcase, ctx context.Context, debug bool) {

	*processing = true

	defer func() {
		*processing = false
	}()

	var (
		status    bool
		totalRows int
		jobStatus = "2"
		files     string
	)

	queue, err := orderUcase.GetJob(ctx, "")
	if err != nil {
		err = errors.Wrap(err, "[cron.ProcessJobUpload] GetJob")
		winLog.Error(err)
		return
	}

	for _, row := range queue {
		if row.JobStatus == "1" {

			if row.JobType == "upload" {
				files = row.File
				status, totalRows, err = cmsUcase.ImportPrize(ctx, api.Job{File: row.File})
			} else {
				files, status, totalRows, err = cmsUcase.DownloadRedeem(ctx, api.Job{JobType: row.JobType, StartDate: row.StartAt, EndDate: row.EndAt})
			}

			if err != nil {
				err = errors.Wrap(err, "[cron.ProcessJobUpload]")
				winLog.Error(err)
			}

			if status {
				jobStatus = "3"
			}

			err = orderUcase.UpdateJob(ctx, api.Job{ID: row.ID, TotalRows: totalRows, JobStatus: jobStatus, File: files})
			if err != nil {
				err = errors.Wrap(err, "[cron.ProcessJobUpload] UpdateJob")
				winLog.Error(err)
			}
		}
	}

	*processing = false

}

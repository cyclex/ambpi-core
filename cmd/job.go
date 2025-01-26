package main

import (
	"context"
	"fmt"

	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain"
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
		fmt.Println(err.Error())
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
				fmt.Println(err.Error())
			}

			if status {
				jobStatus = "3"
			}

			err = orderUcase.UpdateJob(ctx, api.Job{ID: row.ID, TotalRows: totalRows, JobStatus: jobStatus, File: files})
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	*processing = false

}

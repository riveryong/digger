///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package client

import (
	"bytes"
	"digger/crawler"
	"digger/models"
	"errors"
	"fmt"
	"github.com/hetianyi/gox/logger"
)

type InMemLogWriter struct {
	data *bytes.Buffer
}

func (w *InMemLogWriter) Write(p []byte) (n int, err error) {
	w.data.Write(p)
	w.data.WriteString("\n")
	// fmt.Println(string(p))
	return len(p), nil
}

func (w *InMemLogWriter) Get() string {
	return w.data.String()
}

func processQueue(queue *models.Queue) (*models.QueueProcessResult, error) {
	// 生成的新queue
	var returnNewQueues []*models.Queue
	// 产生的爬虫结果
	var results []*models.Result
	var log = &InMemLogWriter{
		data: new(bytes.Buffer),
	}

	errChan := make(chan error)

	project := GetConfigSnapshot(queue.TaskId) //service.CacheService().GetSnapshotConfig(queue.TaskId)
	if project == nil {
		return nil, errors.New("cannot get project config snapshot")
	}

	Put(&job{
		finishChan: errChan,
		job: func() error {
			return crawler.Process(queue, project, log, func(oldQueue *models.Queue, newQueues []*models.Queue, _results []*models.Result, err error) {
				if err != nil {
					logger.Error(err)
					return
				}
				results = append(results, _results...)
				returnNewQueues = append(returnNewQueues, newQueues...)
			})
		},
	})
	err := <-errChan
	if err != nil {
		logger.Error(err)
		log.Write([]byte(fmt.Sprintf("<span style=\"color:#F38F8F\">Err: process queue(%d): %s</span>\n", queue.Id, err.Error())))
		return &models.QueueProcessResult{
			TaskId:    queue.TaskId,
			QueueId:   queue.Id,
			Expire:    queue.Expire,
			Logs:      log.Get(),
			NewQueues: nil,
			Results:   nil,
			Error:     log.Get(),
		}, err
	}
	return &models.QueueProcessResult{
		TaskId:    queue.TaskId,
		QueueId:   queue.Id,
		Expire:    queue.Expire,
		Logs:      log.Get(),
		NewQueues: returnNewQueues,
		Results:   results,
	}, nil
}

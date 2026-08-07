package services

import (
	"sync"

	"payment_verifier/config"
	"payment_verifier/models"
)

type BatchItemProcessor func(item string) (string, error)

func ProcessBatch(items []string, processor BatchItemProcessor, concurrency int) models.BatchVerifyResponse {
	if concurrency <= 0 {
		concurrency = config.Performance.DefaultConcurrency
	}

	total := len(items)
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	var mu sync.Mutex
	validList := make([]string, 0, total)
	failedList := make([]models.BatchVerifyFailedItem, 0, total)

	for _, item := range items {
		wg.Add(1)
		sem <- struct{}{}

		go func(it string) {
			defer func() {
				<-sem
				wg.Done()
			}()

			res, err := processor(it)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				failedList = append(failedList, models.BatchVerifyFailedItem{
					ReceiptID: it,
					Error:     err.Error(),
				})
			} else if res != "" {
				validList = append(validList, res)
			}
		}(item)
	}

	wg.Wait()

	return models.BatchVerifyResponse{
		Result: validList,
		Failed: failedList,
		Summary: models.BatchVerifySummary{
			Total:   total,
			Valid:   len(validList),
			Invalid: len(failedList),
		},
	}
}

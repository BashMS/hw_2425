package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// stopWorkers формирует сигнал для окончания работы воркеров.
func stopWorkers(sgnCh chan int, n int) {
	for i := 0; i < n; i++ {
		sgnCh <- i
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	currTask := -1
	countTask := len(tasks)
	errLimit := false
	errCount := 0
	doneCh := make(chan int, n)
	errCh := make(chan error, countTask)
	mut := sync.Mutex{}
	wg := sync.WaitGroup{}
	wgE := sync.WaitGroup{}

	// Воркеры выполнения задач
	for i := 0; i < n; i++ {
		// Создаем воркер обработки задач в отдельной рутине
		wg.Add(1)
		go func(currTask *int, eLimit bool) {
			defer func() {
				// Уменьшаем кол-во рутин при завершении работы воркера
				wg.Done()
			}()
			for {
				select {
				case <-doneCh:
					return
				default:
					// Получим задачу
					mut.Lock()
					*currTask++
					curr := *currTask
					mut.Unlock()
					if curr < countTask && !eLimit {
						err := tasks[curr]()
						errCh <- err
					} else {
						return
					}
				}
			}
		}(&currTask, errLimit)
	}

	// Запустим рутину подсчета ошибок
	wgE.Add(1)
	go func() {
		defer func() {
			stopWorkers(doneCh, n)
			wg.Wait()
			wgE.Done()
		}()
		resCount := 0
		for {
			err, ok := <-errCh
			if !ok {
				continue
			}
			resCount++
			if err != nil {
				errCount++
				errLimit = (m > 0) && (errCount == m)
			}
			if errLimit || resCount == countTask {
				return
			}
		}
	}()

	// Ожидаем окончания всех воркеров
	wgE.Wait()
	close(errCh)
	close(doneCh)

	if errLimit {
		return ErrErrorsLimitExceeded
	}

	return nil
}

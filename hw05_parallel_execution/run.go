package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	currTask := -1
	countTask := len(tasks)
	errLimit := false
	errCount := 0
	doneCh := make(chan struct{})
	errCh := make(chan error, m)
	mut := sync.Mutex{}
	wg := sync.WaitGroup{}

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
					if curr >= countTask || eLimit {
						return
					}
					err := tasks[curr]()
					if err != nil {
						errCh <- err
					}
				}
			}
		}(&currTask, errLimit)
	}

	// Запустим подсчет ошибок
	for {
		select {
		case err := <-errCh:
			if err != nil {
				errCount++
				errLimit = (m > 0) && (errCount == m)
			}
		default:
			break
		}
		mut.Lock()
		cTask := currTask
		mut.Unlock()
		if errLimit || cTask == countTask {
			close(doneCh)
			// Ожидаем окончания всех воркеров
			wg.Wait()
			break
		}
	}
	close(errCh)

	if errLimit {
		return ErrErrorsLimitExceeded
	}

	return nil
}

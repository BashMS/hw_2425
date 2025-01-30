package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	//ress := make([]Out, 0)
	result := make(Out)
	wg := sync.WaitGroup{}

	select {
	case <-done:
		return nil
	default:
		for i, stg := range stages {
			//result := make(Out)
			//result = stages[0](in)
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case <-done:
					return
				default:
					if i > 0 {
						result = stg(result)
					} else {
						result = stg(in)
					}
				}
			}()

		}
		//ress = append(ress, result)
	}

	wg.Wait()
	select {
	case <-done:
		return nil
	default:
		//results := make(chan interface{}, len(ress))
		//for _, item := range ress {
		//	results <- fmt.Sprintf("%v", <-item)
		//}
		//close(results)
		//return results
		return result
	}

	//return result
}

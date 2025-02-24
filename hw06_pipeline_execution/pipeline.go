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
	result := make(Bi)
	defer close(result)
	wg := sync.WaitGroup{}
	for val := range in {
		chIn := make(Bi,1)
		chIn <- val
		close(chIn)
		wg.Add(1)
		go func(ch In) {
			defer wg.Done()
			select {
			case <-done:
			   return
			default:
			   chOut := make(Out)
		       for i, stg := range stages {
			        select {
			        case <-done:
			            return 
		            default:
			            if i > 0 {
					      chOut = stg(chOut)
			            } else {
					      chOut = stg(ch)
			            }
			        }
			    }
				for val := range chOut {
					result <- val
				}
		    }
	    }(chIn)
	}

	wg.Wait()
	
	return result
}
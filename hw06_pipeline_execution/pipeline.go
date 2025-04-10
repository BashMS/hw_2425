package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var result Out

	for i, stage := range stages {
		if i == 0 {
			result = stage(in)
		} else {
			result = stage(result)
		}
	}

	return prepareResult(result, done)
}

func prepareResult(in In, done In) Out {
	result := make(Bi)

	go func() {
		for {
			select {
			case <-done:
				close(result)
				go func() {
					for i := range in {
						_ = i
					}
				}()
				return
			case v, ok := <-in:
				if !ok {
					close(result)
					return
				}
				result <- v
			}
		}
	}()

	return result
}

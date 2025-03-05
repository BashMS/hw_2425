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
			result = startStage(in, done, stage)
		} else {
			result = startStage(result, done, stage)
		}
	}

	return result
}

func startStage(in In, done In, stage Stage) Out {
	result := make(Bi)

	go func() {
		for {
			select {
			case <-done:
				go func() {
					for i := range in {
						_ = i
					}
				}()
				close(result)
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

	return stage(result)
}

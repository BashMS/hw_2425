package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	result := make(Out)
	go func() {
		<-done
		for v := range in {
			_ = v
		}
		for v := range result {
			_ = v
		}
	}()
	for i, stg := range stages {
		select {
		case <-done:
			return nil
		default:
			if i > 0 {
				result = stg(result)
			} else {
				result = stg(in)
			}
		}
	}

	return result
}

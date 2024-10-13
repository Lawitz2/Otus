package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

/*
TODO: look at TestAllStageStop test, gets stuck on that one
*/

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	channels := make([]Bi, 0, len(stages)+1)

	buffer := make(Bi)
	go func() {
		for v := range in {
			buffer <- v
		}
		close(buffer)
	}()
	channels = append(channels, buffer)

	for range len(stages) {
		ch := make(Bi)
		channels = append(channels, ch)
	}

	for i, stage := range stages {
		go func(stage Stage, i int) {
			ch := stage(channels[i])
			for {
				select {
				case v, ok := <-ch:
					if !ok {
						close(channels[i+1])
						return
					}
					channels[i+1] <- v
				case <-done:
					close(channels[i+1])
					return
				}
			}
		}(stage, i)
	}
	return channels[len(channels)-1]
}

package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	channels := make([]Bi, 0, len(stages)) // слайс каналов, которыми будут соединены stages
	for range len(stages) {
		channels = append(channels, make(Bi))
	}

	// обработка первого stage отдельно, т.к. канал in (type In)
	// не влезает в слайс каналов (type Bi)
	go func() {
		defer close(channels[0])
		for v := range stages[0](in) {
			channels[0] <- v
		}
	}()

	for i, stage := range stages[1:] {
		go func() {
			defer close(channels[i+1])
			for v := range stage(channels[i]) {
				channels[i+1] <- v
			}
		}()
	}

	output := make(Bi)

	go filter(channels[len(channels)-1], done, output)

	return output
}

/*
	filter - функция, которая смотрит за сигналом done.

Если его нет, и мы получаем значение после всех стадий
пайплайна - это значение отправляется в финальный output.
Как только появляется сигнал done - output закрывается, а остальные
значения выбрасываются для graceful shutdown горутин.
*/
func filter(in In, done In, output Bi) {
	defer close(output)
	for {
		select {
		case <-done:
			go func() {
				for range in {
					continue
				}
			}()
			return
		case v, ok := <-in:
			if !ok {
				return
			}
			output <- v
		}
	}
}

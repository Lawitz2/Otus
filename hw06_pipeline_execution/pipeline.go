package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

var wg = &sync.WaitGroup{}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	output := make(Bi)
	for _, stage := range stages {
		out := stage(in)
		in = out
	}

	go filter(in, done, output)

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

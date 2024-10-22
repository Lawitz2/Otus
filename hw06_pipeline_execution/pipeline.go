package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

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

Если мы получаем значение после всех стадий пайплайна - оно идёт в output.
Если мы получаем сигнал done - или входные данные заканчиваются -
мы сразу закрываем output, не задерживая поток который его читает.
Лишние значения будут выброшены, горутины завершатся.
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

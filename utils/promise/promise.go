package promise

type Promise[In any, Out any] struct {
	input  In
	result chan Out
	err    error
}

func NewPromise[In any, Out any](input In, handler func(In) (Out, error)) *Promise[In, Out] {
	promise := &Promise[In, Out]{
		input:  input,
		result: make(chan Out, 1),
	}

	go func(promise *Promise[In, Out]) {
		result, err := handler(input)
		promise.err = err
		promise.result <- result
	}(promise)

	return promise
}

func (promise *Promise[In, Out]) Get() (Out, error) {
	result := <-promise.result
	return result, promise.err
}

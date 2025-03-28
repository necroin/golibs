package generator

import "sync"

type Generator struct {
	value      int
	threadSafe bool
	mutex      *sync.RWMutex
}

func New(threadSafe bool) *Generator {
	return &Generator{
		value:      0,
		threadSafe: threadSafe,
		mutex:      &sync.RWMutex{},
	}
}

func (generator *Generator) Get() int {
	if generator.threadSafe {
		generator.mutex.RLock()
		defer generator.mutex.RUnlock()
	}
	return generator.value
}

func (generator *Generator) Next() int {
	if generator.threadSafe {
		generator.mutex.Lock()
		defer generator.mutex.Unlock()
	}
	result := generator.value
	generator.value += 1
	return result
}

package findex

type File[K comparable] interface {
	Close()
	Index() error
	Find(key K) ([]string, error)
}

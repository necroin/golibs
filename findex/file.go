package findex

type Row[T any] interface {
	UnmarshalCsv([]string) error
	*T
}

type File[K comparable, V any, R Row[V]] interface {
	Close()
	Index() error
	FindKey(key K) (*V, error)
	FindIndex(index int64) (*V, error)
}

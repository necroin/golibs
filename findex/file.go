package findex

type Row[T any] interface {
	UnmarshalCsv([]string) error
	*T
}

type File[K comparable, V any] interface {
	Close()
	Index() error
	RowCount() int64
	FindKey(key K) (*V, error)
	FindIndex(index int64) (*V, error)
}

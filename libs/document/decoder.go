package document

type Decoder interface {
	Decode(v any) (err error)
}

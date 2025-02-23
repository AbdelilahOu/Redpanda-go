package consumer

type HandlersInterface[T any] interface {
	HandleCreate(*T) error
	HandleUpdate(*T) error
	HandleDelete(*T) error
	HandleRead(*T) error
}

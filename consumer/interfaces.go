package consumer

type HandlersInterface[T any] interface {
	HandleCreate(*T) error
}

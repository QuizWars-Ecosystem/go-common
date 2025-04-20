package abstractions

type ConfigSubscriber[T any] interface {
	SectionKey() string
	UpdateConfig(newCfg T) error
}

package bark

type IBark interface {
	Bark() error
}

type Bark struct {
	DogName string
	Runner  func(string) error
}

func (b *Bark) Bark() error {
	err := b.Runner(b.DogName)
	if err != nil {
		return err
	}
	return err
}

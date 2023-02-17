package bark

import (
	"testing"
)

type fakeDog struct {
	Name string
}

func (f fakeDog) Bark() error {
	return nil
}

func Test_Bark(t *testing.T) {
	dog := fakeDog{}
	dog.Bark()
}

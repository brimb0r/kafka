package dalmatian

import (
	"fmt"
)

type Dalmatian struct{}

func (d *Dalmatian) DalmatianRunner() func(n string) error {
	return func(n string) error {
		fmt.Printf("%v - woof\n", n)
		return nil
	}
}

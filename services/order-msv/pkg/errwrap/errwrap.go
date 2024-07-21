package errwrap

import "fmt"

func Wrap(wrapper error, base error) error {
	return fmt.Errorf("%w: %w", wrapper, base)
}

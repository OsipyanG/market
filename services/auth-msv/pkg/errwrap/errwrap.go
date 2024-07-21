package errwrap

import "fmt"

func WrapIfErr(wrapper error, base error) error {
	if base != nil {
		return fmt.Errorf("%w: %w", wrapper, base)
	}

	return nil
}

func Wrap(wrapper error, base error) error {
	return fmt.Errorf("%w: %w", wrapper, base)
}

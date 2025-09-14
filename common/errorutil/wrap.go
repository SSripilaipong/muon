package errorutil

import "fmt"

func Wrapf(format string) func(error) error {
	return func(e error) error {
		return fmt.Errorf(format, e)
	}
}

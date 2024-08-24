package conversion

import "io"

type Converter interface {
	Convert(input io.Reader, output io.Writer) error
}

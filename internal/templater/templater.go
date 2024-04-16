package templater

import "io"

type Templater interface {
	GenerateTemplate(template string, data interface{}) (io.Reader, error)
}

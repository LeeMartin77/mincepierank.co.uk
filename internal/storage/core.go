package storage

import "context"

type Operations interface {
	Shutdown()
	TestData(c context.Context) (*string, error)
	Migrate(c context.Context) error
}

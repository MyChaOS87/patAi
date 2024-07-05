package authorization

import "github.com/MyChaOS87/patAi/pkg/middleware"

type Identity interface {
	GetID() string
}

type identity struct {
	id string
}

type provider struct{}

var _ middleware.AuthorizationProvider[Identity] = &provider{}

func (i *identity) GetID() string {
	return i.id
}

// Static mock as this is out of scope for this example.
func (a provider) GetByAPIKey(key string) (Identity, error) {
	if key == "user2" {
		return &identity{
			id: "mock-user2-id",
		}, nil
	}

	return &identity{
		id: "mock-default-id",
	}, nil
}

func NewMockProvider() middleware.AuthorizationProvider[Identity] {
	return &provider{}
}

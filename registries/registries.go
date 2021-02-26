package registries

import "github.com/rlaskowski/go-email"

type Registries struct {
	Email *email.Email
}

func NewRegistries() *Registries {
	return &Registries{
		Email: email.NewEmail(),
	}
}

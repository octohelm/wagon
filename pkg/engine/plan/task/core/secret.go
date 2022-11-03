package core

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/opencontainers/go-digest"

	"dagger.io/dagger"
)

func init() {
	DefaultFactory.Register(&Secret{})
}

type Secret struct {
	Meta struct {
		Secret struct {
			ID string `json:"id,omitempty"`
		} `json:"secret"`
	} `json:"$wagon"`
}

var secretids = sync.Map{}

func (s Secret) SecretID() dagger.SecretID {
	if id, ok := fsids.Load(s.Meta.Secret.ID); ok {
		return id.(dagger.SecretID)
	}
	return ""
}

func (s *Secret) SetSecretID(id dagger.SecretID) {
	key := digest.FromString(string(id)).String()
	fsids.Store(key, id)
	s.Meta.Secret.ID = key
}

func (v *Secret) SetSecretIDBy(ctx context.Context, secret *dagger.Secret) error {
	id, err := secret.ID(ctx)
	if err != nil {
		return err
	}
	v.SetSecretID(id)
	return nil
}

type SecretOrString struct {
	Value  string
	Secret *Secret
}

func (SecretOrString) OneOf() []any {
	return []any{
		"",
		&Secret{},
	}
}

func (s *SecretOrString) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '{' {
		se := &Secret{}
		if err := json.Unmarshal(data, se); err != nil {
			return err
		}
		s.Secret = se
		return nil
	}
	return json.Unmarshal(data, &s.Value)
}

func (s SecretOrString) MarshalJSON() ([]byte, error) {
	if s.Secret != nil {
		return json.Marshal(s.Secret)
	}
	return json.Marshal(s.Value)
}

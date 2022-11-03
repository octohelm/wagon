package core

import (
	"encoding/json"

	"golang.org/x/net/context"

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

func (s Secret) SecretID() dagger.SecretID {
	return dagger.SecretID(s.Meta.Secret.ID)
}

func (s *Secret) SetSecretID(id dagger.SecretID) {
	s.Meta.Secret.ID = string(id)
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

package core

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/opencontainers/go-digest"

	"dagger.io/dagger"
)

func init() {
	DefaultFactory.Register(&Socket{})
}

type Socket struct {
	Meta struct {
		Socket struct {
			ID string `json:"id,omitempty"`
		} `json:"socket"`
	} `json:"$wagon"`
}

var socketIDs = sync.Map{}

func (s Socket) SocketID() dagger.SocketID {
	if id, ok := fsids.Load(s.Meta.Socket.ID); ok {
		return id.(dagger.SocketID)
	}
	return ""
}

func (s *Socket) SetSocketID(id dagger.SocketID) {
	key := digest.FromString(string(id)).String()
	fsids.Store(key, id)
	s.Meta.Socket.ID = key
}

func (v *Socket) SetSocketIDBy(ctx context.Context, socket *dagger.Socket) error {
	id, err := socket.ID(ctx)
	if err != nil {
		return err
	}
	v.SetSocketID(id)
	return nil
}

type SocketOrString struct {
	Value  string
	Socket *Socket
}

func (SocketOrString) OneOf() []any {
	return []any{
		"",
		&Socket{},
	}
}

func (s *SocketOrString) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '{' {
		se := &Socket{}
		if err := json.Unmarshal(data, se); err != nil {
			return err
		}
		s.Socket = se
		return nil
	}
	return json.Unmarshal(data, &s.Value)
}

func (s SocketOrString) MarshalJSON() ([]byte, error) {
	if s.Socket != nil {
		return json.Marshal(s.Socket)
	}
	return json.Marshal(s.Value)
}

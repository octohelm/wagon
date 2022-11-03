package task

import (
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&Client{})
}

type Client struct {
	Env ClientEnv `json:"env"`

	Filesystem map[string]struct {
		Read  *ClientFilesystemRead  `json:"read,omitempty" cueExtra:"path: X"`
		Write *ClientFilesystemWrite `json:"write,omitempty" cueExtra:"path: X" wagon:"deprecated"`
	} `json:"filesystem"`

	Network map[string]ClientNetwork `json:"network" cueExtra:"address: X"`
}

type ClientFilesystemWrite struct {
	Path     string  `json:"path,omitempty"`
	Contents core.FS `json:"contents"`
}

type ClientNetwork struct {
	Path    string      `json:"address,omitempty"`
	Connect core.Socket `json:"connect"`
}

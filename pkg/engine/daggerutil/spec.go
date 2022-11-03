package daggerutil

import (
	"dagger.io/dagger"
	"github.com/dagger/dagger/core/schema"
)

type Directory struct {
	ID dagger.DirectoryID `json:"id,omitempty"`
}

type Container struct {
	ID dagger.ContainerID `json:"id,omitempty"`

	Platform string `json:"platform,omitempty"`

	Entrypoint  []string `json:"entrypoint,omitempty"`
	DefaultArgs []string `json:"defaultArgs,omitempty"`
	Workdir     string   `json:"workdir,omitempty"`

	User         string               `json:"user,omitempty"`
	EnvVariables []schema.EnvVariable `json:"envVariables,omitempty"`
	Labels       []schema.Label       `json:"labels,omitempty"`

	RootFS Directory `json:"rootfs,omitempty"`
}

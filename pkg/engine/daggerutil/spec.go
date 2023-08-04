package daggerutil

import (
	"dagger.io/dagger"
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

	User         string        `json:"user,omitempty"`
	EnvVariables []EnvVariable `json:"envVariables,omitempty"`
	Labels       []Label       `json:"labels,omitempty"`

	RootFS Directory `json:"rootfs,omitempty"`
}

type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EnvVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

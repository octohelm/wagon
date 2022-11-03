package docker

import (
	"dagger.io/dagger"
	"wagon.octohelm.tech/core"
)

#Ref: string

// A container image
#Image: core.#Image

// An empty container image (same as `FROM scratch` in a Dockerfile)
#Scratch: #Image & {
	rootfs: dagger.#Scratch
	config: {}
}

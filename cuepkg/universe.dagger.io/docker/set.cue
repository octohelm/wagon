package docker

import (
	"dagger.io/dagger/core"
)

#Set: {
	input: #Image

	config: core.#ImageConfig

	_set: core.#Set & {
		"input":  input.config
		"config": config
	}

	output: #Image & {
		rootfs:   input.rootfs
		platform: input.platform
		config:   _set.output
	}
}

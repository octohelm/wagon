package docker

import (
	"wagon.octohelm.tech/core"
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

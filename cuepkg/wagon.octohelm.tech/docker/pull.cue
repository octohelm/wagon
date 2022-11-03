package docker

import (
	"wagon.octohelm.tech/core"
)

#Pull: {
	source: #Ref

	resolveMode: *"default" | "forcePull" | "preferLocal"

	auth?: {
		username: string
		secret:   core.#Secret
	}

	platform?: string

	_pull: core.#Pull & {
		"source":      source
		"resolveMode": resolveMode
		if auth != _|_ {
			"auth": auth
		}
		if platform != _|_ {
			"platform": platform
		}
	}

	output: #Image & {
		rootfs:   _pull.output
		config:   _pull.config
		platform: _pull.platform
	}
}

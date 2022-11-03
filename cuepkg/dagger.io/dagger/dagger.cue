package dagger

import "wagon.octohelm.tech/core"

#Socket:  core.#Socket
#Secret:  core.#Secret
#FS:      core.#FS
#Scratch: core.#FS

#Plan: {
	client: core.#Client
	...
}

package dagger

import "wagon.octohelm.tech/core"

#Secret: core.#Secret

#FS: core.#FS

#Scratch: core.#FS

#Plan: {
	client: core.#Client
	...
}

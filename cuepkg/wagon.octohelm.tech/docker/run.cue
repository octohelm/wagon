package docker

import (
	"wagon.octohelm.tech/core"
)

#Run: {
	input: #Image

	mounts: [name=string]: core.#Mount
	env: [string]:         string | core.#Secret
	workdir?: string
	user?:    string

	entrypoint?: [...string]
	command?: {
		name: string
		args: [...string]
		flags: [string]: (string | true)
	}
	always?: bool

	_run: core.#Run & {
		"input":  input.rootfs
		"config": input.config
		"mounts": {
			for k, v in mounts {
				"\(k)": v
			}
		}
		"env": env
		if workdir != _|_ {
			"workdir": workdir
		}
		if user != _|_ {
			"user": user
		}
		if entrypoint != _|_ {
			"entrypoint": entrypoint
		}
		if command != _|_ {
			"command": command
		}
		if always != _|_ {
			"always": always
		}
	}

	exit: _run.exit

	output: #Image & {
		rootfs:   _run.output
		config:   input.config
		platform: input.platform
	}
}

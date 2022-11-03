package core

import "wagon.octohelm.tech/core"

#Source: core.#Source

#Mkdir:     core.#Mkdir
#ReadFile:  core.#ReadFile
#WriteFile: core.#WriteFile
#Copy:      core.#Copy
#Rm:        core.#Rm
#Merge:     core.#Merge
#Diff:      core.#Diff

#Subdir: {
	input: core.#FS
	path:  string

	_copy: #Copy & {
		contents: input
		source:   path
		dest:     "/"
	}

	output: _copy.output
}

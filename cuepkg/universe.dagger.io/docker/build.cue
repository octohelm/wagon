package docker

import (
	"wagon.octohelm.tech/core"
)

#Pull:       core.#ImagePull
#Run:        core.#ImageRun
#Copy:       core.#ImageCopy
#Set:        core.#ImageSet
#Dockerfile: core.#ImageDockerfile

// Modular build API for Docker containers
#Build: {
	steps: [#Step, ...#Step]
	output: #Image

	_dag: {
		for idx, step in steps if idx == 0 {
			"\(idx)": step
		}

		for idx, step in steps if idx > 0 {
			"\(idx)": {
				_prev: _dag["\(idx-1)"].output

				step & {
					input: _prev
				}
			}
		}
	}

	if len(_dag) > 0 {
		output: _dag["\(len(_dag)-1)"].output
	}
}

// A build step is anything that produces a docker image
#Step: {
	input?: #Image
	output: #Image
	...
}

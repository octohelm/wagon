package docker

import (
	"path"

	"wagon.octohelm.tech/core"
)

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

#Step: {
	input?: #Image
	output: #Image
	...
}

#Copy: {
	input:    #Image
	contents: core.#FS
	source:   string | *"/"
	dest:     string | *"."
	include: [...string]
	exclude: [...string]

	_copy: core.#Copy & {
		"input":    input.rootfs
		"contents": contents
		"source":   source
		"include":  include
		"exclude":  exclude
		"dest":     [
				if path.IsAbs(dest, path.Unix) {
				dest
			},
			if input.config.workdir == _|_ {
				path.Join(["/", dest], path.Unix)
			},
			if input.config.workdir != _|_ {
				path.Join([input.config.workdir, dest], path.Unix)
			},
			dest,
		][0]
	}

	output: #Image & {
		"platform": input.platform
		"config":   input.config
		"rootfs":   _copy.output
	}
}

#Dockerfile: {
	source: core.#FS

	dockerfile: {
		path:      string | *"Dockerfile"
		contents?: string
	}

	auth: [registry=string]: {
		username: string
		secret:   core.#Secret
	}

	platform?: string
	target?:   string
	buildArg: [string]: string
	label: [X=string]:  string

	_build: core.#Dockerfile & {
		"source":     source
		"auth":       auth
		"dockerfile": dockerfile
		if platform != _|_ {
			"platform": platform
		}
		if target != _|_ {
			"target": target
		}
		"buildArg": buildArg
		"label":    label
	}

	output: #Image & {
		platform: _build.platform
		rootfs:   _build.output
		config:   _build.config
	}
}

package docker

import (
	"wagon.octohelm.tech/core"
)

#Ref:     string
#Image:   core.#Image
#Scratch: #Image & {
	rootfs: core.#FS
	config: {}
}

package docker

import (
	"wagon.octohelm.tech/core"
)

#Push: {
	dest:   #Ref
	result: #Ref & _push.result

	auth?: {
		username: string
		secret:   core.#Secret
	}

	_push: core.#Push & {
		"dest": dest
		if auth != _|_ {
			"auth": auth
		}
	}

	image?: core.#Image
	images: [Platform=string]: core.#Image

	if image != _|_ {
		_push: {
			input:    image.rootfs
			config:   image.config
			platform: image.platform
		}
	}

	if image == _|_ {
		_push: {
			inputs: {
				for _p, _image in images {
					"\(_p)": {
						input:  _image.rootfs
						config: _image.config
					}
				}
			}
		}
	}
}

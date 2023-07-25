import (
	"strings"
	"wagon.octohelm.tech/core"

	"github.com/innoai-tech/runtime/cuepkg/tool"
	"github.com/innoai-tech/runtime/cuepkg/golang"
	"github.com/innoai-tech/runtime/cuepkg/debian"
	"github.com/innoai-tech/runtime/cuepkg/imagetool"
)

pkg: version: core.#Version

client: core.#Client & {
	env: {
		GH_USERNAME: string | *""
		GH_PASSWORD: core.#Secret
	}
}

setting: core.#Setting & {
	registry: "ghcr.io": auth: {
		username: client.env.GH_USERNAME
		secret:   client.env.GH_PASSWORD
	}
}

actions: go: golang.#Project & {
	version: "\(pkg.version.output)"

	source: {
		path: "."
		include: [
			"cmd/",
			"pkg/",
			"cuepkg/",
			"internal/",
			"go.mod",
			"go.sum",
		]
	}
	goos: [
		"linux",
		"darwin",
	]
	goarch: [
		"amd64",
		"arm64",
	]
	main: "./cmd/wagon"
	ldflags: [
		"-s -w",
		"-X \(go.module)/pkg/version.version=\(go.version)",
	]

	build: {
		pre: [
			"go mod download",
		]
	}

	ship: {
		name: "\(strings.Replace(go.module, "github.com/", "ghcr.io/", -1))/\(go.binary)"

		from: "docker.io/library/debian:bookworm-slim"

		steps: [
			debian.#InstallPackage & {
				packages: {
					"git":  _
					"wget": _
					"curl": _
					"make": _
				}
			},
			imagetool.#Shell & {
				run: """
					ln -s /wagon /bin/wagon
					ln -s /wagon /bin/dagger
					"""
			},
		]
		config: {
			entrypoint: ["/bin/sh"]
		}
	}
}

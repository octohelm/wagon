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
		GOPROXY:   string | *""
		GOPRIVATE: string | *""
		GOSUMDB:   string | *""

		LINUX_MIRROR: string | *""
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

	env: {
		GOPROXY:   client.env.GOPROXY
		GOPRIVATE: client.env.GOPRIVATE
		GOSUMDB:   client.env.GOSUMDB
	}

	build: pre: [
		"go mod download",
	]

	ship: {
		name: "\(strings.Replace(go.module, "github.com/", "ghcr.io/", -1))/\(go.binary)"

		from: "docker.io/library/debian:bullseye-slim"

		steps: [
			imagetool.#Shell & {
				env: {
					LINUX_MIRROR: client.env.LINUX_MIRROR
				}
				run: """
						if [ "${LINUX_MIRROR}" != "" ]; then
							sed -i "s@http://deb.debian.org@${LINUX_MIRROR}@g" /etc/apt/sources.list
							sed -i "s@http://security.debian.org@${LINUX_MIRROR}@g" /etc/apt/sources.list
						fi
					"""
			},
			debian.#InstallPackage & {
				packages: {
					"git":  _
					"wget": _
					"curl": _
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

	// FIXME remove when all migrated
	mirror: linux: client.env.LINUX_MIRROR
}

pkg: {
	name:    "hello"
	version: "0.1.0"
}

src: #Source & {
	path: "."
}

dest: {
	"etc": #Copy & {
		contents: src.output
		source:   "etc"
		include: ["*"]
		exclude: ["*.log"]
	}
}

//
//vendor: {
//	containerd: {
//		version: "1.6.9"
//
//		#Fetch & {
//			source: "https://github.com/containerd/containerd/releases/download/v\(version)/cri-containerd-cni-\(version)-\(context.target.os)-\(context.target.arch).tar.gz"
//		}
//	}
//}
//
//context: _
//
//#GoBuild: #Exec & {
//	source: _
//	dest:   _
//	env: {
//		GOOS:        context.target.os
//		GOARCH:      context.target.arch
//		CGO_ENABLED: 0
//	}
//	script: "go build -v -o \(output.root)\(dest) \(source)"
//}

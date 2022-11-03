module: "github.com/octohelm/wagon"

require: {
	"github.com/innoai-tech/runtime": "v0.0.0-20230207052751-0171bdb4eff0"
	"wagon.octohelm.tech":            "v0.0.0-20200202235959-0uyuz2sq+r5o"
}

require: {
	"dagger.io":          "v0.3.0" @indirect()
	"universe.dagger.io": "v0.3.0" @indirect()
}

module: "github.com/octohelm/wagon"

require: {
	"github.com/innoai-tech/runtime": "v0.0.0-20230208111820-a3dd976b6379"
	"wagon.octohelm.tech":            "v0.0.0-20200202235959-e46af4e91b7c"
}

require: {
	"dagger.io":          "v0.3.0" @indirect()
	"universe.dagger.io": "v0.3.0" @indirect()
}

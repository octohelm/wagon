package version

var (
	version  = "0.0.0.dev"
	revision = ""
)

func shortSha(sha string) string {
	if len(sha) > 8 {
		return ".shortSha+" + sha[0:8]
	}
	return ".shortSha+00000000"
}

func Version() string {
	return version + shortSha(revision)
}

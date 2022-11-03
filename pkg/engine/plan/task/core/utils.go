package core

import (
	"dagger.io/dagger"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func DefaultPlatform(platform string) dagger.Platform {
	if platform == "" {
		arch := os.Getenv("BUILDKIT_ARCH")
		if arch == "" {
			arch = runtime.GOARCH
		}
		return dagger.Platform(fmt.Sprintf("linux/%s", arch))
	}
	return dagger.Platform(platform)
}

type StringOrBool struct {
	String string
	Bool   *bool
}

func (StringOrBool) OneOf() []any {
	return []any{
		"",
		true,
	}
}

func (s *StringOrBool) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] != '"' {
		b := false
		if err := json.Unmarshal(data, &b); err != nil {
			return err
		}
		s.Bool = &b
		return nil
	}
	return json.Unmarshal(data, &s.String)
}

func (s StringOrBool) MarshalJSON() ([]byte, error) {
	if s.Bool != nil {
		return []byte(strconv.FormatBool(*s.Bool)), nil
	}
	return []byte(strconv.Quote(s.String)), nil
}

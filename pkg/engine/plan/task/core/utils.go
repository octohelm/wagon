package core

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"os"
	"strconv"

	"dagger.io/dagger"
)

type ImagePullPrefixier interface {
	ImagePullPrefix(name string) string
}

type imagePullPrefixierContext struct {
}

func ContextWithImagePullPrefixier(ctx context.Context, p ImagePullPrefixier) context.Context {
	return context.WithValue(ctx, imagePullPrefixierContext{}, p)
}

func ImagePullPrefixierFromContext(ctx context.Context) ImagePullPrefixier {
	if f, ok := ctx.Value(imagePullPrefixierContext{}).(ImagePullPrefixier); ok {
		return f
	}
	return &imagePullPrefixierDiscord{}
}

type imagePullPrefixierDiscord struct {
}

func (imagePullPrefixierDiscord) ImagePullPrefix(name string) string {
	return name
}

func DefaultPlatform(platform string) dagger.Platform {
	if platform == "" {
		arch := os.Getenv("BUILDKIT_ARCH")
		if arch != "" {
			return dagger.Platform(fmt.Sprintf("linux/%s", arch))
		}
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

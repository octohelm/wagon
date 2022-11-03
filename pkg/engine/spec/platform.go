package spec

import (
	"fmt"

	"github.com/containerd/containerd/platforms"
)

func ParsePlatform(p string) (*Platform, error) {
	platform, err := platforms.Parse(p)
	if err != nil {
		return nil, err
	}

	return &Platform{
		OS:         platform.OS,
		Arch:       platform.Architecture,
		Variant:    platform.Variant,
		OSVersion:  platform.OSVersion,
		OSFeatures: platform.OSFeatures,
	}, nil
}

type Platform struct {
	Arch       string   `json:"arch"`
	OS         string   `json:"os"`
	Variant    string   `json:"variant,omitempty"`
	OSVersion  string   `json:"os_version,omitempty"`
	OSFeatures []string `json:"os_features,omitempty"`
}

func (p *Platform) StorageKey() string {
	base := fmt.Sprintf("%s-%s", p.OS, p.Arch)
	if p.Variant != "" {
		return fmt.Sprintf("%s-%s", base, p.Variant)
	}
	return base
}

func (p *Platform) String() string {
	base := fmt.Sprintf("%s/%s", p.OS, p.Arch)

	if p.Variant != "" {
		return fmt.Sprintf("%s/%s", base, p.Variant)
	}

	return base
}

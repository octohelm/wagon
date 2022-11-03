package core

import (
	"encoding/json"
	"reflect"

	"dagger.io/dagger"
)

func init() {
	DefaultFactory.Register(&Mount{})
}

type Mounter interface {
	MountType() string
	MountTo(client *dagger.Client, container *dagger.Container) *dagger.Container
}

type Mount struct {
	Mounter `json:"-"`
}

func (m *Mount) UnmarshalJSON(data []byte) error {
	mt := &struct {
		Type string `json:"type"`
	}{}

	if err := json.Unmarshal(data, mt); err != nil {
		return err
	}

	for _, v := range m.OneOf() {
		if i, ok := v.(Mounter); ok {
			if i.MountType() == mt.Type {
				i = reflect.New(reflect.TypeOf(i).Elem()).Interface().(Mounter)

				if err := json.Unmarshal(data, i); err != nil {
					return err
				}

				m.Mounter = i
				return nil
			}
		}

	}

	return nil
}

func (m *Mount) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Mounter)
}

func (Mount) OneOf() []any {
	return []any{
		&MountCacheDir{},
		&MountTemp{},
		&MountFs{},
		&MountSecret{},
		&MountFile{},
	}
}

var _ Mounter = &MountCacheDir{}
var _ Mounter = &MountTemp{}
var _ Mounter = &MountFs{}
var _ Mounter = &MountSecret{}
var _ Mounter = &MountFile{}

type MountCacheDir struct {
	Type     string   `json:"type" enum:"cache"`
	Dest     string   `json:"dest"`
	Contents CacheDir `json:"contents"`
}

type CacheDir struct {
	ID          string `json:"id"`
	Concurrency string `json:"concurrency" default:"shared" wagon:"deprecated"`
}

func (MountCacheDir) MountType() string {
	return "cache"
}

func (m MountCacheDir) MountTo(client *dagger.Client, c *dagger.Container) *dagger.Container {
	return c.WithMountedCache(m.Dest, client.CacheVolume(m.Contents.ID), dagger.ContainerWithMountedCacheOpts{})
}

type MountTemp struct {
	Type     string  `json:"type" enum:"tmp"`
	Dest     string  `json:"dest"`
	Contents TempDir `json:"contents"`
}

type TempDir struct {
	Size int64 `json:"size" default:"0" wagon:"deprecated"`
}

func (t MountTemp) MountTo(client *dagger.Client, container *dagger.Container) *dagger.Container {
	return container.WithMountedTemp(t.Dest)
}

func (MountTemp) MountType() string {
	return "tmp"
}

type MountFs struct {
	Type     string  `json:"type" enum:"fs"`
	Dest     string  `json:"dest"`
	Contents FS      `json:"contents"`
	Source   *string `json:"source,omitempty"`

	ReadOnly bool `json:"ro,omitempty" wagon:"deprecated"`
}

func (f MountFs) MountTo(client *dagger.Client, container *dagger.Container) *dagger.Container {
	dir := client.Directory(dagger.DirectoryOpts{
		ID: f.Contents.DirectoryID(),
	})

	if source := f.Source; source != nil {
		dir = dir.Directory(*source)
	}

	return container.WithMountedDirectory(f.Dest, dir)
}

func (MountFs) MountType() string {
	return "fs"
}

type MountSecret struct {
	Type     string `json:"type" enum:"secret"`
	Dest     string `json:"dest"`
	Contents Secret `json:"contents"`

	Uid  int `json:"uid" default:"0" wagon:"deprecated"`
	Gid  int `json:"gid" default:"0" wagon:"deprecated"`
	Mask int `json:"mask" default:"0o644" wagon:"deprecated"`
}

func (m MountSecret) MountTo(client *dagger.Client, container *dagger.Container) *dagger.Container {
	return container.WithMountedSecret(m.Dest, client.Secret(m.Contents.SecretID()))
}

func (MountSecret) MountType() string {
	return "secret"
}

type MountFile struct {
	Type        string `json:"type" enum:"file"`
	Dest        string `json:"dest"`
	Contents    string `json:"contents"`
	Permissions int    `json:"mask" default:"0o644"`
}

func (m MountFile) MountTo(client *dagger.Client, container *dagger.Container) *dagger.Container {
	f := client.Container().
		WithNewFile("/tmp", dagger.ContainerWithNewFileOpts{
			Contents:    m.Contents,
			Permissions: m.Permissions,
		}).
		File("/tmp")

	return container.WithMountedFile(m.Dest, f)
}

func (MountFile) MountType() string {
	return "file"
}

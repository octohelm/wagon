package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-courier/logr"
	"golang.org/x/net/context"
)

type Info struct {
	File string
}

type Downloader interface {
	Resolve(ctx context.Context, remote *url.URL) (*Info, error)
}

func NewDownloader(root string) Downloader {
	return &downloader{
		root: root,
	}
}

type downloader struct {
	root string
}

func (d *downloader) Resolve(ctx context.Context, remote *url.URL) (*Info, error) {
	info := Info{}
	info.File = filepath.Join(d.root, "download", remote.Host, remote.Path)

	if _, err := os.Stat(info.File); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err := d.downloadTo(ctx, remote, info.File); err != nil {
			return nil, err
		}
	}

	return &info, nil
}

var downloadingTasks = sync.Map{}

func (d *downloader) downloadTo(ctx context.Context, remote *url.URL, file string) error {
	if ch, ok := downloadingTasks.Load(file); ok {
		return <-ch.(chan error)
	}

	doneCh := make(chan error)
	downloadingTasks.Store(file, doneCh)

	go func() {
		defer func() {
			downloadingTasks.Delete(file)
			close(doneCh)
		}()

		doneCh <- downloadTo(ctx, remote.String(), file)
	}()

	return <-doneCh
}

func downloadTo(ctx context.Context, source string, file string) error {
	l := logr.FromContext(ctx)

	dir := filepath.Dir(file)

	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	f, err := os.CreateTemp(dir, "")
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	c := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, "GET", source, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	total := resp.ContentLength
	written := 0

	if total > 0 {
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()
		go func() {
			for range t.C {
				l.Info(fmt.Sprintf("downloaded %s / %s", FileSize(written), FileSize(total)))
			}
		}()
	}

	if _, err = io.Copy(io.MultiWriter(f, WriteFunc(func(b []byte) (int, error) {
		written += len(b)
		return len(b), nil
	})), resp.Body); err != nil {
		return err
	}

	l.Info(fmt.Sprintf("downloaded."))

	if err := os.Rename(f.Name(), file); err != nil {
		return err
	}

	return nil
}

func WriteFunc(fn func(b []byte) (int, error)) io.Writer {
	return &writeFunc{fn: fn}
}

type writeFunc struct {
	fn func(b []byte) (int, error)
}

func (w *writeFunc) Write(p []byte) (n int, err error) {
	return w.fn(p)
}

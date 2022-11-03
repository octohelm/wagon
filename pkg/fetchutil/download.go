package fetchutil

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/fsutil"
	"golang.org/x/net/context"
)

func NewClient() *http.Client {
	return &http.Client{}
}

func DownloadTo(ctx context.Context, fs fsutil.Fs, source url.URL, dest string) error {
	l := logr.FromContext(ctx)

	uri := source.String()

	if dest == "" {
		dest = filepath.Base(uri)
	}

	dir := filepath.Dir(dest)
	if dir != "" {
		_ = fs.MkdirAll(dir, 0755)
	}

	f, err := fs.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return err
	}

	resp, err := NewClient().Do(req)
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

package strfmt

import "net/url"

type URL struct {
	url.URL
}

func (x *URL) UnmarshalText(text []byte) error {
	u, err := url.Parse(string(text))
	if err != nil {
		return err
	}
	x.URL = *u
	return nil
}

func (x URL) MarshalText() (text []byte, err error) {
	return []byte(x.String()), nil
}

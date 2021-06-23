package email

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime"
	"strings"

	"github.com/paulrosania/go-charset/charset"
)

func decode(encoded string) (string, error) {
	wd := new(mime.WordDecoder)

	wd.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if !strings.Contains(charset, "utf-8") {
			c, err := decodeCharset(charset, input)
			if err != nil {
				return nil, err
			}

			return bytes.NewReader(c), nil
		}

		return input, nil
	}

	return wd.DecodeHeader(encoded)
}

func decodeCharset(ch string, r io.Reader) ([]byte, error) {
	r, err := charset.NewReader(ch, r)
	if err != nil {
		return nil, err
	}

	c, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return c, nil
}

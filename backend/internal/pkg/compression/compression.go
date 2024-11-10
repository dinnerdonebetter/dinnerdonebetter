package compression

import (
	"bytes"
	"fmt"
	"io"

	"github.com/klauspost/compress/zstd"
)

const (
	algoZstd   algo = "zstd"
	algoS2     algo = "s2"     // TODO: implement
	algoSnappy algo = "snappy" // TODO: implement
)

type (
	algo string

	Compressor interface {
		CompressBytes(in []byte) ([]byte, error)
		DecompressBytes(in []byte) ([]byte, error)
	}
)

type compressor struct {
	algo algo
}

func (c *compressor) CompressBytes(in []byte) ([]byte, error) {
	switch c.algo {
	case algoZstd:
		var b bytes.Buffer
		enc, err := zstd.NewWriter(&b)
		if err != nil {
			return nil, err
		}

		if _, err = io.Copy(enc, bytes.NewReader(in)); err != nil {
			return nil, err
		}

		if err = enc.Close(); err != nil {
			return nil, err
		}

		return b.Bytes(), nil
	default:
		return nil, fmt.Errorf("unsupported compression algorithm: %s", c.algo)
	}
}

func (c *compressor) DecompressBytes(in []byte) ([]byte, error) {
	switch c.algo {
	case algoZstd:
		d, err := zstd.NewReader(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}
		defer d.Close()

		var b bytes.Buffer
		if _, err = io.Copy(&b, d); err != nil {
			return nil, err
		}

		return b.Bytes(), nil
	default:
		return nil, fmt.Errorf("unsupported compression algorithm: %s", c.algo)
	}
}

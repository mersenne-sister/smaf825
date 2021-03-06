package chunk

import (
	"io"

	"fmt"
	"strings"

	"bytes"

	"github.com/but80/smaf825/smaf/subtypes"
	"github.com/but80/smaf825/smaf/util"
	"github.com/pkg/errors"
)

type MMMGEXVOChunk struct {
	*ChunkHeader
	Stream    []uint8             `json:"stream"`
	Exclusive *subtypes.Exclusive `json:"exclusive"`
}

func (c *MMMGEXVOChunk) Traverse(fn func(Chunk)) {
	fn(c)
}

func (c *MMMGEXVOChunk) String() string {
	result := "MMMGEXVOChunk: " + c.ChunkHeader.String()
	sub := []string{}
	if c.Exclusive != nil {
		sub = append(sub, fmt.Sprintf("Exclusive: %s", c.Exclusive.String()))
	} else {
		sub = append(sub, fmt.Sprintf("Stream: %s", util.Hex(c.Stream)))
	}
	return result + "\n" + util.Indent(strings.Join(sub, "\n"), "\t")
}

func (c *MMMGEXVOChunk) Read(rdr io.Reader) error {
	c.Stream = make([]uint8, c.ChunkHeader.Size)
	_, err := rdr.Read(c.Stream)
	if err != nil {
		return errors.WithStack(err)
	}
	if !(c.Stream[0] == 0xFF && c.Stream[1] == 0xF0) {
		return nil
	}
	c.Exclusive = subtypes.NewExclusive(false)
	rest := len(c.Stream) - 2
	err = c.Exclusive.Read(bytes.NewReader(c.Stream[2:]), &rest)
	if err != nil {
		return errors.WithStack(err)
	}
	if rest != 0 {
		return fmt.Errorf("Wrong size of EXVO exclusive data")
	}
	return nil
}

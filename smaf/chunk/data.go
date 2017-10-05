package chunk

import (
	"io"

	"fmt"
	"strings"

	"github.com/mersenne-sister/smaf825/smaf/util"
	"github.com/pkg/errors"
)

type DataChunk struct {
	*ChunkHeader
	Stream     []uint8
	HasOptions bool
	Options    struct {
		Vendor          string
		Carrior         string
		Category        string
		Title           string
		Artist          string
		LyricWriter     string
		Composer        string
		Arranger        string
		Copyright       string
		ManagementGroup string
		ManagementInfo  string
		CreatedDate     string
		UpdatedDate     string
		EditStatus      string
		VCard           string
	}
}

func (c *DataChunk) Traverse(fn func(Chunk)) {
	fn(c)
}

func (c *DataChunk) CodeType() int {
	return int(c.ChunkHeader.Signature & 255)
}

func (c *DataChunk) String() string {
	result := "DataChunk: " + c.ChunkHeader.String()
	sub := []string{
		fmt.Sprintf("Code type: 0x%02X", c.CodeType()),
		fmt.Sprintf("Stream: %s", util.Escape(c.Stream)),
		fmt.Sprintf("Options: %+v", c.Options),
	}
	return result + "\n" + util.Indent(strings.Join(sub, "\n"), "\t")
}

func (c *DataChunk) Read(rdr io.Reader) error {
	c.Stream = make([]uint8, c.ChunkHeader.Size)
	n, err := rdr.Read(c.Stream)
	if err != nil {
		return err
	}
	if n < len(c.Stream) {
		return errors.Errorf("Cannot read enough byte length specified in chunk header")
	}
	if c.CodeType() == 0x00 {
		options := map[string]string{}
		i := 0
		for i < len(c.Stream) {
			tag := string(c.Stream[i : i+2])
			i += 2
			size := int(c.Stream[i])<<8 | int(c.Stream[i+1])
			i += 2
			value := util.DecodeShiftJIS(c.Stream[i : i+size])
			i += size
			options[tag] = value
		}
		c.HasOptions = true
		c.Options.Vendor = options["VN"]
		c.Options.Carrior = options["CN"]
		c.Options.Category = options["CA"]
		c.Options.Title = options["ST"]
		c.Options.Artist = options["AN"]
		c.Options.LyricWriter = options["WW"]
		c.Options.Composer = options["SW"]
		c.Options.Arranger = options["AW"]
		c.Options.Copyright = options["CR"]
		c.Options.ManagementGroup = options["GR"]
		c.Options.ManagementInfo = options["MI"]
		c.Options.CreatedDate = options["CD"]
		c.Options.UpdatedDate = options["UD"]
		c.Options.EditStatus = options["ES"]
		c.Options.VCard = options["VC"]
	}
	return nil
}

package clause

import (
	"strconv"
	"strings"
)

type Config struct {
	start      int
	hit        int
	format     string
	rerankSize int
}

const (
	formatJSON     = "json"
	formatFullJSON = "fulljson"
)

const (
	strDefault = 0
	hitDefault = 10
	fmtDefault = "json" // "fulljson"
	rszDefault = 200
)

// NewConfig new config
func NewConfig() *Config {
	return &Config{
		start:      strDefault,
		hit:        hitDefault,
		format:     fmtDefault,
		rerankSize: rszDefault,
	}
}

func (c *Config) SetStart(val int) *Config      { c.start = val; return c }
func (c *Config) SetHit(val int) *Config        { c.hit = val; return c }
func (c *Config) SetRerankSize(val int) *Config { c.rerankSize = val; return c }
func (c *Config) SetFormatJSON() *Config        { c.format = formatJSON; return c }
func (c *Config) SetFormatFullJSON() *Config    { c.format = formatFullJSON; return c }

func (c *Config) Source() string {
	var s []string

	if c.start >= 5000 {
		s = append(s, "start:5000")
	} else if c.start > 0 {
		s = append(s, "start:"+strconv.FormatInt(int64(c.start), 10))
	}

	if c.hit > 500 {
		s = append(s, "hit:500")
	} else if c.hit > 0 {
		s = append(s, "hit:"+strconv.FormatInt(int64(c.hit), 10))
	}

	s = append(s, "format:"+c.format)

	if c.hit > 2000 {
		s = append(s, "rerank_size:2000")
	} else if c.rerankSize > 0 {
		s = append(s, "rerank_size:"+strconv.FormatInt(int64(c.rerankSize), 10))
	}

	return "config=" + strings.Join(s, ",")
}

package pagination

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const ctEncodedCursorItems = 3

type Cursor struct {
	Time       time.Time
	ID         string
	Parameters map[string][]string
}

func (c Cursor) String() string {
	return fmt.Sprintf("%s,%s,%s", c.Time.Format(time.RFC3339Nano), c.ID, c.Parameters)
}

func (c Cursor) encodeParameters() (string, error) {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(c.Parameters); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (c *Cursor) decodeParameters(s string) error {
	r := strings.NewReader(s)
	var m map[string][]string
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return err
	}
	c.Parameters = m
	return nil
}

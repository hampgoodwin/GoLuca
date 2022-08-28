package pagination

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/hampgoodwin/errors"
	"github.com/segmentio/ksuid"
)

func ParseCursor(encodedCursor string) (Cursor, error) {
	c := Cursor{}

	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		err = errors.WithErrorMessage(err, errors.NotValidRequest, "cursor is not valid base64")
		return c, err
	}

	arrStr := strings.SplitN(string(byt), ",", ctEncodedCursorItems)
	if len(arrStr) != ctEncodedCursorItems {
		err = errors.New("cursor is invalid")
		return c, err
	}

	t, err := time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return c, err
	}
	c.Time = t

	id := arrStr[1]
	_, err = ksuid.Parse(id)
	if err != nil {
		return c, err
	}
	c.ID = id

	if err = c.decodeParameters(arrStr[2]); err != nil {
		return c, err
	}

	return c, nil
}

func (c Cursor) EncodeCursor() (string, error) {
	parameters, err := c.encodeParameters()
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("%s,%s,%s", c.Time.Format(time.RFC3339Nano), c.ID, parameters)
	return base64.StdEncoding.EncodeToString([]byte(key)), nil
}

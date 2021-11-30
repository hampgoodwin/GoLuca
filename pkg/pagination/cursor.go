package pagination

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/hampgoodwin/errors"

	"github.com/google/uuid"
)

func DecodeCursor(encodedCursor string) (id string, res time.Time, err error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		err = errors.WithErrorMessage(err, errors.NotValidRequest, "cursor is not valid base64")
		return
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		err = errors.New("cursor is invalid")
		return
	}

	res, err = time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return
	}
	id = arrStr[1]
	_, err = uuid.Parse(id)
	if err != nil {
		return
	}
	return
}

func EncodeCursor(t time.Time, id string) string {
	key := fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), id)
	return base64.StdEncoding.EncodeToString([]byte(key))
}

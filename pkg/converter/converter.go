package converter

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"strings"
	"time"
)

// AnyToBytesBuffer Convert bytes to buffer helper
func AnyToBytesBuffer(i interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(i)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func UnixToDate(unix int64) string {
	return time.UnixMilli(unix).UTC().String()
}

// AnyToBytesStringWithJoin Convert bytes to buffer helper
func AnyToBytesStringWithJoin(i []int64) string {
	return strings.Trim(strings.Replace(fmt.Sprint(i), " ", ",", -1), "[]")
}

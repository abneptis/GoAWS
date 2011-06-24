package aws

import (
	"http"
  "strconv"
	"time"
)
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

var SQSTimestampFormat = "2006-01-02T15:04:05MST"
var ISO8601TimestampFormat = "2006-01-02T15:04:05Z"

// Adds a current timestamp to the request (deleting any present timestamp).
//  - you may optionally pass in a custom function to create the timestamp value,
// otherwise we default to time.UTC().

func Timestamp(req *http.Request, t func() *time.Time) {
	if t == nil {
		req.Form.Set("Timestamp", time.UTC().Format(ISO8601TimestampFormat))
	} else {
		req.Form.Set("Timestamp", t().Format(ISO8601TimestampFormat))
	}
}

// Adds an expiration message to the request (as opposed to a timestamp)
// You may supply a custom time() function, else time.Seconds() is used.
func Expires(req *http.Request, t func() *time.Time, from_now int64) {
	if t == nil {
		req.Form.Set("Expires", strconv.Itoa64(time.Seconds() + from_now))
	} else {
		req.Form.Set("Expires", strconv.Itoa64(t().Seconds()+from_now))
	}
}

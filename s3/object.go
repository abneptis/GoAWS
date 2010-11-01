package s3

import "io"

type Object struct {
  Key string
  Body io.ReadCloser
}

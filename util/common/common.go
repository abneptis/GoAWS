package common

import (
  "flag"
  "os"
)

// Common functionality to ease sub-modules

type UserCall struct {
  Args []flag.Flag
  F    func()(os.Error)
}

type Calls map[string]UserCall

package common

import (
  "flag"
  "os"
)

// Common functionality to ease sub-modules

type UserCall struct {
  Name string
  Args []flag.Flag
  f    func()(os.Error)
}

type Calls map[string]UserCall

package elog

import (
   "os"
)

// RedirectStderr to the file passed in
func RedirectStderr() (err error) {
   logFile, err := os.OpenFile("./test-error.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
   if err != nil {
      return
   }
   // Keep it cross-platform; on Windows syscall.Dup2 is unavailable.
   os.Stderr = logFile
   return
}

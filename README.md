# logfile

## Golang logrotate files helper

### Import

`go get github.com/milQA/logfile`

### Example

```go
package main

import (
    "log"
    "time"

    "github.com/milQA/logfile"
)

func main() {
    logFileName:= logfile
    rotateTickDuration:= time.Minute

    logSave := logfile.NewLogSave(logFileName, rotateTickDuration)
    defer logSave.Close()
    log.SetOutput(logSave)

    log.Println("Start service")
    for i := 0; i < 10; i++ {
        log.Printf("Log #%d", i)
        time.Sleep(10 * time.Second)
    }
    log.Println("Stop service")
}
```
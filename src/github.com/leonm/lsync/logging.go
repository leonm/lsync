package main

import (
    "io"
    "os"
    "io/ioutil"
    "log"
)

var (
    Trace   *log.Logger
    Info    *log.Logger
    Error   *log.Logger
)

func InitLogging(
    traceHandle io.Writer,
    infoHandle io.Writer,
    errorHandle io.Writer) {

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func InitStdOutLogging() {
  InitLogging(ioutil.Discard, os.Stdout, os.Stdout)
}

func InitFileLogging(name string) {
  f, err := os.Create(name)
  check(err)
  InitLogging(ioutil.Discard, f, f)
}

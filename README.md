# Middleware

Middleware provides a convenient way to chain
generic middleware functions.

In short, it transforms

    Middleware1(Middleware2(Middleware3(App)))

to

    middleware.New(Middleware1, Middleware2, Middleware3)

### Why?

There are several middleware solutions out there, but all that I've found deal
primarily with HTTP.  I wanted something where I could use generic non-HTTP
specific middleware components.  This project is influenced by
https://github.com/mitchellh/middleware and https://github.com/justinas/alice

### Usage

Your middleware handlers should have the form of

    func(env map[string]interface{}, next MiddlewareHandler)


```go
package main

import (
    "github.com/jcarley/middleware"
)

func getFileSize(env map[string]interface{}, next MiddlewareHandler) {
  fileSize := env["fileSize"];
  // ... do some work
  env["newResult"] = "Result"
  next(
}

func doSomethingElse(env map[string]interface{}) {
  // ... do some other work
}

func doWork(env map[string]interface{}) {
  result := env["newResult"];
  // ... log it somewhere
}

func main() {

    env := make(map[string]interface{})
    env["initialState"] = "somevalue"

    getFileSizeHandler := MiddlewareHandleFunc(getFileSize)
    doSomethingElseHandler := MiddlewareHandleFunc(doSomethingElse)

    chain := middleware.New(getFileSizeHandler, doSomethingElseHandler).ThenFunc(doWork)
    chain.Call(env)
}
```

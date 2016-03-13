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
https://github.com/mitchellh/middleware and https://github.com/codegangsta/negroni

### Usage

Install the package with

    go get github.com/jcarley/middleware

Then import

    import "github.com/jcarley/middleware"

Your middleware handlers should have the form of

    func(env map[string]interface{}, next middleware.HandlerFunc)


Example

```go
package main

import (
    "github.com/jcarley/middleware"
)

func getFileSize(env map[string]interface{}, next middleware.HandlerFunc) {
  // ... do some work
  env["fileSize"] = 1000
  next(env)
}

func doSomethingElse(env map[string]interface{}, next middleware.HandlerFunc) {
  // ... do some other work
  next(env)
}

func doWork(env map[string]interface{}, next middleware.HandlerFunc) {
  result := env["newResult"];
  // ... log it somewhere
  next(env)
}

func main() {

    env := make(map[string]interface{})
    env["initialState"] = "somevalue"

    chain := New()
    chain.UseFunc(getFileSizeHandler)
    chain.UseFunc(doSomethingElseHandler)
    chain.UseFunc(doWork)
    chain.Call(env)
}
```

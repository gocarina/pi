pi
==

Simple toolkit for API in Go.

It offers a simple but powerful request router with interceptors to execute code before, after or when an error occurs while handling the request.

Complete Example
================

```go

package main

import (
    "github.com/gocarina/pi"
)


func GetEveryUsers(c *pi.RequestContext) error {
    users, err := db.GetEveryUsers()
    if err != nil {
        return pi.NewError(500, err)
    }
    return c.WriteJSON(users)
}

func GetSingleUser(c *pi.RequestContext) error {
    user, err := db.GetUserByID(c.GetRouteVariable("id"))
    if err != nil {
        return pi.NewError(400, err)
    }
    return c.WriteJSON(user)
}

func MainHandler(c *pi.RequestContext) error {
    return c.WriteString("Hello, 世界")
}

func AddUser(c *pi.RequestContext) error {
    user := &User{}
    if err := c.GetJSONObject(user); err != nil {
        return pi.NewError(400, err)
    }
     user, err := db.AddUser(user)
     if err != nil {
        return pi.NewError(500, err)
    }
    return c.WriteJSON(user)
}

func main() {
    p := pi.New()
    p.Router("/",
        p.Route("/users",
            p.Route("/{id}").Get(GetSingleUser)
        ).Get(GetEveryUsers).Post(AddUser)
    ).Any(MainHandler)
}



```

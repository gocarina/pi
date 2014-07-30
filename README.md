pi
==

Simple toolkit for API in Go: http://godoc.org/github.com/gocarina/pi.

It offers a simple but powerful request router with interceptors to execute code before, after or when an error occurs while handling the request.

Complete Example
================

```go

package main

import (
    "github.com/gocarina/pi"
)


// GetEveryUsers matches /users (GET)
func GetEveryUsers(c *pi.RequestContext) error {
    users, err := db.GetEveryUsers()
    if err != nil {
        return pi.NewError(500, err)
    }
    return c.WriteJSON(users)
}

// GetSingleUser matches /users/{id} (GET)
func GetSingleUser(c *pi.RequestContext) error {
    user, err := db.GetUserByID(c.GetRouteVariable("id"))
    if err != nil {
        return pi.NewError(400, err)
    }
    return c.WriteJSON(user)
}

// MainHandler matches / (ANY)
func MainHandler(c *pi.RequestContext) error {
    return c.WriteString("Hello, 世界")
}

// AddUser matches /users (POST)
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

// DeleteUser matches /users/{id} and /users/{id}/delete (DELETE)
func DeleteUser(c *pi.RequestContext) error {
    if err := db.DeleteUser(c.GetRouteVariable("id")); err != nil {
        return http.NewError(404, err)
    }
    return c.WriteJSON(pi.J{"status": "OK"})
}

// AuthorizeUser will be called for each routes starting by "/users"
func AuthorizeUser(c *pi.RequestContext) error {
    // Authorize User...
    return nil
}

func main() {
    p := pi.New()
    p.Router("/",
        p.Route("/users",
            p.Route("/{id}",
                p.Route("/delete").Delete(DeleteUser)).
            Get(GetSingleUser).Delete(DeleteUser)).
        Get(GetEveryUsers).Post(AddUser).BeforeFunc(AuthorizeUser)).
    Any(MainHandler)
    
    p.ListenAndServe(":8080")
}


```


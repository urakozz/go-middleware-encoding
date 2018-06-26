# go-middleware-encoding
Content encoding middleware for the go-json-rest

### Encoders:

 - None Encoder 
 - Deflate Encoder (RFC 1950 and RFC 1951)
 - Gzip Encoder (RFC 1952)

### Optional C encoders:

 - C Gzip Encoder (RFC 1952)
 - C Brotli encoder (no RFC yet)

 
### Examples


#### 1. Using default settings (deflate, gzip)
```go
package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/urakozz/go-middleware-encoding"
    "log"
    "net/http"
)

func main() {

    mw := []rest.Middleware{
    		&rest.PoweredByMiddleware{
    			XPoweredBy: "Golang",
    		},
    		&mwencoding.EncodingMiddleware{},
    	}
    	
    api := rest.NewApi()
    api.Use(mw...)
    api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
        w.WriteJson(map[string]string{"Body": "Hello World!"})
    }))
    log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
```

#### 2. Using (deflate, C-gzip, brotli)
```go
package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/urakozz/go-middleware-encoding"
    "github.com/urakozz/go-middleware-encoding/encoding/cencoding"
    "log"
    "net/http"
)

func main() {
    mwencoding.RegisterEncoder(&cencoding.EncoderCGzip{})
    mwencoding.RegisterEncoder(&cencoding.EncoderCBrotli{})
    
    mw := []rest.Middleware{
    		&rest.PoweredByMiddleware{
    			XPoweredBy: "Golang",
    		},
    		&mwencoding.EncodingMiddleware{},
    	}
    	
    api := rest.NewApi()
    api.Use(mw...)
    api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
        w.WriteJson(map[string]string{"Body": "Hello World!"})
    }))
    log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
```

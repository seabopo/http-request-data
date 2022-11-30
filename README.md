http-request-data
=========================
HTTP Request Data is a small go web server that displays information about the client HTTP request.

The default port is 80, but this is configurable via the `-listen` flag:

```
http-request-data -listen=:8080
```

Then visit http://localhost:8080/ in your browser.

This is modified version of [HashiCorp's HTTP-ECHO](https://github.com/hashicorp/http-echo) repository on Github.

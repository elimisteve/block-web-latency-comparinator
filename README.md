# WebPipes Block: Web Latency Comparinator

[WebPipes](http://www.webpipes.org/) block written in
[Go](http://golang.org) which accepts an array of URLs, performs HEAD
requests to all URLs (in parallel), then returns the time it took (in
milliseconds) to receive each response


## Web Latency Comparinator

    curl -i -X POST \
    -d '{"inputs": {"urls": ["http://google.com", "http://yahoo.com"]}}' \
    http://web-latency-comparinator.herokuapp.com/

To run this example locally, clone the repo and start up the service:

```
git clone https://github.com/elimisteve/block-web-latency-comparinator
cd block-web-latency-comparinator
go run web.go
```

In another terminal, run this command:

    curl -i -X POST \
    -d '{"inputs": {"urls": ["http://google.com/", "http://yahoo.com/"]}}' \
    http://localhost:8080/

You should receive a response similar to the following:

```
HTTP/1.1 200 OK
Content-Length: 98
Date: Sun, 08 Sep 2013 04:18:21 GMT

{"outputs":[{"url":"http://google.com/","latency":147},{"url":"http://yahoo.com/","latency":577}]}
```


## Block Definition

```javascript
{
  "name": "Web Latency Comparinator",
  "url": "http://web-latency-comparinator.herokuapp.com",
  "description": "Performs a HEAD request to all given URLs in parallel and returns the time taken to receive a response from each.",
  "inputs": {
    "name": "urls",
    "type": "Array",
    "description": "URLs to be visited"
  },
  "outputs": [
    {
      "name": "url",
      "type": "String",
      "description": "Visited URL"
    },
    {
      "name": "latency",
      "type": "Number",
      "description": "URL response time (in milliseconds)"
    }
  ]
}
```

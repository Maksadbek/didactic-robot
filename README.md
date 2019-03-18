Interview assignment by Okmeter

# Go HTTPS Client
Utility to send metrics over HTTP(S)

## What it can do
* You can set the timeout for resolve.
* First go through the local resolver, if it fails, fallbacks to 8.8.8.8.
* If it resolves to multiple IP, then select random (if an error is made, make a second attempt through another IP).
* You can set a connection timeout.
* You can set the timeout for TLS handshake.
* Enable TCP KEEPALIVE (idle, interval, fail_after), if the OS does not know TCP KA, set the I/O timeout.

## Usage
```Go
config := Config{
        ResolveTimeout:      time.Second * time.Duration(3),
        TLSHandshakeTimeout: time.Second * 3,
        NameserverAddr:      "8.8.8.8:53",
        ConnectionTimeout:   time.Second * 3,
        TCPKeepIdleTime:     time.Second * 3,
        TCPKeepIntervalTime: time.Second * 3,
        TCPKeepFailAfter:    3,
        ReadTimeout:         time.Second * 60,
        WriteTimeout:        time.Second * 60,
}

client, err := NewClient(config)
if err != nil {
        panic(err)
}
request, err := http.NewRequest(http.MethodGet, "http://httpbin.org/get", nil)
if err != nil {
        panic(err)
}
resp, err := client.Do(request)
if err != nil {
        panic(err)
}

defer resp.Body.Close()

resp, err = ioutil.ReadAll(resp.Body)
if err != nil {
        panic(err)
}

// do whatever with resp...

```

## CLI
Very simple cli that sends file contents to given endpoint
```
maxs-MacBook-Pro:drobot max$ ./drobot -h
Usage of ./drobot:
  -connection-timeout duration
    	set connection timeout
  -endpoint string
    	endpoint where data will be sent (default "https://httpbin.org/post")
  -filename string
    	file which content will be sent with POST request (default "inputs.json")
  -nameserver-addr string
    	set up default fallback address (default "8.8.8.8")
  -read-timeout duration
    	set read timeout if OS do not support TCP Keep-alive, in seconds (default 1m0s)
  -resolve-timeout duration
    	set resolve timeout
  -tcp-keep-failafter duration
    	set TCP keepalive fail after value (default 3s)
  -tcp-keep-idle duration
    	set TCP keepalive idle in seconds (default 1s)
  -tcp-keep-interval duration
    	set TCP keepalive interval, in seconds (default 1s)
  -tls-handshake-timeout duration
    	set TLS handshake timeout
  -write-timeout duration
    	set write timeout if OS do not support TCP Keep-alive, in seconds (default 1m0s)

maxs-MacBook-Pro:drobot max$ cat inputs.json 
{
	"foo": "bar",
	"bar": "baz"
}

maxs-MacBook-Pro:drobot max$ ./drobot 
DROBOT
{
  "args": {}, 
  "data": "{\n\t\"foo\": \"bar\",\n\t\"bar\": \"baz\"\n}\n", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Connection": "close", 
    "Content-Length": "33", 
    "Host": "httpbin.org", 
    "User-Agent": "Go-http-client/1.1"
  }, 
  "json": {
    "bar": "baz", 
    "foo": "bar"
  }, 
  "origin": "46.39.248.63", 
  "url": "https://httpbin.org/post"
}
```

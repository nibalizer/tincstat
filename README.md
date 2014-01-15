# Uptime Server

[![Build Status](https://travis-ci.org/mkinney/uptime-server.png)](https://travis-ci.org/mkinney/uptime-server)

## API Usage

### Uptime

Get the system uptime.

Resource: /uptime  
Method: GET

#### Curl Example
```
curl -i http://127.0.0.1:9000/uptime
```

Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 17 Dec 2013 00:57:36 GMT
Content-Length: 111

{
 "one_minute": 1.0199999809265137,
 "five_minutes": 1.2100000381469727,
 "fifteen_minutes": 1.2300000190734863
}
```

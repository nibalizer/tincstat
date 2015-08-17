# Tinc Status Server

## API Usage

### Uptime

Get some basics of the tinc daemon, must be run as the user of tinc

Resource: /tincstat
Method: GET

#### Curl Example
```
curl -i http://127.0.0.1:9000/tincstat
```

Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 17 Dec 2013 00:57:36 GMT
Content-Length: 111

{
  "total_bytes_in": 115324,
  "total_bytes_out": 67990,
  "connections": [
    {
      "name": "some_random_node",
      "ip": "192.0.2.15",
      "port": 2003
    }
  ]
}

```

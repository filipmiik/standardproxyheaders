# Standard Proxy Headers Traefik Plugin

This plugin constructs and adds standard HTTP headers to outgoing requests according to RFC 7239 (Forwarded) and RFC 9110 (Via).

## Configuration

```yaml
http:
  middlewares:
    standard-proxy-headers:
      plugins:
        standard-proxy-headers:
          forwardedByHostname: true
          forwardedByHeader: ""
          forwardedByValue: "traefik"
          forwardedForRemote: true
          forwardedForHeader: "X-Real-IP"
          forwardedForValue: "client"
```

## Headers

### Forwarded

This header is added to all requests in format according to RFC 7239.
If the request already contains this header, the new value is appended to it with `, ` (comma and space) as a separator.

> If any of the fields would evaluate to an empty string, they will not be included/set in the final header since the
> RFC defines all of them as optional.

#### Forwarded by field

Forwarded `by` field configuration, `hostname > header > value` if more are specified at the same time.

- If `forwardedByHostname` is `true`, the field is set to `os.Hostname()`.
- If `forwardedByHeader` is not empty, the field is set to the value of such request header.
- If `forwardedByValue` is not empty, the field is set to the provided value.

#### Forwarded for field

Forwarded `for` field configuration, `remote > header > value` if more are specified at the same time.

- If `forwardedForRemote` is `true`, the field is set to the hostname of the remote address.
- If `forwardedForHeader` is not empty, the field is set to the value of such request header.
- If `forwardedForValue` is not empty, the field is set to the provided value.

> When using the ProxyProtocol module, the `forwardedForRemote` can be used to get the client IP address since the remote
> address is resolved to the one from the proxy header.

#### Forwarded host field

Set from the request `Host` header.

#### Forwarded proto field

Set to `https` if SSL/TLS connection was used, otherwise `http`.

### Via

This header is added to all requests in format according to RFC 9110.
If the request already contains this header, the new value is appended to it with `, ` (comma and space) as a separator.
The header is constructed from the request's protocol version (e.g. `HTTP/2.0`) and `os.Hostname()`.

## Caveats

- Since the request object provided into the middleware does not contain information about the server IP address nor
  about the server port, hostnames have to be used instead and unless something changes, it is impossible for the 
  `Forwarded` and `Via` headers to contain a full chain of IP addresses and ports.

## References

- RFC 7239: [https://datatracker.ietf.org/doc/rfc7239/](https://datatracker.ietf.org/doc/rfc7239/)
- RFC 9110: [https://datatracker.ietf.org/doc/rfc9110/](https://datatracker.ietf.org/doc/rfc9110/)

## README

This project addresses routing issues encountered when using applications proxied through Code Server. Code Server's proxy functionality often requires application-specific routing adjustments which aren't always feasible. This solution provides a mechanism to rewrite the `proxy/<port>` path to a subdomain, enabling access via a more conventional URL structure.

### Functionality

The Go script acts as a reverse proxy and URL rewriting service. It intercepts requests and performs the following:

1.  **Rewrite:** If a request is made to `proxy/<port>`, the script rewrites the URL to `<port>.domain.com` and redirects the client. For example, a request to `proxy.domain.com/proxy/9000` will redirect to `9000.domain.com`.

2.  **Subdomain Reverse Proxy:** If a request is made directly to a subdomain (e.g., `9000.domain.com`), the script forwards the request to a configurable target base URL with the port appended.

### Configuration

The script's behavior is controlled by the following environment variables:

- `DOMAIN_SUFFIX`: The domain suffix used for subdomain creation (e.g., `example.com`).
- `TARGET_BASE_URL`: The base URL of the target application (e.g., `http://localhost`).
- `SERVER_PORT`: The port on which the Go script listens for incoming requests.

### Example Usage

Assuming the following configuration:

- `DOMAIN_SUFFIX=example.com`
- `TARGET_BASE_URL=http://localhost`
- `SERVER_PORT=8081`

A request to `proxy.example.com/proxy/9000` will redirect to `9000.example.com`. Requests to `9000.example.com` will be proxied to `http://localhost:9000`.

### Integration with Code Server & Wildcard Subdomain Handling

This script is designed to be used in conjunction with Code Server. Configure Code Server to proxy requests to this script.To enable access via wildcard subdomains (e.g., `*.example.com`), configure your DNS to resolve all subdomain requests to the server running this script. The script then handles the URL rewriting and reverse proxying to the underlying applications. The following example demonstrates how a request to `proxy/<port>` is redirected to the wildcard subdomain: A request to `proxy.example.com/proxy/9000` is initially received by the server. The script then redirects the client to `9000.example.com`, effectively utilizing the wildcard DNS configuration.

### Security Considerations

This project focuses on routing and URL rewriting. Security measures, such as authentication and authorization, are the responsibility of the user and should be implemented separately.

### Caddy Configuration Example (Reference)

The following Caddyfile snippet demonstrates how to integrate this script with Caddy as a reverse proxy:

```caddyfile
*.example.com {
    reverse_proxy http://x.x.x.x:8081
}

coder.example.com {
    @proxyPath path_regexp proxyPath ^/proxy/(\d+)
    redir @proxyPath https://proxy.example.com{uri}

    reverse_proxy http://x.x.x.x:8080
}

```

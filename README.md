# Rename Header Traefik Plugin

Plugin to rename HTTP Headers in a Traefik Middleware.

## Configuration

### Static

```yml
experimental:
  plugins:
    renameHeaders:
      modulename: "github.com/beyerleinf/traefik-plugin-rename-header"
      version: "v1.0.0"
```

### Dynamic

```yml
http:
  routes:
    my-router:
      rule: "Host(`localhost`)"
      service: "my-service"
      middlewares:
        - "renameHeaders"
  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1"
  middlewares:
    renameHeaders:
      plugin:
        oldHeader: "X-Old"
        newHeader: "X-New"
```

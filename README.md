# Smart Hub 2 Exporter

Prometheus exporter for BT Smart Hub 2

## Configuration

1. Set the Smart Hub 2 IP in an environment variable if it is not using the default address (`192.168.1.254`)
```sh
SMARTHUB2_IP=ip_address # Optional: defaults to 192.168.1.254
```

## Docker (Compose)

```yaml
version: '3'
services:
  vox3-exporter:
    image: ghcr.io/njallam/smarthub2_exporter
    container_name: smarthub2_exporter
    restart: unless-stopped
    ports:
      - "9906:9906"
```

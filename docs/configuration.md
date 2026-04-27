# Configuration reference

NetManager stores all configuration in the `config/` directory relative to the binary. Files are created with empty defaults on first run if they do not exist.

---

## config/config.json

General application settings.

```json
{
  "service-name": "netmanager-service",
  "server-group-name": "dreammc"
}
```

| Key | Description | Example |
|---|---|---|
| `service-name` | Name of this NetManager instance | `netmanager-service` |
| `server-group-name` | Network group name — must match the `GROUP_NAME` env var used by deployed Paper services | `dreammc` |

---

## config/redis.json

Connection details for the Redis instance used for pub/sub communication with deployed services.

```json
{
  "host": "redis-host",
  "port": 6379,
  "username": "",
  "password": ""
}
```

| Key | Description |
|---|---|
| `host` | Redis hostname or IP |
| `port` | Redis port (default `6379`) |
| `username` | Redis ACL username (leave empty if not used) |
| `password` | Redis password (leave empty if not used) |

---

## config/harbor.json

Credentials for the Harbor container registry. Images built by NetManager are pushed here before being deployed to Kubernetes.

```json
{
  "domain": "registry.example.com",
  "project": "dreammc",
  "username": "robot$account",
  "password": "your-token"
}
```

| Key | Description | Example |
|---|---|---|
| `domain` | Harbor registry hostname | `registry.example.com` |
| `project` | Harbor project name | `dreammc` |
| `username` | Registry username (robot accounts recommended) | `robot$netmanager` |
| `password` | Registry password or token | |

Images are tagged as:

```
{domain}/{project}/{service-name}:{version}
```

---

## config/mongodb.json

MongoDB connection used when injecting the database URI into deployed service containers.

```json
{
  "uri": "mongodb://user:pass@mongo-host:27017"
}
```

| Key | Description |
|---|---|
| `uri` | Standard MongoDB connection string |

---

## Kubernetes access

NetManager detects the cluster connection automatically:

1. **In-cluster** — if running inside a pod, uses the mounted service-account token.
2. **Local** — falls back to `~/.kube/config`.

No additional configuration file is required. Ensure the service account (or your local kubeconfig user) has permissions to manage `StatefulSets`, `Services`, and `Namespaces` in the target cluster.

---

## Environment variables injected into deployed services

When NetManager starts a service on Kubernetes it injects the following environment variables into every container:

| Variable | Source |
|---|---|
| `MONGODB_URI` | Value from `config/mongodb.json` |
| `REDIS_URI` | Constructed from `config/redis.json` |
| `GROUP_NAME` | `server-group-name` from `config/config.json` |
| `SERVICE_NAME` | Service definition name |
| `SERVER_ID` | Replica index |

These match the variables expected by the [DreamMC API](../README.md#related-repositories) library.

---

## Data persistence

In addition to JSON config files, NetManager serialises service definitions to disk using Go's `gob` format. These files live alongside the JSON configs and are loaded automatically on startup. Do not edit them manually.

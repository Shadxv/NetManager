# Usage guide

## Starting NetManager

```bash
./bin/NetManager
```

On first run NetManager creates the required directory structure and empty config files, then drops into the interactive console.

---

## Interactive console

The console accepts commands while background tasks (builds, deployments) run concurrently. Output is colour-coded:

- **White** — general information
- **Yellow** — warnings / in-progress operations
- **Red** — errors
- **Green** — success

---

## Service commands

All service management is done through the `service` command group.

### `service create <name>`

Launches an interactive wizard that guides you through creating a new service definition.

```
> service create lobby
```

The wizard prompts for:

1. **Service type** — `paper` (game server) or `velocity` (proxy)
2. **Minecraft version** — fetched live from the PaperMC API
3. **Build number** — fetched live for the selected version
4. **Replica count** — number of pods to run (Paper only)
5. **Port** — exposed proxy port (Velocity only)

After completion, a service directory is created under `services/` containing a generated `Dockerfile` and the downloaded server `.jar`.

---

### `service list`

Displays all defined services with their type and current status.

```
> service list

┌────────────┬──────────┬──────────┐
│ NAME       │ TYPE     │ STATUS   │
├────────────┼──────────┼──────────┤
│ lobby      │ PAPER    │ running  │
│ proxy      │ VELOCITY │ running  │
│ minigames  │ PAPER    │ stopped  │
└────────────┴──────────┴──────────┘
```

---

### `service build <name>`

Builds a Docker image for the specified service and pushes it to the configured Harbor registry.

```
> service build lobby
```

Steps performed:

1. Reads the `Dockerfile` from the service directory
2. Builds the image via the local Docker daemon
3. Tags the image with an incremented version (`v0.1`, `v0.2`, …)
4. Authenticates with Harbor and pushes the image

The build runs in the background — console output is paused and resumed when the build completes.

---

### `service start <name>`

Deploys the service to the Kubernetes cluster.

```
> service start lobby
```

Creates the following Kubernetes resources:

- **Namespace** — `{server-group-name}` (created if it does not exist)
- **StatefulSet** — one pod per configured replica, using the latest pushed image
- **Service** — ClusterIP for internal communication; NodePort for Velocity proxies

Environment variables (`MONGODB_URI`, `REDIS_URI`, `GROUP_NAME`, `SERVICE_NAME`, `SERVER_ID`) are injected automatically. See [configuration.md](configuration.md#environment-variables-injected-into-deployed-services) for details.

---

### `service stop <name>`

Removes the Kubernetes StatefulSet and Service for the specified service.

```
> service stop lobby
```

This does **not** delete the service definition or built images — the service can be restarted with `service start` at any time.

---

### `service listpods <name>`

Lists all running pod instances for a service.

```
> service listpods lobby

┌───────────────┬─────────┬──────────┐
│ POD           │ READY   │ STATUS   │
├───────────────┼─────────┼──────────┤
│ lobby-0       │ 1/1     │ Running  │
│ lobby-1       │ 1/1     │ Running  │
└───────────────┴─────────┴──────────┘
```

---

## REST API

The embedded HTTP gateway is available on port `4000` while NetManager is running.

### List services

```http
GET /gateway/v1/services/
```

Returns a JSON array of all defined Paper and Velocity services.

**Example response:**

```json
[
  { "name": "lobby",     "type": "PAPER",    "status": "running" },
  { "name": "proxy",     "type": "VELOCITY", "status": "running" },
  { "name": "minigames", "type": "PAPER",    "status": "stopped" }
]
```

---

## Typical workflow

```
1. service create lobby          # define the service
2. service build lobby           # build and push image
3. service start lobby           # deploy to cluster
4. service listpods lobby        # verify pods are healthy
   ...
5. service stop lobby            # tear down when done
```

When you update a service (new Minecraft version, config change):

```
1. service build lobby           # increments image version
2. service stop lobby            # remove old deployment
3. service start lobby           # redeploy with new image
```

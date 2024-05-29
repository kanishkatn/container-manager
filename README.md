# Container Manager

Container Manager provides a service for managing containers. It includes functionalities for creating containers, enqueuing jobs, broadcasting jobs to a peer-to-peer network, and managing a job queue.

## Architecture

### JRPC API

The Container Manager provides a JRPC API for managing containers and jobs. The API includes the following methods:

- `CreateContainer`: Creates a container with the specified image. The job is queued, broadcast into the network and a job ID is returned.
- `Status`: Returns the status of the job with the specified ID.

Example usage:

Create request
```curl
curl -X POST localhost:8080/jrpc \
-H "Content-Type: application/json" \
-d '{
    "jsonrpc": "2.0",
    "method": "ContainerService.Create",
    "params": [{
        "image": "nginx",
        "arguments": [],
        "resources": {},
        "env": {}
    }],
    "id": 1
}'
```

On success, the response will include a job ID:

```json
{
  "result":{
    "job_id":"2c1581c9-1d82-11ef-aa1b-0242ac160003",
    "message":"Job created successfully"
  },
  "error":null,
  "id":1
}
```

Status request
```curl
curl -X POST localhost:8080/jrpc \
-H "Content-Type: application/json" \
-d '{
    "jsonrpc": "2.0",
    "method": "ContainerService.Status",
    "params": [{"job_id":"2c1581c9-1d82-11ef-aa1b-0242ac160003"}],
    "id": 1
}'
```

The response will include the status of the job:

```json
{
  "result":{
    "status":"running",
    "job_id":"2c1581c9-1d82-11ef-aa1b-0242ac160003"
  },
  "error":null,
  "id":1
}
```

### Job Queue Service

The Container Manager includes a job queue for managing jobs. The job queue is implemented using channels. The job queue includes the following methods:

```go
type Queue interface {
	Enqueue(jobID string, container types.Container) error
	GetStatus(jobID string) (types.JobStatus, bool)
	Run(workerCount int)
	Stop()
}
```

- `Enqueue`: Enqueues a job in the job queue.
- `GetStatus`: Returns the status of the job with the specified ID.
- `Run`: Runs the queue and processes the jobs.
- `Stop`: Stops the queue.

### Peer-to-Peer Service

The Container Manager includes a peer-to-peer service for broadcasting jobs to a peer-to-peer network. The peer-to-peer service includes the following methods:

```go
type P2PService interface {
	ID() string
	Start()
	Broadcast(msg Message) error
	Stop()
}
```

- `ID`: Returns the ID of the p2p host.
- `Start`: Starts the peer-to-peer service.
- `Broadcast`: Broadcasts a message to the peer-to-peer network.
- `Stop`: Stops the peer-to-peer service.

### Docker Service

The Container Manager includes a Docker service for managing Docker containers. The Docker service includes the following methods:

```go
type DockerService interface {
	DeployContainer(container types.Container) (string, error)
	GetContainerStatus(containerID string) (string, error)
}
```

- `DeployContainer`: Deploys a container with the specified image.
- `GetContainerStatus`: Returns the status of the container with the specified ID.

### CLI

The Container Manager includes a CLI for interacting with the Container Manager. The CLI is built using Cobra and includes only the root command.

```bash
Usage:
  container-manager [flags]

Flags:
  -h, --help                    help for container-manager
      --listen-address string   the address to listen on (default "0.0.0.0")
      --log-level string        log level (default "info")
      --port string             the port to listen on (default "8080")
      --queue-size int          the size of the job queue (default 100)
      --worker-count int        the number of workers to run (default 10)
```

## Usage

Build the Container Manager:

```bash
go build
```

Run the Container Manager:

```bash
./container-manager
```

## Testing

Run the tests:

```bash
go test ./...
```

The project also includes a Docker Compose file for running the tests in a Docker container, which spins up two containers. By default, the jrpc server
will be running on port 8080 and port 8081.

```bash
docker-compose up
```

Once the containers are up, send a request to the Container Manager:

```bash
curl -X POST localhost:8080/jrpc \
-H "Content-Type: application/json" \
-d '{
    "jsonrpc": "2.0",
    "method": "ContainerService.Create",
    "params": [{
        "image": "nginx",
        "arguments": [],
        "resources": {},
        "env": {}
    }],
    "id": 1
}'
```

Use the returned job ID to check the status of the job in both containers:

```bash
curl -X POST localhost:8080/jrpc \
-H "Content-Type: application/json" \
-d '{
    "jsonrpc": "2.0",
    "method": "ContainerService.Status",
    "params": [{"job_id":"2c1581c9-1d82-11ef-aa1b-0242ac160003"}],
    "id": 1
}'
```

```bash
curl -X POST localhost:8081/jrpc \
-H "Content-Type: application/json" \
-d '{
    "jsonrpc": "2.0",
    "method": "ContainerService.Status",
    "params": [{"job_id":"2c1581c9-1d82-11ef-aa1b-0242ac160003"}],
    "id": 1
}'
```

## Further Improvements

- Add integration tests for the JRPC API.
- Add integration tests for the p2p service.
- Better logging in the packages.
- Retry mechanism for the job queue.
- Viper support for configuration management.
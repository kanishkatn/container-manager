#!/bin/bash

# Start the containers
echo "Starting the local cluster..."
docker-compose up -d

# Wait until the service is up
wait_for_jrpc() {
  local container_name=$1
  local port=$2

  while ! docker exec "$container_name" curl -s "http://localhost:$port/jrpc" > /dev/null; do
    echo "Waiting for JRPC service in $container_name..."
    sleep 5
  done

  echo "JRPC service in $container_name is operational."
}

# Extract JSON value
extract_json_value() {
  local json=$1
  local key=$2
  echo "$json" | grep -o "\"$key\":\"[^\"]*\"" | sed "s/\"$key\":\"\([^\"]*\)\"/\1/"
}

# Check job status
check_job_status() {
  local container_name=$1
  local port=$2
  local jobID=$3

  response=$(docker exec "$container_name" curl -s -X POST -H "Content-Type: application/json" -d "{
      \"jsonrpc\": \"2.0\",
      \"method\": \"ContainerService.Status\",
      \"params\": [{\"job_id\":\"$jobID\"}],
      \"id\": 1
    }" "http://localhost:$port/jrpc")
    status=$(echo "$response" | grep -o '"status":"[^"]*"' | sed 's/"status":"\([^"]*\)"/\1/')
    echo "$status"
}

# Wait until the JRPC services are operational in both containers
echo "Waiting for the services to be operational..."
wait_for_jrpc manager1 8080
wait_for_jrpc manager2 8081

# Send a curl request and get jobID back
echo "Sending a request to create a container..."
response=$(curl -s -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"ContainerService.Create","params":[{"image": "nginx", "arguments": [], "env": {}}],"id":1}' http://localhost:8080/jrpc)
jobID=$(extract_json_value "$response" "job_id")

echo "Job ID: $jobID"

# Function to monitor job status
monitor_job() {
  local container_name=$1
  local port=$2
  local jobID=$3

  total_wait_time=120
  wait_interval=5
  elapsed_time=0

  while [ $elapsed_time -lt $total_wait_time ]; do
    status=$(check_job_status "$container_name" "$port" "$jobID")

    if [ "$status" == "complete" ]; then
      echo "Job $jobID completed successfully on $container_name."
      break
    fi

    echo "Waiting for $wait_interval seconds before checking again on $container_name..."
    sleep $wait_interval
    elapsed_time=$((elapsed_time + wait_interval))
  done

  if [ "$status" != "complete" ]; then
    echo "Error: Job $jobID did not complete within 2 minutes on $container_name."
  fi
}

# Monitor the job status on both managers
echo "Monitoring the job status for manager1"
monitor_job manager1 8080 "$jobID"

echo "Monitoring the job status for manager2"
monitor_job manager2 8081 "$jobID"

# Stop and cleanup the containers
docker-compose down

container_ids=$(docker ps -a -q --filter ancestor=nginx)
if [ -z "$container_ids" ]; then
  exit 0
fi

for container_id in $container_ids; do
  docker stop "$container_id"
  docker rm "$container_id"
done

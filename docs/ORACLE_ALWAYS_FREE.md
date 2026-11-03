# Oracle Cloud Always Free (ARM64) deployment

This is a **single-host demonstration**, not a production HA deployment. Always Free capacity is not guaranteed. Keep the instance within the tenancy's Always Free quotas and verify the billing screen before provisioning.

## 1. Provision
Create an Oracle Cloud Ampere A1 Flex Ubuntu ARM64 VM within Always Free limits (for example 2 OCPUs / 12 GB RAM), with an SSH key and a boot volume that fits the remaining free storage allowance. Use a public IP for SSH. In OCI VCN security rules allow **only TCP 22** from your IP; do not expose database, Kafka, Grafana, Jaeger or Prometheus ports.

## 2. Install Docker on the VM
Follow the current official Docker Engine Ubuntu instructions at https://docs.docker.com/engine/install/ubuntu/ and install the Compose plugin. Confirm `docker compose version`. Enable Docker to start on boot.

## 3. Clone and configure
```bash
git clone https://github.com/jenildhorajiya/cqrs_microservice.git
cd cqrs_microservice
git checkout deploy/oracle-always-free
cp .env.production.example .env
chmod 600 .env
# Edit .env and set three unique, strong passwords.
nano .env
docker compose -f docker-compose.prod.yaml config --quiet
docker compose -f docker-compose.prod.yaml up -d --build
docker compose -f docker-compose.prod.yaml ps
docker compose -f docker-compose.prod.yaml logs --tail=100 gateway writer reader kafka
```
The Compose CLI automatically reads `.env` in the working directory. Do not commit it.

## 4. Initialize PostgreSQL
The repository contains SQL migrations in `migrations/`. Use the golang-migrate CLI from the VM or an ephemeral container on the Compose network. Supply the PostgreSQL password securely, and run the migration against the `products` database **before** exercising product APIs. Do not run the old Makefile migration command unchanged: it assumes localhost and development credentials.

## 5. Initialize MongoDB
Inspect `scripts/init.js` before using it against MongoDB 7. The MongoDB container creates its root user from the environment. If the script is required to create collections or indexes, execute it with `mongosh` and the configured admin credentials.

## 6. Validate and access
On the VM, check `curl -i http://127.0.0.1:5001/live` only if your gateway exposes the liveness probe on the HTTP listener; in this project probes are instead on port 3001 inside the container. For API docs check `curl -I http://127.0.0.1:5001/swagger/index.html`.

From your laptop, create an SSH tunnel (replace IP and key):
```bash
ssh -i ~/.ssh/oci_key -N \
  -L 5001:127.0.0.1:5001 \
  -L 3005:127.0.0.1:3005 \
  -L 9090:127.0.0.1:9090 \
  -L 16686:127.0.0.1:16686 ubuntu@YOUR_VM_PUBLIC_IP
```
Then visit:
- http://localhost:5001/swagger/index.html
- http://localhost:3005 (Grafana)
- http://localhost:9090 (Prometheus)
- http://localhost:16686 (Jaeger)

Add a Grafana Prometheus data source at `http://prometheus:9090` (the Docker-internal address).

## 7. Public HTTPS endpoint
The gateway is intentionally bound to loopback. For a public API, configure a domain, HTTPS reverse proxy such as Caddy/Nginx, authentication/rate limiting, and explicitly open only ports 80/443 in OCI and the host firewall. Never publish the data services or admin UIs.

## 8. Operations
```bash
docker compose -f docker-compose.prod.yaml logs -f --tail=100
docker compose -f docker-compose.prod.yaml ps
docker compose -f docker-compose.prod.yaml pull
docker compose -f docker-compose.prod.yaml up -d --build
```
Persisted volumes hold database and monitoring state. Back them up before upgrades. Watch VM memory/disk and configure an OCI budget alert. This project uses single-node Kafka, databases, and monitoring without redundancy.

## Notes
- ARM64 images should be verified on the VM; image availability and free quotas may change.
- Legacy Go dependencies and service initialization may require modernization; verify all three Go services build and the create-read event flow before treating deployment as complete.
- The original development Compose files are unchanged.

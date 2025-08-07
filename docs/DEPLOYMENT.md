# RouteIQ Deployment Guide

This guide covers local development, staging, and production deployments on Google Cloud Run with supporting GCP services.

## Prerequisites
- Google Cloud project with billing enabled
- gcloud CLI installed and authenticated
- Docker installed
- Terraform (optional, recommended for IaC)

## Services and Resources
- Cloud Run (containers)
- Artifact Registry (container registry)
- Pub/Sub (messaging)
- AlloyDB for PostgreSQL (primary database)
- BigQuery (analytics)
- Secret Manager (secrets)
- Cloud Logging/Monitoring/Trace (observability)

## Environment Configuration
Define environment variables for each service:

- API
  - `PORT` (default 8080)
  - `DATABASE_URL` (Postgres connection string)
  - `PUBSUB_TOPIC`, `PUBSUB_SUBSCRIPTION`
  - `JWT_AUDIENCE`, `JWT_ISSUER` (if using auth between services)
  - `ALLOYDB_INSTANCE_URI` (for auth proxy: `projects/PROJECT/locations/REGION/clusters/CLUSTER/instances/INSTANCE`)
- Frontend
  - `NEXT_PUBLIC_MAPBOX_TOKEN`
  - `WS_URL` (WebSocket URL)
- Shared
  - `GCP_PROJECT_ID`

Store secrets in Secret Manager and mount them as env vars or files in Cloud Run.

## Local Development
- Database: use Dockerized PostgreSQL locally
- Pub/Sub: use the Pub/Sub emulator
- BigQuery: mock locally; run in staging/prod

Example docker-compose stubs (to be added when services exist):
```yaml
version: '3.9'
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: routeiq
      POSTGRES_PASSWORD: routeiq
      POSTGRES_DB: routeiq
    ports:
      - "5432:5432"
  pubsub:
    image: gcr.io/google.com/cloudsdktool/cloud-sdk:slim
    command: gcloud beta emulators pubsub start --host-port=0.0.0.0:8085
    ports:
      - "8085:8085"
```

## Container Build and Push (Artifact Registry)
```bash
REGION=us-central1
PROJECT_ID=YOUR_PROJECT_ID
REPO=routeiq
SERVICE=api

gcloud artifacts repositories create $REPO \
  --repository-format=docker \
  --location=$REGION || true

gcloud auth configure-docker $REGION-docker.pkg.dev --quiet

docker build -t $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$SERVICE:latest ./backend/api

docker push $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$SERVICE:latest
```

## Deploy to Cloud Run
```bash
REGION=us-central1
PROJECT_ID=YOUR_PROJECT_ID
REPO=routeiq
SERVICE=api
IMAGE=$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$SERVICE:latest

gcloud run deploy $SERVICE \
  --image $IMAGE \
  --region $REGION \
  --platform managed \
  --allow-unauthenticated \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10 \
  --set-env-vars "PORT=8080" \
  --set-secrets "DATABASE_URL=projects/$PROJECT_ID/secrets/routeiq-db-url:latest"
```

Repeat for additional services (websocket, exporter, frontend). For private APIs, remove `--allow-unauthenticated` and configure IAM.

## Database Setup
- Provision AlloyDB (cluster + primary instance) with Private IP
  - Create VPC and subnetwork for AlloyDB
  - Create cluster and primary instance in your region
  - Create database and user
- Connectivity options from Cloud Run:
  1) Serverless VPC Access + Private IP (recommended for production)
  2) AlloyDB Auth Proxy sidecar (simple to start; supports IAM DB Auth)
- Apply schemas using a migration tool (e.g., `golang-migrate`)

```bash
migrate -path db/migrations -database "$DATABASE_URL" up
```

## AlloyDB Connectivity

### Option A: Serverless VPC Access (Private IP)
1. Create a Serverless VPC Access connector in the same region as Cloud Run.
```bash
gcloud compute networks vpc-access connectors create routeiq-connector \
  --region=$REGION \
  --network=default \
  --range=10.8.0.0/28
```
2. Deploy Cloud Run with VPC connector and private egress.
```bash
gcloud run deploy $SERVICE \
  --image $IMAGE \
  --region $REGION \
  --allow-unauthenticated \
  --vpc-connector routeiq-connector \
  --vpc-egress private-ranges-only \
  --set-env-vars "DATABASE_URL=postgres://USER:PASSWORD@<ALLOYDB_PRIVATE_IP>:5432/routeiq?sslmode=disable"
```

### Option B: AlloyDB Auth Proxy (local + Cloud Run sidecar)

Local development:
```bash
alloydb-auth-proxy \
  projects/$PROJECT_ID/locations/$REGION/clusters/$CLUSTER/instances/$INSTANCE \
  --port 5432

export DATABASE_URL="postgres://USER:PASSWORD@127.0.0.1:5432/routeiq?sslmode=disable"
```

Cloud Run with proxy sidecar (multi-container):
```bash
gcloud run deploy $SERVICE \
  --region $REGION \
  --allow-unauthenticated \
  --image $IMAGE \
  --container $SERVICE \
  --set-env-vars "DATABASE_URL=postgres://USER:PASSWORD@127.0.0.1:5432/routeiq?sslmode=disable" \
  --add-container alloydb-proxy=gcr.io/cloud-sql-connectors/alloydb-auth-proxy:latest \
  --container alloydb-proxy \
  --container-command alloydb-auth-proxy \
  --container-args "projects/$PROJECT_ID/locations/$REGION/clusters/$CLUSTER/instances/$INSTANCE","--port=5432"
```

Notes:
- For IAM DB Auth, add `--enable_iam_login` to the proxy and use short-lived tokens.
- Ensure the Cloud Run service account has `roles/alloydb.client`.

## Pub/Sub Topics
```bash
gcloud pubsub topics create traffic-events
gcloud pubsub subscriptions create traffic-events-sub \
  --topic=traffic-events
```

## BigQuery
- Create dataset `routeiq`
- Use time-partitioned tables for historical traffic

```bash
gcloud bq --location=US mk --dataset $PROJECT_ID:routeiq
```

## Observability
- Structured JSON logs with correlation IDs
- Export metrics to Cloud Monitoring
- Enable Cloud Trace for request profiling

## Security
- Enforce HTTPS
- Least-privilege IAM for service accounts
- Secrets in Secret Manager only
- VPC for private DB connectivity (Serverless VPC Access)

## Rollouts
- Use traffic-splitting for zero-downtime deploys
- Blue/Green by keeping N-1 revision available

## Cost Notes (MVP)
- Prefer Cloud SQL over AlloyDB for MVP to reduce complexity and cost
- Use min instances = 0 for Cloud Run to scale to zero when idle

---
For production hardening, add Terraform to manage all resources and set up CI/CD to build, test, and deploy on merge to `main`.

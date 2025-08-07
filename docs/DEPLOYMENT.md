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
- Cloud SQL for PostgreSQL (primary database; upgrade to AlloyDB later)
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
- Start with Cloud SQL for PostgreSQL for MVP
- Apply schemas using a migration tool (e.g., `golang-migrate`)

```bash
migrate -path db/migrations -database "$DATABASE_URL" up
```

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

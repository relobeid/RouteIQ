#!/usr/bin/env bash
set -euo pipefail

# RouteIQ GCP bootstrap (AlloyDB-first)
# Usage:
#   export PROJECT_ID=your-project
#   export REGION=us-central1
#   export NETWORK=default
#   export ARTIFACT_REPO=routeiq
#   export PUBSUB_TOPIC=traffic-events
#   export PUBSUB_SUB=traffic-events-sub
#   export BQ_DATASET=routeiq
#   export CLUSTER=routeiq-cluster
#   export INSTANCE=routeiq-primary
#   export ALLOYDB_POSTGRES_PASSWORD='strong-password'
#   bash scripts/gcp_bootstrap.sh

: "${PROJECT_ID:?PROJECT_ID is required}"
: "${REGION:?REGION is required}"
: "${NETWORK:?NETWORK is required}"
: "${ARTIFACT_REPO:?ARTIFACT_REPO is required}"
: "${PUBSUB_TOPIC:?PUBSUB_TOPIC is required}"
: "${PUBSUB_SUB:?PUBSUB_SUB is required}"
: "${BQ_DATASET:?BQ_DATASET is required}"
: "${CLUSTER:?CLUSTER is required}"
: "${INSTANCE:?INSTANCE is required}"
: "${ALLOYDB_POSTGRES_PASSWORD:?ALLOYDB_POSTGRES_PASSWORD is required}"

SA_NAME=routeiq-run-sa
SA_EMAIL="$SA_NAME@$PROJECT_ID.iam.gserviceaccount.com"
CONNECTOR_NAME=routeiq-connector

log() { echo "[routeiq] $*"; }

log "Setting project"
gcloud config set project "$PROJECT_ID"

log "Enabling required services"
gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  pubsub.googleapis.com \
  alloydb.googleapis.com \
  secretmanager.googleapis.com \
  bigquery.googleapis.com \
  vpcaccess.googleapis.com \
  cloudbuild.googleapis.com \
  iam.googleapis.com

log "Creating Artifact Registry (if missing)"
gcloud artifacts repositories create "$ARTIFACT_REPO" \
  --repository-format=docker \
  --location="$REGION" || true

gcloud auth configure-docker "$REGION-docker.pkg.dev" --quiet

log "Creating service account: $SA_EMAIL"
gcloud iam service-accounts create "$SA_NAME" --display-name "RouteIQ Cloud Run SA" || true

log "Granting roles to service account"
# Runtime roles for Cloud Run service to access dependencies
for role in \
  roles/alloydb.client \
  roles/secretmanager.secretAccessor \
  roles/pubsub.publisher \
  roles/pubsub.subscriber \
  roles/logging.logWriter \
  roles/monitoring.metricWriter; do
  gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:$SA_EMAIL" \
    --role="$role" --quiet
done

log "Creating Pub/Sub topic and subscription"
gcloud pubsub topics create "$PUBSUB_TOPIC" || true
gcloud pubsub subscriptions create "$PUBSUB_SUB" \
  --topic="$PUBSUB_TOPIC" || true

log "Creating BigQuery dataset"
bq --location=US mk --dataset "$PROJECT_ID:$BQ_DATASET" || true

log "Creating Serverless VPC Access connector"
gcloud compute networks vpc-access connectors create "$CONNECTOR_NAME" \
  --region="$REGION" \
  --network="$NETWORK" \
  --range=10.8.0.0/28 || true

log "Provisioning AlloyDB cluster"
gcloud beta alloydb clusters create "$CLUSTER" \
  --region="$REGION" \
  --network="projects/$PROJECT_ID/global/networks/$NETWORK" \
  --password="$ALLOYDB_POSTGRES_PASSWORD" || true

log "Provisioning AlloyDB primary instance"
gcloud beta alloydb instances create "$INSTANCE" \
  --region="$REGION" \
  --cluster="$CLUSTER" \
  --instance-type=PRIMARY \
  --cpu-count=2 || true

log "Fetching AlloyDB instance IP"
ALLOYDB_IP=$(gcloud beta alloydb instances describe "$INSTANCE" \
  --region="$REGION" --cluster="$CLUSTER" \
  --format='value(ipAddress)')
if [[ -z "$ALLOYDB_IP" ]]; then
  echo "Failed to resolve AlloyDB IP. Check instance status." >&2
  exit 1
fi
log "AlloyDB IP: $ALLOYDB_IP"

log "Creating Secret Manager secret for DATABASE_URL"
DB_URL="postgres://postgres:${ALLOYDB_POSTGRES_PASSWORD}@${ALLOYDB_IP}:5432/routeiq?sslmode=disable"
if gcloud secrets describe routeiq-db-url >/dev/null 2>&1; then
  echo -n "$DB_URL" | gcloud secrets versions add routeiq-db-url --data-file=-
else
  echo -n "$DB_URL" | gcloud secrets create routeiq-db-url --data-file=-
fi

log "Bootstrap complete. Next steps:"
echo "1) Build and push your container to $REGION-docker.pkg.dev/$PROJECT_ID/$ARTIFACT_REPO/<service>:latest"
echo "2) Deploy Cloud Run with:"
echo "   gcloud run deploy <service> \\\" 
      --image $REGION-docker.pkg.dev/$PROJECT_ID/$ARTIFACT_REPO/<service>:latest \\\" 
      --region $REGION \\\" 
      --allow-unauthenticated \\\" 
      --service-account $SA_EMAIL \\\" 
      --vpc-connector $CONNECTOR_NAME \\\" 
      --vpc-egress private-ranges-only \\\" 
      --set-secrets DATABASE_URL=projects/$PROJECT_ID/secrets/routeiq-db-url:latest"

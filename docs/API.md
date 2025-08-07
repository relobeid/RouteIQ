# RouteIQ API Reference (MVP)

This document outlines the initial API surface for ingestion, route optimization, and realtime updates.

## Conventions
- Base URL: varies by environment (e.g., Cloud Run URL)
- Auth: Bearer JWT for service-to-service (optional in MVP)
- Content-Type: `application/json`
- Error format:
```json
{
  "error": {
    "code": "string",   
    "message": "string",
    "details": {}
  }
}
```

## 1. Traffic Ingestion
Ingest simulated or real traffic events (vehicle updates, incidents).

### POST /api/v1/traffic/vehicle
- Description: Upsert vehicle state
- Body:
```json
{
  "id": "uuid",
  "position": {"x": 10, "y": 15},
  "speed": 25.5,
  "destination": {"x": 19, "y": 4},
  "timestamp": "2024-01-15T10:30:00Z"
}
```
- Responses:
  - 202 Accepted (enqueued)
  - 400 Invalid payload

### POST /api/v1/traffic/incident
- Description: Record or update incident
- Body:
```json
{
  "id": "uuid",
  "type": "accident | closure | construction",
  "position": {"x": 3, "y": 11},
  "severity": 1,
  "timestamp": "2024-01-15T10:30:00Z",
  "resolved": false
}
```
- Responses:
  - 202 Accepted

## 2. Route Optimization

### POST /api/v1/routes/optimal
- Description: Compute optimal route considering live conditions
- Body:
```json
{
  "origin": {"x": 0, "y": 0},
  "destination": {"x": 19, "y": 19},
  "preferences": {
    "avoid_incidents": true,
    "weight": 1.2
  }
}
```
- Response:
```json
{
  "route": {
    "path": [{"x":0,"y":0}, {"x":0,"y":1}],
    "distance": 42.0,
    "eta_seconds": 520
  },
  "metadata": {
    "computed_ms": 120,
    "traffic_multipliers": [1.0, 1.3]
  }
}
```

## 3. Realtime Updates (WebSocket)
- URL: `wss://<host>/ws`
- Heartbeat: ping/pong every 30s
- Messages:
```json
{
  "type": "vehicle_update | incident_update | stats",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {}
}
```

## 4. Health and Metrics

### GET /healthz
- 200 OK

### GET /metrics
- Prometheus metrics endpoint (authenticated or private)

---

Future: publish OpenAPI spec when endpoints solidify; include rate limits and pagination for analytics endpoints.

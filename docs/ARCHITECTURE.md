# RouteIQ System Architecture

## Architecture Overview

RouteIQ follows a microservices architecture pattern optimized for real-time data processing and machine learning inference. The system is designed to handle high-volume traffic data ingestion while providing low-latency route optimization recommendations.

## Core Components

### 1. Traffic Simulation Engine
**Technology**: Go  
**Purpose**: Generate realistic traffic patterns for testing and training

- **Grid System**: 20x20 coordinate system representing city blocks
- **Vehicle Simulation**: Physics-based movement with acceleration/braking
- **Intersection Management**: Traffic light cycles and turn restrictions
- **Incident Generation**: Random accidents and road closures for realistic scenarios

**Data Output**: Vehicle positions, speeds, destinations, and traffic events

### 2. Data Ingestion Layer

#### Go API Server
**Technology**: Go with Gorilla Mux router  
**Purpose**: High-throughput traffic data processing

- **Concurrent Processing**: Goroutines for handling 10,000+ requests/hour
- **Data Validation**: Input sanitization and format verification
- **Rate Limiting**: Protection against traffic spikes
- **Health Monitoring**: Prometheus metrics and health checks

#### Message Queue
**Technology**: Google Cloud Pub/Sub  
**Purpose**: Decouple data ingestion from processing

- **Async Processing**: Non-blocking data flow between services
- **Fault Tolerance**: Message persistence and retry mechanisms
- **Scalability**: Automatic scaling based on queue depth

### 3. Data Storage Layer

#### AlloyDB (Primary Database)
**Technology**: Google AlloyDB (PostgreSQL-compatible)  
**Purpose**: Transactional storage for live traffic state

**Schema Design**:
```sql
-- Vehicle tracking table
CREATE TABLE vehicles (
    id UUID PRIMARY KEY,
    position_x INTEGER NOT NULL,
    position_y INTEGER NOT NULL,
    speed DECIMAL(5,2),
    destination_x INTEGER,
    destination_y INTEGER,
    last_updated TIMESTAMP DEFAULT NOW()
);

-- Intersection state table
CREATE TABLE intersections (
    id INTEGER PRIMARY KEY,
    position_x INTEGER NOT NULL,
    position_y INTEGER NOT NULL,
    light_state VARCHAR(10) CHECK (light_state IN ('red', 'yellow', 'green')),
    last_changed TIMESTAMP DEFAULT NOW()
);

-- Traffic incidents table
CREATE TABLE incidents (
    id UUID PRIMARY KEY,
    incident_type VARCHAR(20),
    position_x INTEGER NOT NULL,
    position_y INTEGER NOT NULL,
    severity INTEGER CHECK (severity BETWEEN 1 AND 5),
    created_at TIMESTAMP DEFAULT NOW(),
    resolved_at TIMESTAMP
);
```

**Performance Optimizations**:
- Spatial indexing for position-based queries
- Connection pooling for concurrent access
- Read replicas for analytics queries

#### BigQuery (Analytics Warehouse)
**Technology**: Google BigQuery  
**Purpose**: Historical analysis and ML training data

**Data Pipeline**:
- Batch export from AlloyDB every 15 minutes
- Time-partitioned tables for efficient querying
- Aggregated views for common analytics patterns

### 4. Machine Learning Layer

#### VertexAI Integration
**Technology**: Google VertexAI  
**Purpose**: Traffic prediction and anomaly detection

**Model Architecture**:
- **Time Series Forecasting**: LSTM networks for traffic volume prediction
- **Classification Models**: Random Forest for incident type prediction
- **Feature Engineering**: Time-based, location-based, and historical pattern features

**Training Pipeline**:
```python
# Feature extraction from historical data
features = [
    'hour_of_day', 'day_of_week', 'month',
    'intersection_id', 'weather_condition',
    'historical_avg_volume', 'recent_incidents'
]

# Model training configuration
model_config = {
    'algorithm': 'automl_tabular',
    'prediction_type': 'regression',
    'target_column': 'traffic_volume',
    'optimization_objective': 'minimize_rmse'
}
```

### 5. Route Optimization Engine

#### Weighted A* Algorithm
**Technology**: Go with custom graph implementation  
**Purpose**: Dynamic route calculation using live traffic conditions

**Algorithm Features**:
- **Dynamic Weights**: Real-time traffic conditions affect edge costs
- **Heuristic Function**: Manhattan distance with traffic congestion factors
- **Caching**: Pre-computed routes for common origin-destination pairs
- **Optimization**: Sub-500ms response time for route calculations

**Weight Calculation**:
```go
func calculateEdgeWeight(from, to Position, currentTime time.Time) float64 {
    baseDistance := manhattanDistance(from, to)
    trafficMultiplier := getTrafficDensity(from, to, currentTime)
    incidentPenalty := getIncidentPenalty(from, to)
    
    return baseDistance * trafficMultiplier * incidentPenalty
}
```

### 6. Real-time Communication Layer

#### WebSocket Server
**Technology**: Go with Gorilla WebSocket  
**Purpose**: Live updates to frontend dashboard

**Connection Management**:
- Connection pooling for multiple dashboard clients
- Heartbeat mechanism for connection health
- Graceful degradation for connection failures
- Rate limiting to prevent client overload

**Message Protocol**:
```json
{
    "type": "vehicle_update",
    "timestamp": "2024-01-15T10:30:00Z",
    "data": {
        "vehicle_id": "uuid",
        "position": {"x": 10, "y": 15},
        "speed": 25.5,
        "status": "moving"
    }
}
```

### 7. Frontend Dashboard

#### Next.js Application
**Technology**: Next.js 14 with TypeScript  
**Purpose**: Interactive traffic visualization and analytics

**Component Architecture**:
- **Map Component**: Mapbox GL JS integration for traffic visualization
- **Analytics Panel**: Real-time statistics and historical trends
- **Route Planner**: Interactive route selection and optimization
- **Admin Dashboard**: System monitoring and configuration

#### Mapbox Integration
**Technology**: Mapbox GL JS  
**Purpose**: Interactive traffic mapping

**Visualization Features**:
- **Heat Maps**: Traffic density visualization using color gradients
- **Vehicle Markers**: Real-time vehicle position updates
- **Route Overlays**: Optimized route recommendations
- **Incident Markers**: Traffic incident locations and severity

### 8. Infrastructure & Deployment

#### Containerization
**Technology**: Docker with multi-stage builds  
**Purpose**: Consistent deployment across environments

**Container Strategy**:
```dockerfile
# Go API container
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

#### Cloud Deployment
**Technology**: Google Cloud Run  
**Purpose**: Serverless container hosting

**Deployment Configuration**:
- Auto-scaling based on request volume
- Zero-downtime deployments with traffic splitting
- Environment-specific configurations
- Integrated logging and monitoring

## Data Flow Architecture

### Real-time Data Flow
1. **Traffic Simulation** generates vehicle events at 1Hz frequency
2. **Go API** receives and validates traffic data
3. **Pub/Sub** queues messages for async processing
4. **AlloyDB** stores current traffic state with sub-100ms latency
5. **WebSocket Server** broadcasts updates to connected clients
6. **Next.js Dashboard** renders real-time traffic visualization

### Analytics Data Flow
1. **AlloyDB** accumulates historical traffic data
2. **Batch Export** transfers data to BigQuery every 15 minutes
3. **BigQuery** processes and aggregates historical patterns
4. **VertexAI** trains ML models on historical data
5. **Prediction API** serves traffic forecasts for route optimization

### Route Optimization Flow
1. **User Request** for route from origin to destination
2. **Current Traffic Data** retrieved from AlloyDB
3. **ML Predictions** fetched for forecasted conditions
4. **A* Algorithm** calculates optimal route with weighted edges
5. **Route Response** returned with estimated travel time

## Performance Characteristics

### Latency Requirements
- **API Response Time**: <200ms for data ingestion
- **WebSocket Updates**: <500ms from data generation to client
- **Route Calculation**: <500ms for optimal path finding
- **ML Predictions**: <1000ms for traffic forecasting

### Throughput Targets
- **Data Ingestion**: 10,000+ traffic data points per hour
- **Concurrent Users**: 100+ simultaneous dashboard connections
- **Database Operations**: 1000+ writes per minute
- **Route Calculations**: 100+ per minute

### Scalability Design
- **Horizontal Scaling**: Stateless services with load balancing
- **Database Scaling**: Read replicas and connection pooling
- **Caching Strategy**: Redis for frequently accessed route calculations
- **CDN Integration**: Static asset delivery optimization

## Security & Monitoring

### Security Measures
- **API Authentication**: JWT tokens for service-to-service communication
- **Network Security**: VPC with private subnets for databases
- **Data Encryption**: TLS 1.3 for all external communications
- **Input Validation**: Comprehensive sanitization of traffic data

### Monitoring & Observability
- **Application Metrics**: Prometheus for custom metrics collection
- **Infrastructure Monitoring**: Google Cloud Monitoring for resource usage
- **Distributed Tracing**: Cloud Trace for request flow analysis
- **Error Tracking**: Structured logging with correlation IDs

This architecture ensures RouteIQ can handle the demands of real-time traffic management while maintaining the flexibility to scale and evolve with changing requirements.
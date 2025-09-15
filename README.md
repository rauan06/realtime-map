<p align="center">
    <img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/ec559a9f6bfd399b82bb44393651661b08aaf7ba/icons/folder-markdown-open.svg" align="center" width="30%">
</p>
<p align="center"><h1 align="center">REALTIME-MAP</h1></p>
<p align="center">
	<em><code>Realtime scalable map using Golang and Kafka. Other tools: Redis, gRPC, WebSocket, GORM, PostgreSQL</code></em>
</p>
<p align="center">
	<img src="https://img.shields.io/github/license/rauan06/realtime-map?style=default&logo=opensourceinitiative&logoColor=white&color=0080ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/rauan06/realtime-map?style=default&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/rauan06/realtime-map?style=default&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/rauan06/realtime-map?style=default&color=0080ff" alt="repo-language-count">
</p>
<p align="center"><!-- default option, no dependency badges. -->
</p>
<p align="center">
	<!-- default option, no dependency badges. -->
</p>
<br>

##  Table of Contents

- [ Overview](#-overview)
- [ Features](#-features)
- [ Project Structure](#-project-structure)
- [ Architecture](#-architecture)
- [ Getting Started](#-getting-started)
  - [ Prerequisites](#-prerequisites)
  - [ Installation](#-installation)
  - [ Usage](#-usage)
  - [ Testing](#-testing)
- [ API Documentation](#-api-documentation)
- [ Project Roadmap](#-project-roadmap)
- [ Contributing](#-contributing)
- [ License](#-license)

---

##  Overview

**Realtime Map** is a scalable, distributed system for real-time location tracking and visualization. The system processes location data from IoT devices (OBU - On-Board Units) through a Kafka-based streaming pipeline, stores it in Redis for fast retrieval, and provides real-time visualization through WebSocket connections.

The architecture leverages microservices pattern with separate components for data ingestion, processing, storage, and presentation. Built with Go, it provides high-performance real-time location tracking suitable for fleet management, logistics, and IoT applications.

---

##  Features

- **Real-time Location Tracking**: Process GPS coordinates from multiple devices simultaneously
- **Scalable Architecture**: Microservices-based design with horizontal scaling capabilities
- **High-Performance Data Processing**: Kafka-based streaming for handling high-throughput location data
- **Fast Data Retrieval**: Redis clustering for sub-millisecond location lookups
- **Real-time Visualization**: WebSocket-powered live map updates using Leaflet.js
- **gRPC API**: High-performance API for device communication and data exchange
- **RESTful Gateway**: HTTP/REST API gateway for easy integration
- **Multi-Device Support**: Handle thousands of concurrent device connections
- **Data Persistence**: PostgreSQL for historical data storage and analytics
- **Containerized Deployment**: Docker Compose setup for easy development and deployment

---

##  Project Structure

```scharp
Devices -> Kafka Topic (obu_positions)  ──┐
                                         │
                                 ┌───────▼────────┐
                                 │  Consumer Pool │  (group of Go workers)
                                 └───────┬────────┘
                                         │
         ┌───────────────────────────────┼──────────────────────────────┐
         │                               │                              │
   Update Redis shard A             Update Redis shard B           Publish to Fan-out
   key: obu:{device_id}              key: obu:{device_id}            (optional)
   (Redis Cluster)                   (Redis Cluster)                 └─> internal pubsub (Redis/ NATS)
                                                                   ┌─> gRPC/WebSocket Servers (stateless, horizontal)
                                                                   └─> Optional: Kafka topic for derived events
Clients connect:
  1) On connect: gRPC server reads snapshot from Redis (get keys / mget / scan by hash tag)
  2) After snapshot: client subscribes to live updates (server pushes) or connects to stream

```

##  Architecture

### System Components

- **Seeder Service**: Generates simulated GPS data from multiple devices for testing
- **Producer Service**: Consumes location data from Kafka, processes it, and stores in Redis
- **API Gateway**: Provides gRPC, REST, and WebSocket APIs with web interface
- **Go Commons**: Shared protocol buffers, database utilities, and common libraries

### Data Flow

1. **Data Ingestion**: Devices send GPS coordinates via gRPC to API Gateway
2. **Message Queuing**: Location data is published to Kafka topic `obu_positions`
3. **Stream Processing**: Producer service consumes from Kafka and processes location updates
4. **Caching**: Current device positions are stored in Redis for fast retrieval
5. **Real-time Updates**: WebSocket connections push live updates to web clients
6. **Persistence**: Historical data is stored in PostgreSQL for analytics

### Technology Stack

- **Backend**: Go 1.21+, gRPC, Protocol Buffers
- **Message Broker**: Apache Kafka for event streaming
- **Caching**: Redis cluster for high-performance data access  
- **Database**: PostgreSQL for persistent storage
- **Frontend**: HTML5, JavaScript, Leaflet.js for map visualization
- **Infrastructure**: Docker, Docker Compose for containerization

##  Getting Started

###  Prerequisites

Before getting started with realtime-map, ensure your runtime environment meets the following requirements:

- **Programming Language:** Go 1.21 or higher
- **Container Runtime:** Docker and Docker Compose
- **Protocol Buffers:** Protocol Buffers compiler (protoc)
- **Build Tools:** Make utility
- **Message Broker:** Apache Kafka (provided via Docker)
- **Database:** PostgreSQL (provided via Docker)
- **Cache:** Redis (provided via Docker)


###  Installation

Install realtime-map using the following steps:

**1. Clone the Repository:**
```sh
❯ git clone https://github.com/rauan06/realtime-map
❯ cd realtime-map
```

**2. Install Development Dependencies:**
```sh
❯ make deps
❯ make buf-install
```

**3. Generate Protocol Buffer Files:**
```sh
❯ make generate
```

**4. Start Infrastructure Services:**
```sh
❯ make up
# or manually: docker compose up -d
```

This will start:
- Apache Kafka broker on port 9092
- Redis cache on port 6378
- PostgreSQL database on port 5431

**5. Build Services (Optional):**
```sh
❯ go build ./api-gateway/cmd/app
❯ go build ./producer/cmd/app  
❯ go build ./seeder/cmd/app
```




###  Usage

The realtime-map system consists of multiple microservices that work together:

**1. Start the Infrastructure:**
```sh
❯ make up
```

**2. Run the API Gateway:**
```sh
❯ cd api-gateway && go run cmd/app/main.go
```
The API gateway provides:
- gRPC server for device communication
- HTTP gateway for REST API access  
- WebSocket endpoints for real-time map updates
- Web interface at http://localhost:8080

**3. Run the Producer Service:**
```sh
❯ cd producer && go run cmd/app/main.go
```
The producer consumes location data from Kafka and processes it for storage and real-time distribution.

**4. Run the Data Seeder (for testing):**
```sh
❯ cd seeder && go run cmd/app/main.go
```
The seeder generates simulated device location data for testing and demonstration.

**Service Endpoints:**
- **Web Interface:** http://localhost:8080
- **gRPC API:** localhost:50051 (default)
- **REST API Gateway:** http://localhost:8080/api
- **Kafka Broker:** localhost:9092
- **Redis Cache:** localhost:6378
- **PostgreSQL:** localhost:5431


###  Testing

Run the test suite using the following commands:

**Run All Tests:**
```sh
❯ go test ./...
```

**Run Tests for Specific Service:**
```sh
❯ cd api-gateway && go test ./...
❯ cd producer && go test ./...
❯ cd seeder && go test ./...
```

**Generate Test Coverage:**
```sh
❯ go test -coverprofile=coverage.out ./...
❯ go tool cover -html=coverage.out
```

**Integration Testing:**
1. Start the infrastructure services: `make up`
2. Run the seeder to generate test data: `cd seeder && go run cmd/app/main.go`
3. Start the producer to process data: `cd producer && go run cmd/app/main.go`
4. Start the API gateway: `cd api-gateway && go run cmd/app/main.go`
5. Open http://localhost:8080 to see real-time location updates on the map


---

##  API Documentation

### gRPC Service

The system provides a gRPC service for high-performance device communication:

**Service Definition:**
```protobuf
service LocationService {
  rpc StartSession(DeviceID) returns (google.protobuf.Empty);
  rpc SendLocation(OBUData) returns (google.protobuf.Empty);
  rpc GetDeviceLocations(google.protobuf.Empty) returns (stream OBUData);
}
```

**Message Types:**
```protobuf
message OBUData {
  bytes device_id = 1;       // 16-byte UUID
  double latitude = 2;       // GPS latitude
  double longitude = 3;      // GPS longitude  
  google.protobuf.Timestamp timestamp = 4;
}

message DeviceID {
  bytes device_id = 1;       // 16-byte UUID
}
```

### REST API Gateway

The HTTP gateway provides RESTful access to the gRPC services:

- **POST /api/v1/devices/{device_id}/session** - Start device session
- **POST /api/v1/locations** - Send location data
- **GET /api/v1/locations/stream** - Stream location updates
- **GET /api/v1/health** - Health check endpoint

### WebSocket API

Real-time location updates via WebSocket:

- **Endpoint:** `ws://localhost:8080/ws`
- **Message Format:** JSON-encoded OBUData
- **Connection:** Automatic reconnection with exponential backoff

---

##  Project Roadmap

- [X] **`Task 1`**: <strike>Implement Kafka producer and data seeder.</strike>
- [X] **`Task 2`**: <strike>Implement Kafka consumer for location data processing.</strike>
- [X] **`Task 3`**: <strike>Basic WebSocket implementation and HTML templates.</strike>
- [X] **`Task 4`**: <strike>Integration with Leaflet.js for map visualization.</strike>
- [ ] **`Task 5`**: Add Grafana dashboards for analytics and monitoring.
- [ ] **`Task 6`**: Implement device authentication and authorization.
- [ ] **`Task 7`**: Add geofencing and location-based alerts.
- [ ] **`Task 8`**: Performance optimization and horizontal scaling.
- [ ] **`Task 9`**: Comprehensive test coverage and CI/CD pipeline.
---

##  Contributing

- **💬 [Join the Discussions](https://github.com/rauan06/realtime-map/discussions)**: Share your insights, provide feedback, or ask questions.
- **🐛 [Report Issues](https://github.com/rauan06/realtime-map/issues)**: Submit bugs found or log feature requests for the `realtime-map` project.
- **💡 [Submit Pull Requests](https://github.com/rauan06/realtime-map/blob/main/CONTRIBUTING.md)**: Review open PRs, and submit your own PRs.

<details closed>
<summary>Contributing Guidelines</summary>

1. **Fork the Repository**: Start by forking the project repository to your github account.
2. **Clone Locally**: Clone the forked repository to your local machine using a git client.
   ```sh
   git clone https://github.com/rauan06/realtime-map
   ```
3. **Create a New Branch**: Always work on a new branch, giving it a descriptive name.
   ```sh
   git checkout -b new-feature-x
   ```
4. **Make Your Changes**: Develop and test your changes locally.
5. **Commit Your Changes**: Commit with a clear message describing your updates.
   ```sh
   git commit -m 'Implemented new feature x.'
   ```
6. **Push to github**: Push the changes to your forked repository.
   ```sh
   git push origin new-feature-x
   ```
7. **Submit a Pull Request**: Create a PR against the original project repository. Clearly describe the changes and their motivations.
8. **Review**: Once your PR is reviewed and approved, it will be merged into the main branch. Congratulations on your contribution!
</details>

<details closed>
<summary>Contributor Graph</summary>
<br>
<p align="left">
   <a href="https://github.com/rauan06/realtime-map/graphs/contributors">
      <img src="https://contrib.rocks/image?repo=rauan06/realtime-map">
   </a>
</p>
</details>

---

##  License

This project is protected under the MIT License. For more details, refer to the [LICENSE](https://github.com/rauan06/realtime-map/blob/main/LICENSE) file.

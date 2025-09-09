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
- [ Getting Started](#-getting-started)
  - [ Prerequisites](#-prerequisites)
  - [ Installation](#-installation)
  - [ Usage](#-usage)
  - [ Testing](#-testing)
- [ Project Roadmap](#-project-roadmap)
- [ Contributing](#-contributing)
- [ License](#-license)

---

##  Overview

<code>❯ REPLACE-ME</code>

---

##  Features

<code>❯ REPLACE-ME</code>

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

##  Getting Started

###  Prerequisites

Before getting started with realtime-map, ensure your runtime environment meets the following requirements:

- **Programming Language:** Go version 1.24.6
- **Package Manager:** Go modules
- **Container Runtime:** Docker


###  Installation

Install realtime-map using one of the following methods:

**Build from source:**

1. Clone the realtime-map repository:
```sh
❯ git clone https://github.com/rauan06/realtime-map
```

2. Navigate to the project directory:
```sh
❯ cd realtime-map
```

3. Install the project dependencies:


**Using `go modules`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Go-00ADD8.svg?style={badge_style}&logo=go&logoColor=white" />](https://golang.org/)

```sh
❯ go build
```


**Using `docker`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Docker-2CA5E0.svg?style={badge_style}&logo=docker&logoColor=white" />](https://www.docker.com/)

```sh
❯ docker build -t rauan06/realtime-map .
```




###  Usage
Run realtime-map using the following command:
**Using `go modules`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Go-00ADD8.svg?style={badge_style}&logo=go&logoColor=white" />](https://golang.org/)

```sh
❯ go run {entrypoint}
```


**Using `docker`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Docker-2CA5E0.svg?style={badge_style}&logo=docker&logoColor=white" />](https://www.docker.com/)

```sh
❯ docker run -it {image_name}
```


###  Testing
Run the test suite using the following command:
**Using `go modules`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Go-00ADD8.svg?style={badge_style}&logo=go&logoColor=white" />](https://golang.org/)

```sh
❯ go test ./...
```


---
##  Project Roadmap

- [X] **`Task 1`**: <strike>Implement kafka producer and data seeder.</strike>
- [ ] **`Task 2`**: Implement kafka consumer.
- [ ] **`Task 3`**: Integrate WebSocket protocol and Grafana.
- [ ] **`Task 4`**: Write HTML templates, integration with Leaflet.js.
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
   <a href="https://github.com{/rauan06/realtime-map/}graphs/contributors">
      <img src="https://contrib.rocks/image?repo=rauan06/realtime-map">
   </a>
</p>
</details>

---

##  License

This project is protected under the MIT License. For more details, refer to the [LICENSE](https://choosealicense.com/licenses/) file.

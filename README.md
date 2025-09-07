<p align="center">
    <img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/ec559a9f6bfd399b82bb44393651661b08aaf7ba/icons/folder-markdown-open.svg" align="center" width="30%">
</p>
<p align="center"><h1 align="center">REALTIME-MAP</h1></p>
<p align="center">
	<em><code>❯ REPLACE-ME</code></em>
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
  - [ Project Index](#-project-index)
- [ Getting Started](#-getting-started)
  - [ Prerequisites](#-prerequisites)
  - [ Installation](#-installation)
  - [ Usage](#-usage)
  - [ Testing](#-testing)
- [ Project Roadmap](#-project-roadmap)
- [ Contributing](#-contributing)
- [ License](#-license)
- [ Acknowledgments](#-acknowledgments)

---

##  Overview

<code>❯ REPLACE-ME</code>

---

##  Features

<code>❯ REPLACE-ME</code>

---

##  Project Structure

```sh
└── realtime-map/
    ├── Dockerfile
    ├── Makefile
    ├── api-gateway
    │   ├── cmd
    │   ├── config
    │   ├── go.mod
    │   ├── go.sum
    │   ├── internal
    │   └── template
    ├── docker-compose.yml
    ├── go-commons
    │   ├── buf.gen.yaml
    │   ├── buf.lock
    │   ├── buf.yaml
    │   ├── gen
    │   ├── go.mod
    │   ├── go.sum
    │   ├── pkg
    │   └── proto
    ├── producer
    │   ├── README.md
    │   ├── cmd
    │   ├── config
    │   ├── go.mod
    │   ├── go.sum
    │   └── internal
    └── seeder
        ├── cmd
        ├── go.mod
        ├── go.sum
        └── internal
```


###  Project Index
<details open>
	<summary><b><code>REALTIME-MAP/</code></b></summary>
	<details> <!-- __root__ Submodule -->
		<summary><b>__root__</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/Makefile'>Makefile</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/docker-compose.yml'>docker-compose.yml</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/Dockerfile'>Dockerfile</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			</table>
		</blockquote>
	</details>
	<details> <!-- seeder Submodule -->
		<summary><b>seeder</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/seeder/go.mod'>go.mod</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/seeder/go.sum'>go.sum</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			</table>
			<details>
				<summary><b>cmd</b></summary>
				<blockquote>
					<details>
						<summary><b>app</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/seeder/cmd/app/main.go'>main.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
				</blockquote>
			</details>
			<details>
				<summary><b>internal</b></summary>
				<blockquote>
					<details>
						<summary><b>repo</b></summary>
						<blockquote>
							<details>
								<summary><b>grpcclient</b></summary>
								<blockquote>
									<table>
									<tr>
										<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/seeder/internal/repo/grpcclient/client.go'>client.go</a></b></td>
										<td><code>❯ REPLACE-ME</code></td>
									</tr>
									</table>
								</blockquote>
							</details>
						</blockquote>
					</details>
					<details>
						<summary><b>domain</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/seeder/internal/domain/models.go'>models.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
				</blockquote>
			</details>
		</blockquote>
	</details>
	<details> <!-- producer Submodule -->
		<summary><b>producer</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/go.mod'>go.mod</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/go.sum'>go.sum</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			</table>
			<details>
				<summary><b>config</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/config/config.go'>config.go</a></b></td>
						<td><code>❯ REPLACE-ME</code></td>
					</tr>
					</table>
				</blockquote>
			</details>
			<details>
				<summary><b>cmd</b></summary>
				<blockquote>
					<details>
						<summary><b>app</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/cmd/app/main.go'>main.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
				</blockquote>
			</details>
			<details>
				<summary><b>internal</b></summary>
				<blockquote>
					<details>
						<summary><b>usecase</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/internal/usecase/contracts.go'>contracts.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
							<details>
								<summary><b>producer</b></summary>
								<blockquote>
									<table>
									<tr>
										<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/internal/usecase/producer/producer.go'>producer.go</a></b></td>
										<td><code>❯ REPLACE-ME</code></td>
									</tr>
									</table>
								</blockquote>
							</details>
						</blockquote>
					</details>
					<details>
						<summary><b>repo</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/internal/repo/contracts.go'>contracts.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
							<details>
								<summary><b>eventbus</b></summary>
								<blockquote>
									<table>
									<tr>
										<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/internal/repo/eventbus/kafka.go'>kafka.go</a></b></td>
										<td><code>❯ REPLACE-ME</code></td>
									</tr>
									</table>
								</blockquote>
							</details>
							<details>
								<summary><b>cache</b></summary>
								<blockquote>
									<table>
									<tr>
										<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/internal/repo/cache/redis.go'>redis.go</a></b></td>
										<td><code>❯ REPLACE-ME</code></td>
									</tr>
									</table>
								</blockquote>
							</details>
						</blockquote>
					</details>
					<details>
						<summary><b>domain</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/internal/domain/obu.go'>obu.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/internal/domain/kafka_msg.go'>kafka_msg.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
					<details>
						<summary><b>controller</b></summary>
						<blockquote>
							<details>
								<summary><b>grpcrouter</b></summary>
								<blockquote>
									<table>
									<tr>
										<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/internal/controller/grpcrouter/router.go'>router.go</a></b></td>
										<td><code>❯ REPLACE-ME</code></td>
									</tr>
									</table>
									<details>
										<summary><b>v1</b></summary>
										<blockquote>
											<table>
											<tr>
												<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/internal/controller/grpcrouter/v1/controller.go'>controller.go</a></b></td>
												<td><code>❯ REPLACE-ME</code></td>
											</tr>
											<tr>
												<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/internal/controller/grpcrouter/v1/router.go'>router.go</a></b></td>
												<td><code>❯ REPLACE-ME</code></td>
											</tr>
											</table>
										</blockquote>
									</details>
								</blockquote>
							</details>
						</blockquote>
					</details>
					<details>
						<summary><b>app</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/producer/internal/app/app.go'>app.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
				</blockquote>
			</details>
		</blockquote>
	</details>
	<details> <!-- api-gateway Submodule -->
		<summary><b>api-gateway</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/api-gateway/go.mod'>go.mod</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/api-gateway/go.sum'>go.sum</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			</table>
			<details>
				<summary><b>config</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/api-gateway/config/config.go'>config.go</a></b></td>
						<td><code>❯ REPLACE-ME</code></td>
					</tr>
					</table>
				</blockquote>
			</details>
			<details>
				<summary><b>cmd</b></summary>
				<blockquote>
					<details>
						<summary><b>app</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/api-gateway/cmd/app/main.go'>main.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
				</blockquote>
			</details>
			<details>
				<summary><b>template</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/api-gateway/template/map.html'>map.html</a></b></td>
						<td><code>❯ REPLACE-ME</code></td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/api-gateway/template/index.html'>index.html</a></b></td>
						<td><code>❯ REPLACE-ME</code></td>
					</tr>
					</table>
				</blockquote>
			</details>
			<details>
				<summary><b>internal</b></summary>
				<blockquote>
					<details>
						<summary><b>controller</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/api-gateway/internal/controller/server.go'>server.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
					<details>
						<summary><b>app</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/api-gateway/internal/app/app.go'>app.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
				</blockquote>
			</details>
		</blockquote>
	</details>
	<details> <!-- go-commons Submodule -->
		<summary><b>go-commons</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/go.mod'>go.mod</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/buf.yaml'>buf.yaml</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/go.sum'>go.sum</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/buf.gen.yaml'>buf.gen.yaml</a></b></td>
				<td><code>❯ REPLACE-ME</code></td>
			</tr>
			</table>
			<details>
				<summary><b>proto</b></summary>
				<blockquote>
					<details>
						<summary><b>route</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/proto/route/route.proto'>route.proto</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
				</blockquote>
			</details>
			<details>
				<summary><b>gen</b></summary>
				<blockquote>
					<details>
						<summary><b>proto</b></summary>
						<blockquote>
							<details>
								<summary><b>route</b></summary>
								<blockquote>
									<table>
									<tr>
										<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/gen/proto/route/route.pb.go'>route.pb.go</a></b></td>
										<td><code>❯ REPLACE-ME</code></td>
									</tr>
									<tr>
										<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/gen/proto/route/route.pb.gw.go'>route.pb.gw.go</a></b></td>
										<td><code>❯ REPLACE-ME</code></td>
									</tr>
									<tr>
										<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/gen/proto/route/route_grpc.pb.go'>route_grpc.pb.go</a></b></td>
										<td><code>❯ REPLACE-ME</code></td>
									</tr>
									</table>
								</blockquote>
							</details>
						</blockquote>
					</details>
					<details>
						<summary><b>openapiv2</b></summary>
						<blockquote>
							<details>
								<summary><b>proto</b></summary>
								<blockquote>
									<details>
										<summary><b>route</b></summary>
										<blockquote>
											<table>
											<tr>
												<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/gen/openapiv2/proto/route/route.swagger.json'>route.swagger.json</a></b></td>
												<td><code>❯ REPLACE-ME</code></td>
											</tr>
											</table>
										</blockquote>
									</details>
								</blockquote>
							</details>
						</blockquote>
					</details>
				</blockquote>
			</details>
			<details>
				<summary><b>pkg</b></summary>
				<blockquote>
					<details>
						<summary><b>logger</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/pkg/logger/logger.go'>logger.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
					<details>
						<summary><b>postgres</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/pkg/postgres/options.go'>options.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/pkg/postgres/postgres.go'>postgres.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
					<details>
						<summary><b>grpcserver</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/pkg/grpcserver/options.go'>options.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/rauan06/realtime-map/blob/master/go-commons/pkg/grpcserver/server.go'>server.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
				</blockquote>
			</details>
		</blockquote>
	</details>
</details>

---
##  Getting Started

###  Prerequisites

Before getting started with realtime-map, ensure your runtime environment meets the following requirements:

- **Programming Language:** Go
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

- [X] **`Task 1`**: <strike>Implement feature one.</strike>
- [ ] **`Task 2`**: Implement feature two.
- [ ] **`Task 3`**: Implement feature three.

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

This project is protected under the [SELECT-A-LICENSE](https://choosealicense.com/licenses) License. For more details, refer to the [LICENSE](https://choosealicense.com/licenses/) file.

---

##  Acknowledgments

- List any resources, contributors, inspiration, etc. here.

---
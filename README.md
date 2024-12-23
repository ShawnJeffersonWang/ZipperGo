# Campus Navigation Project

This project is a campus navigation system implemented in Go, designed as a data structures course project. It uses the Gin framework and Redis, with a focus on graph usage and the implementation of Dijkstra's algorithm for shortest path finding.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Setup](#setup)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Introduction

The Campus Navigation Project is designed to help students and visitors navigate a campus efficiently. The system utilizes graph data structures to represent the campus layout and implements Dijkstra's algorithm to find the shortest path between locations.

## Features

- **Graph Representation:** Models the campus using graph data structures.
- **Shortest Path Calculation:** Implements Dijkstra's algorithm to find the shortest path between two points.
- **Redis Integration:** Uses Redis for caching and fast data retrieval.
- **RESTful API:** Provides a RESTful API to interact with the navigation system.

## Technology Stack

- **Backend:** Go, Gin framework
- **Database:** Redis
- **Algorithm:** Dijkstra's algorithm for shortest path

## Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/ShawnJeffersonWang/campus-guide.git
    cd campus-guide
    ```

2. Ensure you have Go and Redis installed on your system.

3. Install the necessary Go dependencies:
    ```sh
    go mod tidy
    ```

4. Start Redis server if it is not already running:
    ```sh
    redis-server
    ```

5. Run the application:
    ```sh
    go run main.go
    ```

## Usage

Once the application is running, you can access the API endpoints to interact with the navigation system.

### Example Endpoints

- **Get Shortest Path:** `/api/shortest-path?start={start}&end={end}`

## Contributing

Contributions are welcome! Please fork this repository and submit pull requests with your changes.

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -am 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Create a new Pull Request.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

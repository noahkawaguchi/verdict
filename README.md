# Verdict

### [🔗 Live interactive docs here](https://verdict.noahkawaguchi.com)

---

Verdict is a serverless REST API (with interactive Swagger UI) for working with ranked choice voting. Consumers can create polls, cast ballots, and retrieve results.

Instead of selecting a single choice, voters rank all choices in order of preference. The instant runoff algorithm then calculates a winner by repeatedly eliminating the last-place choice, using a recursive tie-breaking sub-poll when needed, and redistributing votes until one choice achieves a strict majority.

## Table of Contents

1. [Tech Stack](#tech-stack)
2. [Architecture](#architecture)
3. [Project Structure](#project-structure)
4. [Local Development](#local-development)
5. [Testing](#testing)
6. [About Ranked Choice Voting](#about-ranked-choice-voting)

## Tech Stack

|                |                                                                                                                                                                                                                                                                           |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Language       | ![Static Badge](https://img.shields.io/badge/Go-00ADD8)                                                                                                                                                                                                                   |
| Infrastructure | ![Static Badge](https://img.shields.io/badge/Amazon_API_Gateway-FF4F8B) ![Static Badge](https://img.shields.io/badge/AWS_Lambda-FF9900) ![Static Badge](https://img.shields.io/badge/Amazon_DynamoDB-4053D6) ![Static Badge](https://img.shields.io/badge/AWS_SAM-232F3E) |
| Docs           | ![Static Badge](https://img.shields.io/badge/Swagger_UI-85EA2D) ![Static Badge](https://img.shields.io/badge/GitHub_Pages-222222)                                                                                                                                         |

## Architecture

The project is fully serverless, with API Gateway routing requests to a single Go Lambda function that reads from and writes to DynamoDB. Since there's no persistent server to manage, the cost scales to zero at rest.

Key design decisions:

- **Go**: Go compiles to a small, self-contained binary with low memory usage and fast cold starts, making it a natural fit for Lambda.
- **Single Lambda function**: All routes are handled by one function with an internal router, avoiding the overhead of managing multiple function deployments for a small API.
- **DynamoDB**: A simple key-value access pattern (get poll by ID, get all ballots for a poll) maps cleanly to DynamoDB without needing a relational schema.

## Project Structure

The design of the codebase generally follows the common three-layer pattern of API, domain, and infrastructure, but it allows for some pragmatic compromises due to the project's relatively small scale.

```
verdict/
├── cmd/              # Lambda entrypoint and dev/prod config initialization
├── internal/
│   ├── api/          # HTTP routing and request handlers
│   ├── voting/       # Domain types and their associated logic
│   └── datastore/    # DynamoDB read/write
├── docs/             # Swagger UI static site (served via GitHub Pages)
├── template.yaml     # AWS SAM/CloudFormation infrastructure definition
└── justfile          # Project-specific commands and config (similar to Makefile)
```

## Local Development

### Prerequisite Installations

- **[Go](https://go.dev/doc/install)**: The main language of the project.
- **[Python](https://www.python.org/downloads/)**: Only used to serve the static docs locally, but likely already installed on most systems.
- **[Docker](https://docs.docker.com/get-started/get-docker/)**: For the local Lambda/API Gateway emulator and DynamoDB container.
- **[AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)**: For local emulation of SAM infrastructure and deployment.
- **[Just](https://github.com/casey/just#installation)**: Command runner used to save and run the commands used in development.

### Running the Project Locally

1. With the Docker daemon running, execute the following commands:

```sh
just       # Starts the Docker network, DynamoDB container, and API (port 3000)
just docs  # In a separate terminal, serves the Swagger UI docs (port 8000)
```

2. View the interactive docs in a browser at [http://localhost:8000](http://localhost:8000).
3. Use the server dropdown to select the `localhost` option.
4. Use the "try it out" functionality to hit the local endpoints.

All available recipes are documented in the [`justfile`](justfile). (Run `just --list` to see a summary.)

## Testing

```sh
just test     # Runs all tests
just test -v  # Runs all tests with verbose output
```

Tests cover all three packages (`api`, `voting`, `datastore`) and focus on the following strategies:

- **Dependency inversion and dependency injection**: The `api` and `datastore` packages depend on interfaces rather than concrete implementations, making it straightforward to substitute mock implementations in tests.
- **Flexible mocking**: Hand-rolled mocks use function fields with `nil` fallbacks, allowing each test case to inject only the behavior relevant to it without having to implement the full interface every time.
- **Table-driven tests**: Tests are organized using Go's table-driven pattern to efficiently run multiple test cases through the same logic.
- **External test packages**: All test files use external test packages (e.g. `package api_test`), testing public APIs as consumers use them rather than having access to internal implementation details.

## About Ranked Choice Voting

Verdict implements a voting algorithm known as ranked choice voting or instant runoff voting. There are several variations, but this is the one used here:

### Voting process

Instead of selecting a single choice, voters rank all choices in order of preference.

### Determining a winner

Instead of immediately selecting a choice that has only a plurality of votes, the algorithm first checks if any choice has a strict majority of votes. If no choice does, the choice with the fewest votes is eliminated, and its votes are redistributed to the voters' next highest choices. This process of elimination continues until a single choice has a strict majority of votes.

### What about ties for last?

In each round, the choice with the fewest votes is eliminated, but what if multiple choices are tied for last place? In this case, a sub-poll is simulated between only the tied choices. This is possible because voters provide a rank for every choice, allowing the algorithm to determine their preferences amongst any subset of choices. If there is another tie for last place, the tie-breaking algorithm continues recursively.

While unlikely unless the numbers of choices and voters are very small, it is possible that multiple choices are tied for last place and received perfectly equivalent rankings. In this case, one of these lowest-ranking choices is eliminated by pseudorandom number generation.

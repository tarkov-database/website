![Lint & Test](https://github.com/tarkov-database/website/workflows/Lint%20&%20Test/badge.svg)
[![Lines of code](https://tokei.rs/b1/github/tarkov-database/website)](https://github.com/XAMPPRocky/tokei)

# Tarkov Database Website

The main presence, the website of [Tarkov Database](https://tarkov-database.com) which presents the data of the [API](https://github.com/tarkov-database/rest-api) in a friendly way.

The project follows the [KISS principle](https://en.wikipedia.org/wiki/KISS_principle) and focuses on the same things as the REST API: scalability, efficiency, performance, and simplicity.
It is also designed as a stateless service for use in a distributed system.

## Set up a local development environment

### Prerequisites

- Git
- [Go](https://golang.org/doc/install) (latest)
- [NodeJS](https://nodejs.org/en/) (>=v15)
- [Podman](https://podman.io/getting-started/installation) (>=v2)
- Linux environment

### Set up

#### 1. [Set up local REST API](https://github.com/tarkov-database/resources/tree/master/rest-api#rest-api)

#### 2. Clone this repository

```BASH
git clone git@github.com:tarkov-database/website.git && cd website
```

#### 3. Install NPM modules

```BASH
npm i
```

#### 4. Build static files

```BASH
make statics
```

#### 5. Create environment variable file

Create a file with the name `.env` in the root directory and the following content (add the appropriate variables)

```SH
export HOST=localhost:8080
export API_URL=http://localhost:9000/v2
export API_TOKEN=<TOKEN>
export SEARCH_HOST=http://localhost:9100
```

#### 6. Start the server

```BASH
make run
```

The server should now be accessible via [localhost:8080](http://localhost:8080)

#### 7. Install Revive linter (optional)

```BASH
go install github.com/mgechev/revive@v1
```

#### 8. [Install Prettier editor integration](https://prettier.io/docs/en/editors.html) (optional)

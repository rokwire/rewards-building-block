# Rewards building block
The Rewards Building Block handles the management of reward inventory and transactions for the Rokwire platform.

## Documentation

The functionality provided by this application is documented in the [Wiki](https://github.com/rokwire/rewards-building-block/wiki).

The API documentation is available here: https://api.rokwire.illinois.edu/rewards/doc/ui/

## Set Up

### Prerequisites

MongoDB v4.2.2+

Go v1.16+

### Environment variables
The following Environment variables are supported. The service will not start unless those marked as Required are supplied.

Name|Format|Required|Description
---|---|---|---
PORT | < int  > | yes | Port to be used by this application
HOST | < string > | yes | URL where this application is being hosted
MONGO_AUTH | <mongodb://USER:PASSWORD@HOST:PORT/DATABASE NAME> | yes | MongoDB authentication string. The user must have read/write privileges.
MONGO_DATABASE | < string > | yes | MongoDB database name
MONGO_TIMEOUT | < int > | no |  MongoDB timeout in milliseconds. Defaults to 500
DEFAULT_CACHE_EXPIRATION_SECONDS | < int > | false | Default cache expiration time in seconds. Defaults to 120
CORE_BB_HOST | < string > | yes | Core BB host URL
REWARDS_SERVICE_URL | < string > | yes | Rewards base URL
INTERNAL_API_KEY | < string > | yes | Internal API key for the corresponding environment.

### Run Application

#### Run locally without Docker

1. Clone the repo (outside GOPATH)

2. Open the terminal and go to the root folder
  
3. Make the project  
```
$ make
...
▶ building executable(s)… 1.9.0 2020-08-13T10:00:00+0300
```

4. Run the executable
```
$ ./bin/content
```

#### Run locally as Docker container

1. Clone the repo (outside GOPATH)

2. Open the terminal and go to the root folder
  
3. Create Docker image  
```
docker build -t content .
```
4. Run as Docker container
```
docker-compose up
```

#### Tools

##### Run tests
```
$ make tests
```

##### Run code coverage tests
```
$ make cover
```

##### Run golint
```
$ make lint
```

##### Run gofmt to check formatting on all source files
```
$ make checkfmt
```

##### Run gofmt to fix formatting on all source files
```
$ make fixfmt
```

##### Cleanup everything
```
$ make clean
```

##### Run help
```
$ make help
```

##### Generate Swagger docs
```
$ make swagger
```

### Test Application APIs

Verify the service is running as calling the get version API.

#### Call get version API

curl -X GET -i http://localhost/rewards/version

Response
```
1.9.0
```

## Contributing
If you would like to contribute to this project, please be sure to read the [Contributing Guidelines](CONTRIBUTING.md), [Code of Conduct](CODE_OF_CONDUCT.md), and [Conventions](CONVENTIONS.md) before beginning.

### Secret Detection
This repository is configured with a [pre-commit](https://pre-commit.com/) hook that runs [Yelp's Detect Secrets](https://github.com/Yelp/detect-secrets). If you intend to contribute directly to this repository, you must install pre-commit on your local machine to ensure that no secrets are pushed accidentally.

```
# Install software 
$ git pull  # Pull in pre-commit configuration & baseline 
$ pip install pre-commit 
$ pre-commit install
```
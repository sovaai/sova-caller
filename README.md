[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)

# SOVA Caller App

## Installation

### Requirements

- Go 1.18.3 (https://go.dev/doc/install)
- Node v14.18 (https://nodejs.org/en/download/)
- npm latest (https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
- PJSIP v.2.12 ([Download PJSIP](https://www.pjsip.org/download.htm) and [Install](https://trac.pjsip.org/repos/wiki/Getting-Started)). Install missing libraries if necessary (libssl-dev, uuid-dev, libasound2-dev)

### Installing dependencies

To install project dependencies in the sova-caller root directory, run the following commands:

```
npm i
cd ./ui
npm i
cd ./server
npm i
```

To run the project in development mode in the root directory, run:

```
npm run start
```

The backend part will launch locally on port 4000 [http://localhost:4000](http://localhost:4000)
The UI part will be available in the browser at [http://localhost:3000](http://localhost:3000)

## Building the image and running the container

### Requirements

- Node v14.18 (https://nodejs.org/en/download/)
- npm latest (https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
- Docker latest (https://www.docker.com/)

The image builds both parts of the project, it contains the necessary dependencies (Go and PJSIP), so you only need to run the following commands to run the application from the container.

To install project dependencies in the root directory, run the commands:

```
npm i
```

To build an image from a Dockerfile, run:

```
docker build -f ./docker/Dockerfile -t sova-caller .
```

To start a container and give it a name my-sova-caller run the command:

```
docker run --name my-sova-caller -p 4000:4000 sova-caller
```

The application is raised locally on port 4000 [http://localhost:4000](http://localhost:4000)

You can run the following commands to stop and start a container:

```
docker stop my-sova-caller
docker start my-sova-caller
```
## Licenses

SOVA Caller is licensed under Apache License 2.0 by Virtual Assistant, LLC.

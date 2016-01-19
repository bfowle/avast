# avast
...ye maties!
...ya-harrrr!

![Avast](https://media.giphy.com/media/E8KFBhPh2s3ra/giphy.gif)

# Server-side

### Building

```
get get ./src/
# recommended:
# gb vendor fetch --all

go build -o bin/avast ./src/
```

### Running

```
AVAST_API_VERSION=v1 \
  AVAST_ADDR=:8080 \
  DOCKER_HOST=tcp://1.1.1.1:2375 \
  DOCKER_API_VERSION=v1.21 \
  CONSUL_HTTP_ADDR=1.1.1.1:8500 \
  bin/avast
```

# Client-side

> Located in the `client/` folder

### Dependencies

- `node` (v4.1.x+) and `npm` (2.14.x+)
- (Development) global install of the following:
```
npm i --global webpack \
  webpack-dev-server \
  karma \
  protractor \
  typings \
  typescript
```

### Installing

```
npm i
# `typings install` will be run in the postinstall hook
```

### Running

Starting the development server

```
npm run server
# http://0.0.0.0:3000/
# or IPv6 http://[::1]:3000
```

Watch and build files

```
npm run watch
```

Running tests

````
npm run test
```

Starting the e2e webdriver

```
npm run webdriver:update
npm run webdriver:start
```

Running e2e tests

```
npm run e2e
```

### Building

```
# development
npm run build:dev

# production
npm run build:prod
```

### Generating Documentation

```
npm run docs
# located at `client/docs/index.html`
```

# TODO

- Split repo into `avast-server` and `avast-ui` repos

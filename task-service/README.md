Docker is a prerequisite for running this application.

### How to run

#### Using makefile
Build
```shell
make build-backend
```

Run
```shell
make start-backend
```

Stop and clean database (even after pressing Ctrl+C)
```shell
make stop-all
```

#### Using docker-compose
Just use commands from `Makefile` directly

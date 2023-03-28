Docker is a prerequisite for running this application.

### How to run

#### Using makefile
Build
```bash
make build-backend
```

Run
```bash
make start-backend
```

Stop and clean database (even after pressing Ctrl+C)
```bash
make stop-all
```

#### Using docker-compose
Just use commands from `Makefile` directly

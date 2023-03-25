Docker is a prerequisite for running this application.

### How to run

#### Using makefile
Build
```bash
make rebuild-backend
```

Run
```bash
make start-backend
```

Stop and clean database (even after pressing Ctrl+C)
```bash
make rebuild-backend
```

#### Using docker-compose
Just use commands from `Makefile`
# K6 load test

## How to run

#### Using makefile
Optionally to see logs, start up the `task` service
```shell
make start-backend
```

Run exemplary load test or adjust the parameters in `Makefile` and run it again.
```shell
make run-k6
```

Access the grafana dashboard at `http://localhost:3000`

To stop all running services
```shell
make stop-all
```

#### Using docker-compose
Just use commands from `Makefile` directly

## Improved data generation
There is more sophisticated way to generate data using `faker.js`. 
Setup process is described in [this article](https://dev.to/k6/performance-testing-with-generated-data-using-k6-and-faker-2e).

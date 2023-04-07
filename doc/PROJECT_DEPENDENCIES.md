# Project Dependencies

- [CockroachDB](https://www.cockroachlabs.com): the project uses  this database to store information. in project, functions works with database indirectly. 
- [pgadmin](https://www.pgadmin.org): to monitor CockroachDB
- [redpanda](https://redpanda.com): redpanda is a queue service. requests and responses stored at redpanda's storage
- [otelcol](https://opentelemetry.io/docs/collector/): opentelemetry collector. we use opentelemetry for tracing this project. tracing information sent to opentelemetry collector.
- [jaeger](https://www.jaegertracing.io): jaeger is tracing monitoring system that gets tracing information from opentelemetry collector
- [swagger](https://swagger.io): used for testing api's

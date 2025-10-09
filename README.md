# pathfinder

Data Platform project

## Get Started

### Docker Compose

- Create a .env file in the deploy/docker directory

```shell
# InfluxDB
SESSION_SECRET_KEY=<Generated with openssl rand -hex 32>
INFLUXDB_ADMIN_TOKEN=<Generated token with docker exec -it influxdb influxdb3 create token --admin>
INFLUXDB_WRITE_TOKEN=<Generated token with docker exec -it influxdb influxdb3 create token>

# RabbitMQ
RABBITMQ_ADMIN_USER=<Rabbitmq root user>
RABBITMQ_ADMIN_PASSWORD=<RabbitMQ root password>

# Postgres
POSTGRES_ADMIN_PASSWORD=<postgress admin password>
```

- Verify the haproxy domain names and include them in your hosts file if you want to address the different services via domain names

- Execute docker compose up -d

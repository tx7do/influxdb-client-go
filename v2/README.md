# influxdb-client-go v2

## Docker Deploy InfluxDB Server

```bash
docker pull bitnami/influxdb:latest

docker run -itd \
    --name influxdb-test \
    -p 8083:8083 \
    -p 8086:8086 \
    -e INFLUXDB_HTTP_AUTH_ENABLED=true \
    -e INFLUXDB_ADMIN_USER=admin \
    -e INFLUXDB_ADMIN_USER_PASSWORD=123456789 \
    -e INFLUXDB_ADMIN_USER_TOKEN=admintoken123 \
    -e INFLUXDB_DB=my_database \
    bitnami/influxdb:latest
```

- Admin UI: <http://localhost:8086/>

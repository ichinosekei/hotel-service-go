
# README для Delivery Service и Notification Service

## Структура сервисов

### Delivery Service
```
delivery-service/
    delivery-service/
        GlobalSign.pem
        go.mod
        Dockerfile
        go.sum
        README.md
        .env.dev
        cmd/
            delivery_service/
                main.go
        tests/
            integration/
                grpc_test.go
        internal/
            grpc/
                handler.go
                proto/
                    delivery_grpc.pb.go
                    delivery.proto
                    delivery.pb.go
            sms/
                sms.go
        configs/
            config.go
        deployments/
            docker/
                docker-compose.yaml
```

### Notification Service
```
notification-service/
    notification-service/
        go.mod
        Dockerfile
        go.sum
        README.md
        .env.dev
        cmd/
            notification_service/
                main.go
        tests/
            integration/
                kafka_test.go
                publish_test.go
                grpc_test.go
                consume_test.go
        internal/
            grpc/
                client.go
                proto/
                    delivery_grpc.pb.go
                    delivery.proto
                    delivery.pb.go
            kafka/
                producer.go
                deduplicator.go
                kafka.go
                consumer.go
        configs/
            config.go
        deployments/
            docker/
                docker-compose.yaml
                kafka-data/
                zookeeper-data/
```

---

## Инструкции по работе

### 1. Создание сети Docker
Перед началом работы создайте сеть для взаимодействия контейнеров:

```bash
docker network create app-network
```

### 2. Запуск Delivery Service

   ```bash
   cd delivery-service/deployments/docker
   docker-compose down
   docker-compose build --no-cache
   docker-compose up -d
   ```

### 3. Запуск Notification Service
   
   ```bash
   cd notification-service/deployments/docker
   docker-compose down
   docker-compose build --no-cache
   docker-compose up -d
   ```

### 4. Работа с Kafka

#### Просмотр списка топиков:
```bash
docker exec -it docker-kafka-1 kafka-topics --list --bootstrap-server kafka:9092
```

#### Создание топика:
```bash
docker exec -it docker-kafka-1 kafka-topics --create --topic booking-events --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
```

#### Заполнение топика из консоли:
```bash
docker exec -it docker-kafka-1 kafka-console-producer --bootstrap-server kafka:9092 --topic booking-events
{"id": "0", "client_phone":"89118026917","check_in":"2024-01-01","check_out":"2024-01-10","room":"101","hotel_phone":"89118026917"}
```

### 5. Работа с логами

#### Просмотр логов контейнера:
```bash
docker logs -f <container_name>
```

---

Просмотр запущенных и приостановленных контейнеров
```bash
docker ps --all
```

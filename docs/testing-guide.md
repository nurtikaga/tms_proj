# Testing Guide

## What is worth checking

In this service, Kafka is used as a producer only.

The main flow is:

1. `CreateShipment` creates a shipment in Postgres.
2. `AddEvent` changes shipment status.
3. `AddEvent` also publishes a JSON event to Kafka topic `shipment-events`.
4. `GetShipment` and `GetHistory` let you verify the saved state and timeline.

## Start the stack

From the project root:

```powershell
docker compose -f deployments/docker-compose.yml up --build
```

The app will be available on `localhost:50051`.

Dependencies from config:

- gRPC: `localhost:50051`
- Postgres: `localhost:5432`
- Redis: `localhost:6379`
- Kafka: `localhost:9092`
- Kafka topic: `shipment-events`

## How to test handlers in Postman

This project exposes gRPC, not REST.

In Postman:

1. Create a new `gRPC` request.
2. Server address: `localhost:50051`
3. Load proto file: `tms-protos/ShipmentManaging.proto`
4. Select service `tmsmanagerpb.v1.ShipmentService`
5. Run methods using the payloads from `docs/postman-grpc-samples.json`

## Happy path scenario

### 1. CreateShipment

Method: `ShipmentService/CreateShipment`

Request:

```json
{
  "reference_number": "REF-POSTMAN-001",
  "origin": "Almaty",
  "destination": "Astana",
  "driver_name": "Nurlan",
  "unit_number": "KZ-001",
  "shipment_amount": 1250000.5,
  "driver_revenue": 210000.25
}
```

Expected result:

- response contains non-empty `id`
- `status` equals `SHIPMENT_STATUS_PENDING`

Save returned `id`. It will be used as `shipment_id` below.

### 2. GetShipment

Method: `ShipmentService/GetShipment`

Request:

```json
{
  "id": "{{shipment_id}}"
}
```

Expected result:

- same shipment is returned
- current status is still `SHIPMENT_STATUS_PENDING`

### 3. AddEvent -> PICKED_UP

Method: `ShipmentService/AddEvent`

Request:

```json
{
  "shipment_id": "{{shipment_id}}",
  "status": "SHIPMENT_STATUS_PICKED_UP",
  "note": "Cargo picked up from shipper warehouse"
}
```

Expected result:

- gRPC success
- shipment status changes to `PICKED_UP`
- one Kafka message appears in topic `shipment-events`

### 4. AddEvent -> IN_TRANSIT

Method: `ShipmentService/AddEvent`

Request:

```json
{
  "shipment_id": "{{shipment_id}}",
  "status": "SHIPMENT_STATUS_IN_TRANSIT",
  "note": "Truck left origin city"
}
```

### 5. AddEvent -> DELIVERED

Method: `ShipmentService/AddEvent`

Request:

```json
{
  "shipment_id": "{{shipment_id}}",
  "status": "SHIPMENT_STATUS_DELIVERED",
  "note": "Shipment delivered to consignee"
}
```

### 6. GetHistory

Method: `ShipmentService/GetHistory`

Request:

```json
{
  "shipment_id": "{{shipment_id}}"
}
```

Expected result:

- 3 events in order:
  - `PICKED_UP`
  - `IN_TRANSIT`
  - `DELIVERED`

## Negative cases for handlers

### CreateShipment with missing required fields

Request:

```json
{
  "reference_number": "",
  "origin": "",
  "destination": ""
}
```

Expected result:

- gRPC code: `InvalidArgument`
- message mentions required fields

### GetShipment without id

Request:

```json
{
  "id": ""
}
```

Expected result:

- gRPC code: `InvalidArgument`

### AddEvent with invalid transition

This service allows only:

- `PENDING -> PICKED_UP`
- `PICKED_UP -> IN_TRANSIT`
- `IN_TRANSIT -> DELIVERED`

So this request must fail:

```json
{
  "shipment_id": "{{shipment_id}}",
  "status": "SHIPMENT_STATUS_DELIVERED",
  "note": "Invalid direct delivery"
}
```

Expected result:

- gRPC code: `FailedPrecondition`
- error text about invalid status transition

### AddEvent without shipment_id

Request:

```json
{
  "shipment_id": "",
  "status": "SHIPMENT_STATUS_PICKED_UP",
  "note": "Invalid request"
}
```

Expected result:

- gRPC code: `InvalidArgument`

## How to verify Kafka

Because the service publishes JSON to Kafka after `AddEvent`, the easiest check is to start a consumer and then call `AddEvent`.

If Kafka is running in Docker, open consumer output with:

```powershell
docker exec -it deployments-kafka-1 kafka-console-consumer --bootstrap-server kafka:9092 --topic shipment-events --from-beginning
```

If the container name differs, first inspect it with:

```powershell
docker ps
```

Expected Kafka message format:

```json
{
  "id": "event-uuid",
  "shipment_id": "shipment-uuid",
  "status": "PICKED_UP",
  "note": "Cargo picked up from shipper warehouse",
  "occurred_at": "2026-03-18T10:00:00Z"
}
```

Important detail:

- gRPC enum names are like `SHIPMENT_STATUS_PICKED_UP`
- Kafka payload status is plain domain value like `PICKED_UP`

## Quick verification checklist

- `CreateShipment` returns shipment with `PENDING`
- `GetShipment` returns created shipment
- valid `AddEvent` changes status
- invalid `AddEvent` returns `FailedPrecondition`
- `GetHistory` returns all applied events
- each successful `AddEvent` writes one JSON message to Kafka

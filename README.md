# Anymind BTC

Simple service to receive btc balance across timezones

## System Requirement

- Docker
- docker-compose CLI

## Installation

```bash
#cd to directory project
docker-compose up
```

Success installation may show this results:
```bash
Creating network "anymind_default" with the default driver
Creating mongodb ... done
Creating api     ... done
Attaching to mongodb, api
api        | time="2020-03-11T05:59:20Z" level=info msg="starting http application.."
```

## API

Local Environment: `http://localhost:9091`

List of available endpoint

- Add Balance 
- Get Balances Hourly Range

### A. Add Balance

Add your BTC balance across different timezone. PS: datatime request should not less than current hour

- Endpoint: `/v1/balance`
- Methods: POST

**Request**

```bash
curl --location --request POST 'http://localhost:9091/v1/balance' \
--header 'Content-Type: application/json' \
--data-raw '{
	"datetime": "2020-03-11T13:01:05+07:00",
	"amount": 10.1
}'
```

**Response**

- Success (200)
```json
{
    "data": "a26df952-f1eb-4ae9-801d-33a2f32f9a57",
    "message": "balance was added successfully"
}
```

- Bad request (400)
```json
{
    "data": null,
    "message": "request time should not older than last hour"
}
```

### B. Get Balances Hourly Range

Retrieve accumulated balance in hourly range

- Endpoint: `/v1/balance/hours/{start}/{end}`
- Methods: GET

PS: start & end should use RFC 3339 time format

**Request**

```bash
curl --location --request GET 'http://localhost:9091/v1/balance/hours/2020-03-11T12:10:05+07:00/2020-03-11T14:20:05+07:00'
```

**Response**

- Success (200):

```json
{
    "data": [
        {
            "datetime": "2020-03-11T06:00:05Z",
            "amount": "10.0"
        },
        {
            "datetime": "2020-03-11T07:00:05Z",
            "amount": "2010.0"
        }
    ],
    "message": "2 record(s) retrieved"
}
```

- Bad Request (400)

```json
{
    "data": null,
    "message": "wrong end datetime request"
}
```
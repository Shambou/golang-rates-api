# Golang Rest api playground

## Running locally
- install docker
- run `docker-compose up --build`

If you have task file installed `https://taskfile.dev/` you can just call `task run` from project

## Running tests
- `task test` or `go test -v ./...`


## API Reference

#### Get latest rate for currency

```http
  GET /api/v1/rates/latest?quote_currency={currency}
```

| Parameter | Type     | Description                     |
| :-------- | :------- |:--------------------------------|
| `quote_currency` | `string` | **Required**. Currency ISO code |

#### Get rates for currency in date range

```http
  GET /api/v1/rates/range?quote_currency={currency}&from={from_date}&to={to_date}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `quote_currency`      | `string` | **Required**. Currency ISO code |
| `from`      | `string` | **Required**. Date in format "2006-01-02" |
| `to`      | `string` | **Required**. Date in format "2006-01-02" |

#### Get the all available rates on date

```http
  GET /api/v1/rates/timeseries?date={date}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `date`      | `string` | **Required**. Date in format "2006-01-02" |

#### Create new rate

```http
  POST /api/v1/rates/{currency}
```

| Parameter  | Type     | Description                     |
|:-----------| :------- |:--------------------------------|
| `currency` | `string` | **Required**. Currency ISO code |

##### Example post data

```json
{
	"date": "2020-12-25",
	"rate": "1.022600"
}
```

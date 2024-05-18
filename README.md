# Genesis API Project

This service implements the following API:

## API Endpoints

### `GET` /rate

This endpoint returns the current `USD to UAH` exchange rate using the Coinbase API.

#### Parameters

``No parameters``

#### Response Codes

```
200: Returns the actual USD to UAH exchange rate.
400: Invalid status value.
```

---

### `POST` /subscribe

This endpoint adds an email address to the database and automatically subscribes it to the USD to UAH exchange rate newsletter.

_The code includes the ability to subscribe to other rates for future development, but this functionality is not currently used to fulfill the requirements._

#### Parameters

``email`` **string** (formData): The email address to be added to the database and the mailing list.

#### Response Codes

```
200: The email address is added to the database and subscribed to the mailing list.
409: The email address already exists.
```

_Not mentioned in the task, but arose during the development process:_
```
400: The provided data (such as email address) is invalid.
500: Internal error status.
```

---

### `POST` /sendEmails

This endpoint sends the current `USD to UAH` exchange rate to subscribed email addresses using goroutines.

#### Parameters

``No parameters``

#### Response Codes

```
200: Emails were sent.
```

## Usage:

- Using Makefile
```
git clone https://github.com/vladyslavpavlenko/genesis-api-project.git
cd genesis-api-project
make up_build
```

Now you can reach an API using [`http://localhost:8080/`](http://localhost:8081/).
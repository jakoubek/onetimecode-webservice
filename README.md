# onetimecode-webservice
Webservice for one-time codes. Access the API at https://api.onetimecode.net/.

**onetimecode-webservice** is a webservice that encapsulates the [onetimecode package](https://github.com/jakoubek/onetimecode) for Golang.

## Usage

```
curl -X GET 'https://api.onetimecode.net/number'

curl -X GET 'https://api.onetimecode.net/number?length=10'

curl -X GET 'https://api.onetimecode.net/number.txt'

curl -X GET 'https://api.onetimecode.net/alphanumeric'

curl -X GET 'https://api.onetimecode.net/alphanumeric?length=40'

curl -X GET 'https://api.onetimecode.net/alphanumeric?length=40&case=uppercase'

curl -X GET 'https://api.onetimecode.net/alphanumeric?case=lowercase'

curl -X GET 'https://api.onetimecode.net/ksuid'

curl -X GET 'https://api.onetimecode.net/uuid'

curl -X GET 'https://api.onetimecode.net/uuid?withoutdashes=true'
```

### Response

```json
{
  "code": 648197
}
```

```json
{
  "code": "vff8GQ"
}
```

## Endpoints

- `/number` returns a numerical code
- `/alphanumeric` returns an alphanumerical code
- `/uuid` returns an UUID
- `/ksuid` returns a KSUID (see [segmentio/ksuid](https://github.com/segmentio/ksuid))
- `/dice` is a shortcut for `/number?min=1&max=6`
- `/coin` is a shortcut for `/number?min=0&max=1` and returns `heads` or `tails`  
- `/status` returns a status object
- `/healthz` is the health endpoint

## Optional parameters

### length

The `length` parameter determines the length of the returned code. The default length is 6. Setting the length to 0 or less uses the default of 6. The highest allowed length is 100. A length higher than 100 gets reduced to 100.

### min/max

Instead of the `length` parameter you can submit a pair of `min` and `max` parameters for the lower and upper threshold of the numerical onetimecode. 

### format

Append a `.txt` to the route name (i.e. `/alphanumeric.txt?length=12`) to get the result as a plain text.

## Installation

```
go get -u github.com/jakoubek/onetimecode-webservice
```

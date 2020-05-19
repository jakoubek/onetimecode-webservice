# onetimecode-webservice
Webservice for one-time codes. Access the API at https://api.onetimecode.net/onetime.

**onetimecode-webservice** is a webservice that encapsulates the [onetimecode package](https://github.com/jakoubek/onetimecode) for Golang.

## Usage

```
curl -X GET 'https://api.onetimecode.net/onetime'

curl -X GET 'https://api.onetimecode.net/onetime?mode=numbers&length=10'

curl -X GET 'https://api.onetimecode.net/onetime?mode=alphanum&length=18'

curl -X GET 'https://api.onetimecode.net/onetime?mode=alphanumuc&length=40'

curl -X GET 'https://api.onetimecode.net/onetime?format=txt'

curl -X GET 'https://api.onetimecode.net/ksuid'
```

### Response

```json
{
  "result": "OK",
  "code": "648197",
  "mode": "numbers",
  "length": 6
}
```

## Endpoints

- `/onetime` returns a one-time code
- `/ksuid` returns a KSUID (see [segmentio/ksuid](https://github.com/segmentio/ksuid))
- `/status` returns a status object

## Optional parameters

### mode

The `mode` parameter determines the kind of onetimecode that is returned:

- numbers (default): only numbers (i.e. 123456)
- alphanum: alphanumeric with lowercase and uppercase characters (i.e. rtXV7u)
- alphanumuc: alphanumeric with only uppercase characters (i.e. U9KPJM)

### length

The `length` parameter determines the length of the returned code. The default length is 6. Setting the length to 0 or less uses the default of 6. The highest allowed length is 100. A length higher than 100 gets reduced to 100.

### format

The `format` parameter determines the format of the response. Default is `json` - it returns a JSON object.
The only other format `txt` returns only the onetimecode as a text/plain answer.

## Installation

```
go get -u github.com/jakoubek/onetimecode-webservice
```

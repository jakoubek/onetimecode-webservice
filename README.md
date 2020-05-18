# onetimecode-webservice
Webservice for one-time codes

**onetimecode-webservice** is a webservice that encapsulates the [onetimecode package](https://github.com/jakoubek/onetimecode) for Golang.

## Usage

```
curl -X GET 'http://api.baneland.de/onetime'

curl -X GET 'http://api.baneland.de/onetime?mode=numbers&length=10'

curl -X GET 'http://api.baneland.de/onetime?mode=alphanum&length=18'

curl -X GET 'http://api.baneland.de/onetime?mode=alphanumuc&length=40'

curl -X GET 'http://api.baneland.de/onetime?format=txt'
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

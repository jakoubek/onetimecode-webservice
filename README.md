# onetimecode-webservice
Webservice for one-time codes. Access the API at https://api.onetimecode.net/.

## Usage

### Number codes

```shell
# get a number (default length: 6)
curl -X GET 'https://api.onetimecode.net/number'
{"result":"883647"}

# get a number with the length of 10 digits
curl -X GET 'https://api.onetimecode.net/number?length=10'
{"result":"1916735822"}

# get a number between 10 and 30
curl -X GET 'https://api.onetimecode.net/number?min=10&max=30'
{"result":"16"}

# get a number with the length of 12 digits, grouped every 4 digits with a dash
curl -X GET 'https://api.onetimecode.net/number?length=12&group_by=-&group_every=4'
{"result":"3650-3264-6315"}
```

### Alphanumeric codes

```shell
# get an alphanumeric code (default length: 6)
curl -X GET 'https://api.onetimecode.net/alphanumeric'
{"result":"crG0Jr"}

# get an alphanumeric code with the length of 40
curl -X GET 'https://api.onetimecode.net/alphanumeric?length=40'
{"result":"8HrEYY2QPAmaKnrAXE1N6oJM7PgvF8LPnRfhfAym"}

# get an alphanumeric code with the length of 20 with all chars UPPERcased
curl -X GET 'https://api.onetimecode.net/alphanumeric?length=20&case=upper'
{"result":"Q41HQOWcEwUakThSA8U7"}

# get an alphanumeric code with the default length with all chars lowerCASED
curl -X GET 'https://api.onetimecode.net/alphanumeric?case=lower'
{"result":"2kf301"}

# get an alphanumeric 40 chars uppercased code, grouped every 4 digits with a dash
curl -X GET 'https://api.onetimecode.net/alphanumeric?length=40&case=upper&group_by=-&group_every=4'
{"result":"YG63-6PPZ-WKFU-7GYG-D63W-LG4M-YQ6D-91S2-3N3A-0KZJ"} 
```

### K-Sortable Globally Unique IDs (ksuid)

```shell
curl -X GET 'https://api.onetimecode.net/ksuid'
{"result":"28iEBXva5OhYVzvDK4iMmBZrP74"}
```

### UUIDs

```shell
# get an UUID
curl -X GET 'https://api.onetimecode.net/uuid'
{"result":"7cc1e97d-6bed-4d63-be64-d0c74dc0c587"}

# get an UUID but without dashes
curl -X GET 'https://api.onetimecode.net/uuid?withoutdashes'
{"result":"2dd9c101d16644eab2df9c048b5d5285"}
```

### Shortcuts for numbers

#### Roll the dice

The `/dice` endpoint is an alias for `/number?min=1&max=6`.
```shell
curl -X GET 'https://api.onetimecode.net/dice'
{"result":"5"}
```

#### Toss a coin

The `/coin` endpoint is an alias for `number?min=0&max=1`. It returns either 0 or 1.
The `/coin` endpoint is the only endpoint that returns an additional attribute `side` whereas 0 = *head* and 1 = *tails*.

```shell
curl -X GET 'https://api.onetimecode.net/coin'
{"result":"1", "side":"tails"}
```

### Response

The API returns a JSON object with the key `result` and a string value.

```json
{
  "result": "648197"
}
```

For the `/coin` endpoint:

```json
{
  "result": "0",
  "side": "head"
}
```

## Endpoints

- `/number` returns a numerical code
- `/alphanumeric` returns an alphanumerical code
- `/uuid` returns an UUID
- `/ksuid` returns a KSUID (see [segmentio/ksuid](https://github.com/segmentio/ksuid))
- `/dice` is a shortcut for `/number?min=1&max=6`
- `/coin` is a shortcut for `/number?min=0&max=1` and returns `heads` or `tails`  
- `/healthz` is the health endpoint

## Optional parameters

### length

The `length` parameter determines the length of the returned code. The default length is 6. Setting the length to 0 or less uses the default of 6. The highest allowed length is 100. A length higher than 100 gets reduced to 100.

### min/max

Instead of the `length` parameter you can submit a pair of `min` and `max` parameters for the lower and upper threshold of the numerical onetimecode. 

### group_by/group_every

Groups the returned code in segments of the length of `group_every`, divided by a single character given by `group_by`.
Works with both or one of the params. `group_every` has a default of `4`, `group_by` has a dash (`-`) as default value. Does not work if `group_every` is larger than the requested `length` of the code.

Applies only to `/number` and `/alphanumeric`.

## Installation

```
go get -u github.com/jakoubek/onetimecode-webservice
```

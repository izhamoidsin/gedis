# Gedis
Test task implementation for *GoJuno* company

Redis-like storage implementation
Contains:
- server with HTTP API
- Go Lang client

## Stores key-value pairs where key is always string and value could be:
- string
- array
- or dictionary (string -> string)

## Supported operations:
- Get
- Set
- Update
- Remove
- Keys
- Get i element on list
- Get value by key from dict

## Per key TTL
All the entries will be expired automatically after a certain period of time (1 minute by default) after the last read operation (Set or Update)

# API spec (simplified)

| URI | METHOD | Description |
| --- | --- | --- |
|`/heartbeat`| GET, HEAD | Check the server state |
|`/keys`| GET | Get all the keys |
|`/entries/{key}`| GET | Get stored value by the key|
|`/entries/{key}`| HEAD | Check if there is a value stored with the key|
|`/entries/{key}`| PUT | Update existing value by the key |
|`/entries/{key}`| POST | Store a new with the key|
|`/entries/{key}`| DELETE | Delete stored value by the key |
|`/entries/{key}/elements/{ind}`| GET | Get `i` element of a list entry stored with the key |
|`/entries/{key}/entries/{subKey}`| GET | Get value by `subKey` from dictionary entry stored with the key |


# Build info
To build the app and run test suite execute following command
```
make
```

# Deployment info
To run the server execute following command
```
go run main.go config.go
```
The server will start at `htpp://localhost:8081`

# Examples
## Check Server state
```
curl http://localhost:8081/heartbeat -v
```
```
< HTTP/1.1 200 OK
I'm ok sinse Thursday, 13-Jul-17 13:39:47 +03
```

## Store new string value
```
curl -XPOST http://localhost:8081/entries/test -d '"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."' -v
```
```
< HTTP/1.1 201 Created
```

## Store new dictionary value
```
curl -XPOST http://localhost:8081/entries/233 -d '{"7": "Seven", "8": "Eight", "9": "ooop... Surprise!"}' -v
```
```
< HTTP/1.1 201 Created
```

## Get all keys
```
curl http://localhost:8081/keys -v
```
```
< HTTP/1.1 200 OK
["233","test"]
```

## Get string value back
```
curl http://localhost:8081/entries/test -v
```
```
< HTTP/1.1 200 OK
"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
```

## Get dictionary nested value
```
curl http://localhost:8081/entries/233/entries/8 -v
```
```
< HTTP/1.1 200 OK
"Eight"
```


##Get dictionary nested value (ERROR)
```
curl http://localhost:8081/entries/test/entries/8 -v
```
```
< HTTP/1.1 400 Bad Request
Stored value is not a dictionary
```

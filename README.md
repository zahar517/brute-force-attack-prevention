# Brute Force Attack Prevention Tool

The service is designed to prevent password brute force attack during authorization in any system.

The service is called before the user is authorized and can either allow or block the attempt.

It is assumed that the service is used only for server-server, i.e. hidden from the end user.

## Work algorithm

The service limits the frequency of authorization attempts for various combinations of parameters, for example:
- no more than N = 10 attempts per minute for a given login.
- no more than M = 100 attempts per minute for a given password (reverse brute-force protection).
- no more than K = 1000 attempts per minute for a given IP (the number is large, because NAT).

The service uses the leaky bucket algorithm. Moreover, the service supports many buckets, one for each login/password/ip.

White/black lists contain network addresses, which are processed in a simpler way:
- If the incoming IP is in the blacklist, then the service unconditionally rejects authorization (ok=false);
- If the incoming IP is in the whitelist - allows authorization (ok=true).

## Architecture

The microservice consists of an API, a database for storing settings, and black/white lists. The service provides a GRPC API.

## Configure

For local development, use the following config `.env`, for docker development use the `.env` file and docker-compose files in `./build/` directory.

Main configuration parameters: `LOGIN_LIMIT`, `PASSWORD_LIMIT`, `IP_LIMIT` - upon reaching these limits the service considers an attempt to be brute force.

## Command line interface

A command-line interface is also available for manual administration of the service. Through the CLI, it is possible to call a bucket reset and manage the whitelist / blacklist. For instance:
```
./bin/cli add black 192.168.0.1/32"
./bin/cli rm white 192.168.0.2/32"
```

## Service deployment
Use `make up` command (`docker-compose up` under the hood) for service deployment.

## Testing
For unit-testing use `make test` command from the main project directory.

Also integration tests are available by using `make integration-test` command.

## Api methods

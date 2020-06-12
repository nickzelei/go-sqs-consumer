# Simple SQS Consumer

Currently only prints out messages

## Config

* SQS_QUEUE_URL - the https url to the SQS queue
* MAX_WORKERS - the number of goroutines to spawn (default:1)
* AWS_PROFILE - the AWS profile to pull from for absorbing credentials (default:default)

The config can be set by either creating a `config.yml` file next to the binary or by setting the config values as environment variables

## Install
go install

## Build
make build

## Clean
make clean

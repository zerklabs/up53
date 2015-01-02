up53
====

## Usage

```
Usage of ./bin/linux-amd64/up53:
  -name="": Record name to update
  -ttl=3600: TTL of record to update
  -type="": Record type to update
  -zoneid="": Route53 Hosted Zone ID
```

## Example

See the [goamz documentation](https://github.com/go-amz/amz/blob/v1/aws/aws.go#L186-L205) for AWS environment authentication

```
AWS_ACCESS_KEY_ID=[...] AWS_SECRET_ACCESS_KEY=[...] up53 -zoneid="ABCD.." -ttl=3600 -type="A" -name="www"
```

this pulls the response from canihazip.com for determining the public IP address

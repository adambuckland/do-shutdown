# do-shutdown

Deletes any DigitalOcean droplets on the account. Simple as that!

## Installation

```
go install github.com/adambuckland/do-shutdown
```

You'll also either need to create a [DigitalOcean access token](https://cloud.digitalocean.com/settings/api/tokens), and set it to the `DIGITALOCEAN_ACCESS_TOKEN` environment variable, or provide the token via the `-token` flag

## Usage

```
do-shutdown
```

**Available flags**
```
-dryrun
    Will list the droplets, but won't actually delete them
-token <token>
    The DigitalOcean personal access token to access the API
```

## License

MIT-licensed
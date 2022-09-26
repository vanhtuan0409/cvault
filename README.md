# Cold Vault

An encryption tool that use AWS KMS to perform envelop encryption

The encrypted file then pushed into a vault storage of either:

- Local dir: `local://{dir_name}`
- S3: `s3://{bucket_name}`

Encryption is performed in memory so original file is expected to not be very large

### Usage

```
A cold vault encryption tool

Usage:
  cvault [flags]
  cvault [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  decrypt     Decrypt a file from storage
  encrypt     Encrypt a file and push it into storage
  help        Help about any command
  list        List encrypted file in store
  remove      Remove encrypted file from store

Flags:
  -h, --help            help for cvault
  -k, --key-id string   KMS key id
  -s, --store string    Location of storage (default "local://.")

Use "cvault [command] --help" for more information about a command.
```

### FAQ

```
- Q: Why cold vault
  A: I want to have an infrequently-access vault to store my recovery code, TOTP seed, etc...

- Q: Why not GPG
  A: I dont have a secure way to sync GPG private across machines
```

# Cold Vault

An encryption tool that use Google Tink to perform envelop encryption

The encrypted file then pushed into a vault storage of either:

- Local dir: `local://{dir_name}`
- S3: `s3://{bucket_name}`

Supported KMS:

- AWS KMS: aws-kms://
- GCP KMS: gcp-kms://
- Hashicorp vault: hcvault://
- Passphrase AES: aesgcm://

Encryption is performed in memory so original file is expected to not be very large

### Usage

```
Usage:
  cvault [flags]
  cvault [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  decrypt     Decrypt a file from storage
  encrypt     Encrypt a file and push it into storage
  help        Help about any command
  list        List encrypted file in store
  peek        Peek an encrypted file to stdout
  remove      Remove encrypted file from store

Flags:
  -h, --help                 help for cvault
  -k, --key-id string        KMS key id
      --pass-prompt string   AES passphrase when using aesgcm:// key
  -s, --store string         Location of storage (default "local://.")
      --vault-token string   HC vault token

Use "cvault [command] --help" for more information about a command.
```

### FAQ

```
- Q: Why cold vault
  A: I want to have an infrequently-access vault to store my recovery code, TOTP seed, etc...

- Q: Why not GPG
  A: I dont have a secure way to sync GPG private across machines
```

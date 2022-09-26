# Cold Vault

An encryption tool that use AWS KMS to perform envelop encryption

The encrypted file then pushed into a vault storage of either:

- Local dir: `local://{dir_name}`
- S3: `s3://{bucket_name}`

Encryption is performed in memory so original file is expected to not be very large

### FAQ

- Q: Why cold vault
  A: I want to have an infrequently-access vault to store my recovery code, TOTP seed, etc...

- Q: Why not GPG
  A: I dont have a secure way to sync GPG private across machines

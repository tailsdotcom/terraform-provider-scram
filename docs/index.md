# SCRAM Provider

A provider to hash passwords as RFC 7677 SCRAM.

At the moment it only implements SHA-256 as used by PostgreSQL (PRs welcome).

## Example Usage

```hcl
resource "scram_password" "me" {
  password = "changeme"
}
```

# Password Resource

Provides a method of generating SCRAM from password.



## Example Usage

```hcl
resource "scram_password" "me" {
  password = "changeme"
}

resource "postgresql_role" "me" {
  name     = "me"
  login    = true
  password = "${scram_password.me.scram_mech}\u0024${scram_password.me.iter_count}:${scram_password.me.salt}\u0024${scram_password.me.stored_key}:${scram_password.me.server_key}"
}
```

## Argument Reference

* `password` - (Required) The secret to hash.
* `scram_mech` - Currently only "SCRAM-SHA-256" is implemented.
* `iter_count` - Number of iterations for hash.

## Attribute Reference

* `salt` - A random salt.
* `stored_key` - The hashed client key.
* `server_key` - The sever key.

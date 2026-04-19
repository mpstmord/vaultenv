# vaultenv

CLI tool to inject secrets from HashiCorp Vault into process environments.

## Installation

```bash
go install github.com/yourusername/vaultenv@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultenv/releases).

## Usage

`vaultenv` reads secrets from Vault and injects them as environment variables into a subprocess.

```bash
vaultenv --addr https://vault.example.com \
         --token s.xxxxxxxx \
         --secret secret/data/myapp \
         -- ./myapp --start
```

The secrets stored at the given path will be available as environment variables in the spawned process.

### Options

| Flag | Description | Default |
|------|-------------|---------|
| `--addr` | Vault server address | `$VAULT_ADDR` |
| `--token` | Vault token | `$VAULT_TOKEN` |
| `--secret` | Secret path to read | *(required)* |
| `--prefix` | Prefix for injected env vars | *(none)* |

### Example

```bash
# Inject database credentials into a migration script
vaultenv --secret secret/data/prod/db -- ./migrate up
```

Inside `./migrate`, the secrets are accessible as normal environment variables:

```bash
echo $DB_PASSWORD
```

## Requirements

- Go 1.21+
- HashiCorp Vault 1.x

## License

MIT © [yourusername](https://github.com/yourusername)
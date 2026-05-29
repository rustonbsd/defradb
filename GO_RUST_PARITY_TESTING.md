The following tests should run without failures, that is a good indicator that badger db corekv backend is working:

```sh
DEFRA_CLIENT_CLI=true go test ./cli/test/integration/... -race -shuffle=on -timeout 10m
```
# Example plugin

This is an example scanner plugin.

## Build

```shell
$ go build -buildmode=plugin ./customscanner.go
```

Depending on your OS, it will create a plugin file (e.g. `customscanner.so` for linux).

## Usage

You can now use this plugin with phoneinfoga.

```shell
$ phoneinfoga scan -n <number> --plugin ./customscanner.so

Running scan for phone number <number>...

Results for customscanner
Valid: true
Info: This number is known for scams!

...
```

The `--plugin` flag can be used multiple times to use several plugins at once.

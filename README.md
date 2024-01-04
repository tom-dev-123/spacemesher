# Spacemesher

CLI tool for PoST initialization

## Build
```bash
cd spacemesher
go build -o "./bin/server ./server/*"
go build -o "./bin/client ./client/*"
```

## Usage
### Server
```bash
./bin/server -h
```
```bash
Usage:
  main [OPTIONS]

Application Options:
      --listen=   address for listening (default: 0.0.0.0:8081)
  -d, --data-dir= the post data path

Help Options:
  -h, --help      Show this help message
```
```bash
# On server machine
./bin/server -d /path/to/post_data --listen 0.0.0.0:8088
```
### Client
```bash
./bin/client -h
```
```bash
Usage:
  client [OPTIONS]

Application Options:
      --address=  Address to connect (default: 0.0.0.0:8081)
  -p, --provider= Binging gpu
      --bin=      the path of postcli binary (default: postcli)

Help Options:
  -h, --help      Show this help message
```
```bash
# On client machines
./bin/client --provider 0 --bin /path/to/postcli --address 0.0.0.0:8088
```

## NOTICE
1. The `postdata_metadata.json` file you need to generate with `postcli` first.
2. If the `postcli` binary file is exported on the env on the client machines, you could use the `--bin` tag by default.
3. It should be noted that the `post_data` directory must be pre mounted on the client machines, Recommend using Ansible or Fabric3.

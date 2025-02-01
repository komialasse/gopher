# gopher

A Go networking library for tunneling. You can use this tool to expose local ports to a remote server, allowing TCP connections that bypass standard NAT firewalls.

## Usage

For simple usage, start the server with `gopher run . server`. Then start the client with `go run . local <PORT>` (e.g. `go run . local 5050`). This will expose your local port at `localhost:5050` to traffic that can access machine running the server.

Generally, you would run the server on a machine that you have access to, registered at a particular domain such as `myserver.com`. Then run the client to connect and start the proxy TCP connection.

## Command Line Options

The client command has a small number of options for configuration.

```bash
go run . local <LOCAL_PORT> [-localhost host] [-port p] [-to t]
Usage of local:
  -l, -localhost string
        the local host to expose (default "localhost")
  -p, -port int
        port of remote server (default 8081)
  -to string
        address of remote server (default "localhost")
```

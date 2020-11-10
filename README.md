# Glua-netcl - simple network client for gopher-lua.

Ver. 0.1.0

## Installation
    go get github.com/uamana/glua-net

## Usage
### Go
    import (
        lua "github.com/yuin/gopher-lua"
        nc "github.com/uamana/glua-netcl"
    )

    func main() {
        L := lua.NewState()
        defer L.Close()
        L.PreloadModule("net", nc.Loader)

        ...
    }

### Lua
Glua-netcl uses Go net package. netcl.dial, conn:read, conn:readline, conn:write, conn:set_timeouts and conn:Close ara available.

    local net = require("netcl")

    local conn, err = net.dial("tcp", "127.0.0.1:8899", 1000)
    if err then error(err) end

    local count, err = conn:write("test_a")
    if err then error(err) end

    local resp, err = conn:read(256)
    if err then error(err) end

**netcl.dial(network *string*, address *string*, timeout *number*) -> *net_client*, *string***

  Arguments:
  - network *string* - name of the network ("tcp", "udp", "ip", "ip4", "ip6")
  - adress *string* - address of the remote host
  - timeout *number* - timeout in milliseconds

  Returns:
  - *net_client* - user defined type for connection
  - *string* - error description if any, nil if no errors

**conn:read(count *number*) -> string, string**

  Arguments:
  - count *number* - maximum number of bytes to read

  Returns:
  - string - readed data
  - string - error message or nil

**conn:readline([delim *string*]) -> *string*, *string***

Reads one line ending with a delim. If delim not set, use "\n" as delim.

  Arguments:
  - delim *string* - line delimiter

  Returns:
  - *string* - readed line
  - *string* - error message or nil

**conn:write(data *string*) -> *number*, *string***

  Arguments:
  - data *string* - data to write

  Returns:
  - *number* - number of bytes readed
  - *string* - error message or nil

**conn:set_timeouts(read_timeout *number*, write_timeout *number*)**

Set connection timeouts. If timeouts are zeroes, no timeout used.

  Arguments:
  - read_timeout *number* - timeout for read in milliseconds
  - write_timeout *number* - timeout for write in milliseconds

**conn:close()**

Close connection.
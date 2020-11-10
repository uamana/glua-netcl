package netcl

import (
	"bufio"
	"net"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type luaNetClient struct {
	net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func checkLuaNetClient(L *lua.LState, n int) *luaNetClient {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaNetClient); ok {
		return v
	}
	L.ArgError(n, luaNetClientTypeName+" expected")
	return nil
}

// Dial lua: netcl.Dial(network, address, timeout)
// Connect to address in network with timeout. If timeot is 0, then no timeout is used.
// Returns: net_client_ud
func Dial(L *lua.LState) int {
	network := L.CheckString(1)
	addr := L.CheckString(2)
	timeout := L.CheckInt64(3)

	var (
		conn net.Conn
		err  error
	)
	n := luaNetClient{}
	if timeout > 0 {
		conn, err = net.DialTimeout(network, addr, time.Duration(timeout)*time.Millisecond)
	} else {
		conn, err = net.Dial(network, addr)
	}
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	n.Conn = conn
	ud := L.NewUserData()
	ud.Value = &n
	L.SetMetatable(ud, L.GetTypeMetatable(luaNetClientTypeName))
	L.Push(ud)
	return 1
}

// SetTimeouts lua: net_client_ud.set_timeouts(read_timeout, write_timeout).
// Sets connection timeouts in milliseconds.
// Returns none
func SetTimeouts(L *lua.LState) int {
	n := checkLuaNetClient(L, 1)
	rd := L.CheckInt64(2)
	wd := L.CheckInt64(3)
	n.readTimeout = time.Duration(rd) * time.Millisecond
	n.writeTimeout = time.Duration(wd) * time.Millisecond
	return 0
}

// Read lua: net_client_ud:read(size).
// Reads size of bytes from connection.
// Returns data, err.
// Where data is a string of readed bytes.
func Read(L *lua.LState) int {
	n := checkLuaNetClient(L, 1)
	size := L.CheckInt(2)

	buf := make([]byte, size)
	if (n.readTimeout) > 0 {
		n.SetReadDeadline(time.Now().Add(n.readTimeout))
	}
	count, err := n.Read(buf)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(string(buf[0:count])))
	return 1
}

// ReadLine lua: net_client_ud:read_line([delim])
// Reads line from network connection, terminated by delim. If no delim set, use '\n' as delim.
// Returns: data, err
func ReadLine(L *lua.LState) int {
	n := checkLuaNetClient(L, 1)

	var delim byte
	if L.GetTop() > 1 {
		d := L.CheckString(2)
		if len(d) != 1 {
			L.ArgError(2, "Delim must be string of length 1")
		}
		delim = byte(d[0])
	} else {
		delim = '\n'
	}

	reader := bufio.NewReader(n.Conn)
	if n.readTimeout > 0 {
		n.SetReadDeadline(time.Now().Add(n.readTimeout))
	}
	data, err := reader.ReadString(delim)
	if err != nil {
		L.Push(lua.LString(data))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(data))
	L.Push(lua.LNil)
	return 2
}

// Write lua: net_client_ud:write(data)
// Writes data(string) to network.
// Returns: count, err.
// Where count is a number of bytes readed.
func Write(L *lua.LState) int {
	n := checkLuaNetClient(L, 1)
	data := L.CheckString(2)

	if (n.writeTimeout) > 0 {
		n.SetWriteDeadline(time.Now().Add(n.readTimeout))
	}
	count, err := n.Write([]byte(data))
	if err != nil {
		L.Push(lua.LNumber(count))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(count))
	return 1
}

// Close lua: net_client_ud:close().
// Closes network connection
func Close(L *lua.LState) int {
	n := checkLuaNetClient(L, 1)
	if n != nil {
		n.Close()
	}
	return 0
}

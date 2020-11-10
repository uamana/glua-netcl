package netcl

import (
	"log"
	"net"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const (
	listenAddress = "127.0.0.1:8899"
	bufSize       = 256
)

const testScript = `
local net = require("netcl")
local conn, err = net.dial("tcp", "127.0.0.1:8899", 1000)
if err then error(err) end
-- write and read test
local count, err = conn:write("test_a")
if err then error(err) end
local resp, err = conn:read(256)
if err then error(err) end
if resp ~= "test_a" then error("Expected 'test_a', got " .. resp) end
--readline test
local count, err = conn:write("test_b\n")
if err then error(err) end
local resp, err = conn:readline()
if err then error(err) end
if resp ~= "test_b\n" then error("Expected 'test_b', got " .. resp) end
conn:close()
`

func init() {
	go echoSrv()
	time.Sleep(1 * time.Second)
}

func echoSrv() {
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	buf := make([]byte, bufSize)
nextconn:
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		for {
			count, err := conn.Read(buf)
			if err != nil {
				log.Printf("Read error: %s", err)
				continue nextconn
			}
			_, err = conn.Write(buf[0:count])
			if err != nil {
				log.Printf("Write error: %s", err)
				continue nextconn
			}
		}
	}
}

func TestNetCl(t *testing.T) {
	L := lua.NewState()
	Preload(L)
	err := L.DoString(testScript)
	if err != nil {
		t.Fatal(err)
	}
}

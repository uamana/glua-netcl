package netcl

import lua "github.com/yuin/gopher-lua"

const (
	luaNetClientTypeName = "net_conn"
	luaDefaultModName    = "netcl"
)

var exports = map[string]lua.LGFunction{
	"dial": Dial,
}

func registerNetClientType(L *lua.LState) {
	methods := map[string]lua.LGFunction{
		"read":         Read,
		"readline":     ReadLine,
		"write":        Write,
		"close":        Close,
		"set_timeouts": SetTimeouts,
	}

	mt := L.NewTypeMetatable(luaNetClientTypeName)
	L.SetGlobal(luaNetClientTypeName, mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), methods))
}

// Loader loads module into givel Lua state
func Loader(L *lua.LState) int {
	registerNetClientType(L)

	t := L.NewTable()
	L.SetFuncs(t, exports)
	L.Push(t)
	return 1
}

// Preload add netcl to given lua state's package preload table.
// Lua: local netcl = requirte("netcl")
func Preload(L *lua.LState) {
	L.PreloadModule(luaDefaultModName, Loader)
}

package email

import (
	lua "github.com/yuin/gopher-lua"
)

var VERSION string = "1.0.8"
var LNAME string = "email"

var API = map[string]lua.LGFunction{
	// General
	"new":     Lemail_LNew,
	"newPool": LemailPool_LNew,
}

var API_Lemail = map[string]lua.LGFunction{
	// General
	"from":    Lemail_From,
	"to":      Lemail_To,
	"bcc":     Lemail_BCC,
	"cc":      Lemail_CC,
	"subject": Lemail_Subject,
	"text":    Lemail_Text,
	"html":    Lemail_Html,
	"attach":  Lemail_Attach,
	"send":    Lemail_Send,
	"sendTLS": Lemail_SendTLS,
}

var API_LemailPool = map[string]lua.LGFunction{
	// General
	"send":  LemailPool_Send,
	"close": LemailPool_Close,
}

func Preload(L *lua.LState) {
	L.PreloadModule(LNAME, Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	luaEmail := L.NewTypeMetatable("email")
	L.SetField(luaEmail, "__index", L.SetFuncs(L.NewTable(), API_Lemail))
	t.RawSetH(lua.LString("email"), luaEmail)

	luaPool := L.NewTypeMetatable("pool")
	L.SetField(luaPool, "__index", L.SetFuncs(L.NewTable(), API_LemailPool))
	t.RawSetH(lua.LString("pool"), luaPool)

	t.RawSetH(lua.LString("__version__"), lua.LString(VERSION))

	L.SetFuncs(t, API)
	L.Push(t)
	return 1
}

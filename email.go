package email

import (
	"crypto/tls"
	"errors"
	"net/smtp"
	"time"

	gomail "github.com/jordan-wright/email"
	lua "github.com/yuin/gopher-lua"
)

type Lemail struct {
	Req *gomail.Email
}

func Lemail_New(L *lua.LState, gs *gomail.Email) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = gs
	L.SetMetatable(ud, L.GetTypeMetatable("email"))
	return ud
}
func Lemail_LNew(L *lua.LState) int {
	L.Push(Lemail_New(L, gomail.NewEmail()))
	return 1
}

func LemailPool_New(L *lua.LState, gs *gomail.Pool) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = gs
	L.SetMetatable(ud, L.GetTypeMetatable("pool"))
	return ud
}
func LemailPool_LNew(L *lua.LState) int {
	// !TODO
	addr := L.CheckString(1)
	size := L.CheckInt(2)
	auth := L.CheckTable(3)
	tlscf := L.NewTable()
	if L.GetTop() > 3 {
		tlscf = L.CheckTable(4)
	}
	authConfig, err := util_GetAuth(L, auth, 3)
	tlsConfig := util_GetTLSConfig(L, tlscf)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	pool, perr := gomail.NewPool(
		addr,
		size,
		authConfig,
		tlsConfig,
	)
	if perr != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(perr.Error()))
		return 2
	}
	L.Push(LemailPool_New(L, pool))
	L.Push(lua.LNil)
	return 2
}

func check_LemailPool(L *lua.LState, index int) *gomail.Pool {
	ud := L.CheckUserData(index)
	if v, ok := ud.Value.(*gomail.Pool); ok {
		return v
	}
	L.ArgError(1, "email.pool object expected")
	return nil
}

func check_Lemail(L *lua.LState, index int) *gomail.Email {
	ud := L.CheckUserData(index)
	if v, ok := ud.Value.(*gomail.Email); ok {
		return v
	}
	L.ArgError(1, "email object expected")
	return nil
}

func check_Lemail_UD(L *lua.LState, index int) *lua.LUserData {
	ud := L.CheckUserData(index)
	if _, ok := ud.Value.(*gomail.Email); ok {
		return ud
	}
	L.ArgError(1, "email object expected")
	return nil
}

func check_LDuration(L *lua.LState, index int) time.Duration {
	v := L.CheckAny(index)
	if v.Type() == lua.LTUserData {
		ud, _ := v.(*lua.LUserData)
		if d, ok := ud.Value.(time.Duration); ok {
			return d
		}
	}
	if v.Type() == lua.LTNumber {
		rd, _ := v.(lua.LNumber)
		d := time.Millisecond * time.Duration(rd)
		return d
	}
	L.ArgError(1, "email object expected")
	return time.Duration(0)
}

func LemailPool_Close(L *lua.LState) int {
	pool := check_LemailPool(L, 1)
	pool.Close()
	return 0
}
func LemailPool_Send(L *lua.LState) int {
	pool := check_LemailPool(L, 1)
	eml := check_Lemail(L, 2)
	timeout := check_LDuration(L, 3)
	if timeout < time.Duration(1) {
		L.ArgError(3, "Invalid timeout. Use time userdate or Lua number as milliseconds value.")
		return 0
	}
	err := pool.Send(eml, timeout)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func Lemail_From(L *lua.LState) int {
	eml := check_Lemail(L, 1)
	v := L.CheckString(2)
	eml.From = v
	ud := check_Lemail_UD(L, 1)
	L.Push(ud)
	return 1
}

func Lemail_To(L *lua.LState) int {
	eml := check_Lemail(L, 1)
	v := L.CheckTable(2)
	opt := []string{}
	v.ForEach(func(a lua.LValue, b lua.LValue) {
		if b.Type() == lua.LTString {
			opt = append(opt, b.String())
		}
	})
	eml.To = opt
	ud := check_Lemail_UD(L, 1)
	L.Push(ud)
	return 1
}

func Lemail_BCC(L *lua.LState) int {
	eml := check_Lemail(L, 1)
	v := L.CheckTable(2)
	opt := []string{}
	v.ForEach(func(a lua.LValue, b lua.LValue) {
		if b.Type() == lua.LTString {
			opt = append(opt, b.String())
		}
	})
	eml.Bcc = opt
	ud := check_Lemail_UD(L, 1)
	L.Push(ud)
	return 1
}

func Lemail_CC(L *lua.LState) int {
	eml := check_Lemail(L, 1)
	v := L.CheckTable(2)
	opt := []string{}
	v.ForEach(func(a lua.LValue, b lua.LValue) {
		if b.Type() == lua.LTString {
			opt = append(opt, b.String())
		}
	})
	eml.Cc = opt
	ud := check_Lemail_UD(L, 1)
	L.Push(ud)
	return 1
}

func Lemail_Subject(L *lua.LState) int {
	eml := check_Lemail(L, 1)
	v := L.CheckString(2)
	eml.Subject = v
	ud := check_Lemail_UD(L, 1)
	L.Push(ud)
	return 1
}

func Lemail_Text(L *lua.LState) int {
	eml := check_Lemail(L, 1)
	v := L.CheckString(2)
	eml.Text = []byte(v)
	ud := check_Lemail_UD(L, 1)
	L.Push(ud)
	return 1
}

func Lemail_Html(L *lua.LState) int {
	eml := check_Lemail(L, 1)
	v := L.CheckString(2)
	eml.HTML = []byte(v)
	ud := check_Lemail_UD(L, 1)
	L.Push(ud)
	return 1
}

func Lemail_Attach(L *lua.LState) int {
	eml := check_Lemail(L, 1)
	ud := check_Lemail_UD(L, 1)
	v := L.CheckString(2)
	_, err := eml.AttachFile(v)
	if err != nil {
		L.Push(ud)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(ud)
	L.Push(lua.LNil)
	return 2
}

func Lemail_Send(L *lua.LState) int {
	eml := check_Lemail(L, 1)
	addr := L.CheckString(2)
	auth := L.CheckTable(3)
	authConfig, err := util_GetAuth(L, auth, 3)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	err = eml.Send(
		addr,
		authConfig,
	)
	//ud := check_Lemail_UD(L, 1)
	//L.Push(ud)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push((lua.LNil))
	}
	return 1
}

func Lemail_SendTLS(L *lua.LState) int {
	eml := check_Lemail(L, 1)
	addr := L.CheckString(2)
	auth := L.CheckTable(3)
	tlsCfg := L.CheckTable(4)
	authConfig, err := util_GetAuth(L, auth, 3)
	tlsConfig := util_GetTLSConfig(L, tlsCfg)
	if err != nil {
		return 0
	}
	err = eml.SendWithTLS(
		addr,
		authConfig,
		tlsConfig,
	)
	ud := check_Lemail_UD(L, 1)
	L.Push(ud)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push((lua.LNil))
	}
	return 2
}

func util_GetTLSConfig(L *lua.LState, t lua.LValue) *tls.Config {
	cfg := &tls.Config{}
	if t.Type() == lua.LTTable {
		tx, _ := t.(*lua.LTable)
		v := tx.RawGetH(lua.LString("insecureSkipVerify"))
		if v.Type() == lua.LTBool {
			vx, ok := v.(lua.LBool)
			if ok {
				cfg.InsecureSkipVerify = bool(vx)
			}
		}
		v2 := tx.RawGetH(lua.LString("serverName"))
		if v2.Type() == lua.LTString {
			cfg.ServerName = v2.String()
			// !TODO make luax crypto x509 cert lib to be able to
			// add root CAs so we don't have to use insecureSkipVerify
			// in our TLSConfig!
		}
	}
	// !TODO
	return cfg
}

func util_GetAuth(L *lua.LState, auth *lua.LTable, param int) (smtp.Auth, error) {
	identity := auth.RawGet(lua.LString("identity"))
	username := auth.RawGet(lua.LString("username"))
	password := auth.RawGet(lua.LString("password"))
	host := auth.RawGet(lua.LString("host"))
	if identity.Type() == lua.LTNil {
		identity = lua.LString("")
	}
	if identity.Type() != lua.LTString ||
		username.Type() != lua.LTString ||
		password.Type() != lua.LTString ||
		host.Type() != lua.LTString {
		L.ArgError(param, "email: Invalid parameter for smtp auth You must provide a username, password and host!")
		return nil, errors.New("invalid smtp auth!")
	}
	return smtp.PlainAuth(identity.String(), username.String(), password.String(), host.String()), nil
}

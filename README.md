# An email client (smtp) for gopher-lua

A simple email client for [gopher-lua](https://github.com/yuin/gopher-lua).
This library is based on [https://github.com/jordan-wright/email](https://github.com/jordan-wright/email) for its email functionallity.

## Example

```lua
	local email = require("email")
	
	local myEmail = email.new()
	myEmail:from("address@domain.com")
	myEmail:to({
		"someemail@domain.com"
	})
	myEmail:html([[
		<h1>Hello</h1>
		<p>User</p>
	]])
	
	myEmail:text([[
hello
user
	]])

	myEmail:Lemail_SendTLS("account@email.com", {
		-- Auth config
		identity = "address@fomain.com";
		username = "username";
		password = "password";
		host = "mail.domain.com";
	}, {
		-- TLS config
		insecureSkipVerify = true;
	})
	
```

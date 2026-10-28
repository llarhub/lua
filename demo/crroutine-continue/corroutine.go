package main

import (
	"github.com/goplus/lib/c"
	"github.com/llarhub/lua"
)

func countdownContinue(L *lua.State, status c.Int, ctx lua.KContext) c.Int {
	return countdown(L)
}

func countdown(L *lua.State) c.Int {
	arg := L.Tointeger(lua.Upvalueindex(1))
	c.Printf(c.Str("resume called with %d\n"), arg)
	if arg > 0 {
		L.Pushinteger(arg - 1)
		L.Replace(lua.Upvalueindex(1))
		L.Pushinteger(arg)
		return L.Yieldk(1, lua.KContext(0), countdownContinue)
	}
	return 0
}

func createCountdown(L *lua.State) c.Int {
	max := L.LCheckinteger(1)
	L.Pushinteger(max)
	L.Pushcclosure(countdown, 1)
	return 1
}

func main() {
	L := lua.LNewstate()
	L.LOpenlibs()
	defer L.Close()
	L.Register(c.Str("create_countdown"), createCountdown)

	testcode := c.Str(`
 local countdown = create_countdown(5)
 local co = coroutine.create(countdown)
 while true do
     local success, value = coroutine.resume(co)
     if not success then
         print('Error:', value)
         break
     end
     if value == nil then break end
     print('Lua received:', value)
 end
 print('Countdown finished');
	`)

	if L.LDostring(testcode) != lua.OK {
		c.Printf(c.Str("Error: %s\n"), L.Tostring(-1))
	}

	L.LOpenlibs()
}

/* Expected output:
resume called with 5
Lua received:   5
resume called with 4
Lua received:   4
resume called with 3
Lua received:   3
resume called with 2
Lua received:   2
resume called with 1
Lua received:   1
resume called with 0
Countdown finished
*/

package main

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
	"github.com/llarhub/lua"
)

func main() {
	L := lua.LNewstate()
	defer L.Close()

	L.LOpenlibs()
	if res := L.LLoadstring(c.Str("function doubleNumber(x) return x * 2 end")); res != lua.OK {
		c.Printf(c.Str("error: %s\n"), L.Tostring(-1))
	}
	if res := L.Pcall(0, 0, 0); res != lua.OK {
		c.Printf(c.Str("error: %s\n"), L.Tostring(-1))
	}

	L.Getglobal(c.Str("doubleNumber"))
	L.Pushnumber(10)

	if res := L.Pcall(1, 1, 0); res != lua.OK {
		c.Printf(c.Str("error: %s\n"), L.Tostring(-1))
	}

	if res := L.Isnumber(-1); res != 0 {
		result := L.Tointeger(-1)
		c.Printf(c.Str("result: %lld\n"), result)
	} else {
		c.Printf(c.Str("error: %s\n"), L.Tostring(-1))
	}
}

/* Expected output:
result: 20
*/

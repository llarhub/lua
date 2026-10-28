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
	if res := L.LDostring(c.Str("function combineParams(num, str) return 'Result: ' .. str .. '  ' .. num end")); res != lua.OK {
		c.Printf(c.Str("error: %s\n"), L.Tostring(-1))
	}
	L.Getglobal(c.Str("combineParams"))
	L.Pushnumber(3.14159)
	L.Pushstring(c.Str("Hello, World!"))
	if res := L.Pcall(2, 1, 0); res != lua.OK {
		c.Printf(c.Str("error: %s\n"), L.Tostring(-1))
	}
	if res := L.Isstring(-1); res != 0 {
		result := L.Tostring(-1)
		c.Printf(result)
	}
}

/* Expected output:
Result: Hello, World!  3.14159
*/

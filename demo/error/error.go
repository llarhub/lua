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
	if res := L.LLoadstring(c.Str("function doubleNumber(x) ! return x * 2 end")); res != lua.OK {
		c.Printf(c.Str("error: %s\n"), L.Tostring(-1))
	}
}

/* Expected output:
error: [string "function doubleNumber(x) ! return x * 2 end"]:1: unexpected symbol near '!'
*/

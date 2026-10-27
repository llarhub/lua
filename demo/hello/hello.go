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
	if res := L.LDostring(c.Str("print('hello world')")); res != lua.OK {
		println("error")
	}
}

/* Expected output:
hello world
*/

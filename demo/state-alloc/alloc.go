package main

import (
	"github.com/goplus/lib/c"
	"github.com/llarhub/lua"
)

func alloc(ud c.Pointer, ptr c.Pointer, osize c.SizeT, nsize c.SizeT) c.Pointer {
	if nsize == 0 {
		c.Free(ptr)
		return nil
	} else {
		return c.Realloc(ptr, uintptr(nsize))
	}
}

func main() {
	L := lua.Newstate(alloc, nil, 0)
	defer L.Close()
	L.LOpenlibs()
	if res := L.LDostring(c.Str("print('new state success')")); res != lua.OK {
		println("newstate error")
	}

	allocf := L.Getallocf(nil)
	L.Setallocf(allocf, nil)

	if res := L.LDostring(c.Str("print('set newstate success')")); res != lua.OK {
		println("set newstate error")
	}
}

/* Expected output:
new state success
set newstate success
*/

package lua

import (
	"github.com/goplus/lib/c"
)

// #define lua_pop(L,n)		lua_settop(L, -(n)-1)
func (L *State) Pop(n c.Int) {
	L.Settop(-(n) - 1)
}

func (L *State) LDofile(filename *c.Char) c.Int {
	if loadResult := L.LLoadfile(filename); loadResult != 0 {
		return loadResult
	}
	return L.Pcall(c.Int(0), c.Int(MULTRET), c.Int(0))
}

func (L *State) LDostring(str *c.Char) c.Int {
	if loadResult := L.LLoadstring(str); loadResult != 0 {
		return loadResult
	}
	return L.Pcall(c.Int(0), c.Int(MULTRET), c.Int(0))
}

package lua

import (
	"github.com/goplus/lib/c"
)

func (L *State) Isfunction(n c.Int) bool      { return L.Type(n) == c.Int(TFUNCTION) }
func (L *State) Istable(n c.Int) bool         { return L.Type(n) == c.Int(TTABLE) }
func (L *State) Islightuserdata(n c.Int) bool { return L.Type(n) == c.Int(TLIGHTUSERDATA) }
func (L *State) Isnil(n c.Int) bool           { return L.Type(n) == c.Int(TNIL) }
func (L *State) Isboolean(n c.Int) bool       { return L.Type(n) == c.Int(TBOOLEAN) }
func (L *State) Isthread(n c.Int) bool        { return L.Type(n) == c.Int(TTHREAD) }
func (L *State) Isnone(n c.Int) bool          { return L.Type(n) == c.Int(TNONE) }
func (L *State) Isnoneornil(n c.Int) bool     { return L.Type(n) <= 0 }

/*
 * Pseudo-indices
 * (-LUAI_MAXSTACK is the minimum valid index; we keep some free empty
 * space after that to help overflow detection)
 */
func Upvalueindex(i c.Int) c.Int {
	return c.Int(REGISTRYINDEX) - i
}

// #define lua_pop(L,n)		lua_settop(L, -(n)-1)
func (self *State) Pop(n c.Int) {
	self.Settop(-(n) - 1)
}

// #define lua_remove(L,idx) (lua_rotate(L, (idx), -1), lua_pop(L, 1))
func (L *State) Remove(idx c.Int) {
	L.Rotate(idx, -1)
	L.Pop(1)
}

// #define lua_replace(L,idx)  (lua_copy(L, -1, (idx)), lua_pop(L, 1))
func (L *State) Replace(idx c.Int) {
	L.Copy(-1, idx)
	L.Pop(1)
}

// #define lua_register(L,n,f) (lua_pushcfunction(L, (f)), lua_setglobal(L, (n)))
func (L *State) Register(name *c.Char, f CFunction) {
	L.Pushcfunction(f)
	L.Setglobal(name)
}

func (self *State) LDofile(filename *c.Char) c.Int {
	if loadResult := self.LLoadfile(filename); loadResult != 0 {
		return loadResult
	}
	return self.Pcall(c.Int(0), c.Int(MULTRET), c.Int(0))
}

func (self *State) LDostring(str *c.Char) c.Int {
	if loadResult := self.LLoadstring(str); loadResult != 0 {
		return loadResult
	}
	return self.Pcall(c.Int(0), c.Int(MULTRET), c.Int(0))
}

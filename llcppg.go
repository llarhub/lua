package lua

import (
	"unsafe"

	"github.com/goplus/lib/c"
)

/*
@@ LUA_EXTRASPACE defines the size of a raw memory area associated with
** a Lua state with very fast access.
** CHANGE it if you need a different size.
*/
const EXTRASPACE = unsafe.Sizeof(c.Pointer(nil))

func (self *State) Getextraspace() c.Pointer {
	return c.Pointer(uintptr(c.Pointer(self)) - EXTRASPACE)
}

func (self *State) Isfunction(n c.Int) bool      { return self.Type(n) == c.Int(TFUNCTION) }
func (self *State) Istable(n c.Int) bool         { return self.Type(n) == c.Int(TTABLE) }
func (self *State) Islightuserdata(n c.Int) bool { return self.Type(n) == c.Int(TLIGHTUSERDATA) }
func (self *State) Isnil(n c.Int) bool           { return self.Type(n) == c.Int(TNIL) }
func (self *State) Isboolean(n c.Int) bool       { return self.Type(n) == c.Int(TBOOLEAN) }
func (self *State) Isthread(n c.Int) bool        { return self.Type(n) == c.Int(TTHREAD) }
func (self *State) Isnone(n c.Int) bool          { return self.Type(n) == c.Int(TNONE) }
func (self *State) Isnoneornil(n c.Int) bool     { return self.Type(n) <= 0 }

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
func (self *State) Remove(idx c.Int) {
	self.Rotate(idx, -1)
	self.Pop(1)
}

// #define lua_replace(L,idx)  (lua_copy(L, -1, (idx)), lua_pop(L, 1))
func (self *State) Replace(idx c.Int) {
	self.Copy(-1, idx)
	self.Pop(1)
}

// #define lua_register(L,n,f) (lua_pushcfunction(L, (f)), lua_setglobal(L, (n)))
func (self *State) Register(name *c.Char, f CFunction) {
	self.Pushcfunction(f)
	self.Setglobal(name)
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

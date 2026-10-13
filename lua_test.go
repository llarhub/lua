//go:build llgo

package lua_test

import (
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
	lua "github.com/llarhub/lua"
)

const registryIndex = c.Int(-1073742823)

func newState(t *testing.T) *lua.State {
	L := lua.LNewstate()
	if L == nil {
		t.Fatal("luaL_newstate returned nil")
	}
	L.Openselectedlibs(c.Int(-1), c.Int(0))
	t.Cleanup(L.Close)
	return L
}

func pop(L *lua.State, n int) { L.Settop(c.Int(-n - 1)) }

func pcall(L *lua.State, nargs, nresults int) c.Int {
	return L.Pcallk(c.Int(nargs), c.Int(nresults), 0, 0, nil)
}

func call(L *lua.State, nargs, nresults int) { L.Callk(c.Int(nargs), c.Int(nresults), 0, nil) }

func doString(L *lua.State, source string) c.Int {
	if status := L.Loadstring(c.AllocaCStr(source)); status != lua.OK {
		return status
	}
	return pcall(L, 0, lua.MULTRET)
}

func luaString(L *lua.State, idx c.Int) string {
	var n c.SizeT
	p := L.Tolstring(idx, &n)
	if p == nil {
		return ""
	}
	return c.GoString(p)
}

func luaNumber(L *lua.State, idx c.Int) lua.Number {
	return L.Tonumberx(idx, nil)
}

func luaInteger(L *lua.State, idx c.Int) lua.Integer {
	return L.Tointegerx(idx, nil)
}

func setGlobalFunction(L *lua.State, name string, fn lua.CFunction) {
	L.Pushcclosure(fn, 0)
	L.Setglobal(c.AllocaCStr(name))
}

func upvalueIndex(n int) c.Int { return registryIndex - c.Int(n) }

func TestHello(t *testing.T) {
	L := newState(t)
	if status := doString(L, `return "hello world"`); status != lua.OK {
		t.Fatalf("hello: %s", luaString(L, -1))
	}
	if got := luaString(L, -1); got != "hello world" {
		t.Fatalf("result = %q, want hello world", got)
	}
}

func TestLoadCall(t *testing.T) {
	L := newState(t)
	if status := L.Loadstring(c.Str(`return "hello world"`)); status != lua.OK {
		t.Fatalf("load: %s", luaString(L, -1))
	}
	if status := pcall(L, 0, 1); status != lua.OK {
		t.Fatalf("call: %s", luaString(L, -1))
	}
	if got := luaString(L, -1); got != "hello world" {
		t.Fatalf("result = %q", got)
	}
}

func TestSyntaxError(t *testing.T) {
	L := newState(t)
	if status := L.Loadstring(c.Str(`function doubleNumber(x) ! return x * 2 end`)); status != lua.ERRSYNTAX {
		t.Fatalf("status = %d, want syntax error", status)
	}
	if msg := luaString(L, -1); msg == "" {
		t.Fatal("syntax error message is empty")
	}
}

func TestFunctionNumber(t *testing.T) {
	L := newState(t)
	if status := doString(L, `function doubleNumber(x) return x * 2 end`); status != lua.OK {
		t.Fatalf("define: %s", luaString(L, -1))
	}
	L.Settop(0)
	L.Getglobal(c.Str("doubleNumber"))
	L.Pushnumber(10)
	if status := pcall(L, 1, 1); status != lua.OK {
		t.Fatalf("call: %s", luaString(L, -1))
	}
	if got := luaInteger(L, -1); got != 20 {
		t.Fatalf("result = %d, want 20", got)
	}
}

func TestFunctionString(t *testing.T) {
	L := newState(t)
	code := `function processStrings(a, b, c) return a .. b .. c end`
	if status := doString(L, code); status != lua.OK {
		t.Fatalf("define: %s", luaString(L, -1))
	}
	L.Settop(0)
	L.Getglobal(c.Str("processStrings"))
	L.Pushstring(c.Str("Hello, World!"))
	L.Pushlstring(c.Str("Hello Lua In LLGO"), 17)
	L.Pushfstring(c.Str("Hello %s In %d"), c.Str("LLGO"), c.Int(2024))
	if status := pcall(L, 3, 1); status != lua.OK {
		t.Fatalf("call: %s", luaString(L, -1))
	}
	want := "Hello, World!Hello Lua In LLGOHello LLGO In 2024"
	if got := luaString(L, -1); got != want {
		t.Fatalf("result = %q, want %q", got, want)
	}
}

func TestFunctionConcat(t *testing.T) {
	L := newState(t)
	if status := doString(L, `function combineParams(num, str) return "Result: " .. str .. "  " .. num end`); status != lua.OK {
		t.Fatalf("define: %s", luaString(L, -1))
	}
	L.Settop(0)
	L.Getglobal(c.Str("combineParams"))
	L.Pushnumber(3.14159)
	L.Pushstring(c.Str("Hello, World!"))
	if status := pcall(L, 2, 1); status != lua.OK {
		t.Fatalf("call: %s", luaString(L, -1))
	}
	if got := luaString(L, -1); got != "Result: Hello, World!  3.14159" {
		t.Fatalf("result = %q", got)
	}
}

func TestStack(t *testing.T) {
	L := newState(t)
	L.Pushstring(c.Str("Hello"))
	L.Pushstring(c.Str("LLGO"))
	L.Pushnumber(2024)
	if got := L.Gettop(); got != 3 {
		t.Fatalf("top = %d", got)
	}
	if got := L.Absindex(-2); got != 2 {
		t.Fatalf("absindex = %d", got)
	}
	L.Pushvalue(2)
	L.Rotate(1, -1)
	L.Copy(-1, 2)
	if got := L.Gettop(); got != 4 || luaString(L, 1) != "LLGO" || luaString(L, 2) != "Hello" {
		t.Fatalf("unexpected rotated stack: top=%d", got)
	}
	if L.Checkstack(2) == 0 {
		t.Fatal("stack did not grow")
	}
	L.Pushnumber(3.14)
	L.Pushstring(c.Str("Lua"))
	L.Settop(5)
	if got := L.Gettop(); got != 5 || luaString(L, 5) != "3.14" {
		t.Fatalf("settop result: top=%d value=%q", got, luaString(L, 5))
	}
	L1 := L.Newthread()
	pop(L, 1)
	L.Xmove(L1, 2)
	if L1.Gettop() != 2 {
		t.Fatalf("moved top = %d", L1.Gettop())
	}
}

func TestTable(t *testing.T) {
	L := newState(t)
	L.Createtable(0, 0)
	L.Pushstring(c.Str("name"))
	L.Pushstring(c.Str("John"))
	L.Settable(-3)
	L.Pushstring(c.Str("age"))
	L.Pushnumber(30)
	L.Settable(-3)
	L.Pushstring(c.Str("John Doe"))
	L.Setfield(-2, c.Str("fullname"))
	L.Pushinteger(123)
	L.Seti(-2, 1)
	key := c.AllocaCStr("pointer key")
	L.Pushstring(c.Str("pointer value"))
	L.Rawsetp(-2, unsafe.Pointer(key))
	L.Getfield(-1, c.Str("name"))
	if luaString(L, -1) != "John" {
		t.Fatalf("name = %q", luaString(L, -1))
	}
	pop(L, 1)
	L.Pushstring(c.Str("fullname"))
	L.Rawget(-2)
	if luaString(L, -1) != "John Doe" {
		t.Fatalf("fullname = %q", luaString(L, -1))
	}
	pop(L, 1)
	L.Getfield(-1, c.Str("age"))
	if luaNumber(L, -1) != 30 {
		t.Fatalf("age = %v", luaNumber(L, -1))
	}
	pop(L, 1)
	L.Geti(-1, 1)
	if luaInteger(L, -1) != 123 {
		t.Fatalf("index = %d", luaInteger(L, -1))
	}
}

func metatableString(L *lua.State) c.Int {
	L.Pushstring(c.Str("Hello from metatable!"))
	return 1
}

func TestMetatable(t *testing.T) {
	L := newState(t)
	L.Createtable(0, 0)
	L.Createtable(0, 1)
	L.Pushcclosure(metatableString, 0)
	L.Setfield(-2, c.Str("__tostring"))
	if L.Setmetatable(-2) == 0 {
		t.Fatal("Setmetatable failed")
	}
	L.Setglobal(c.Str("obj"))
	L.Getglobal(c.Str("obj"))
	if L.Getmetatable(-1) == 0 {
		t.Fatal("metatable missing")
	}
	L.Pushstring(c.Str("__tostring"))
	L.Gettable(-2)
	if L.Iscfunction(-1) == 0 {
		t.Fatal("__tostring is not a C function")
	}
	call(L, 0, 1)
	if got := luaString(L, -1); got != "Hello from metatable!" {
		t.Fatalf("result = %q", got)
	}
}

func TestUserdata(t *testing.T) {
	L := newState(t)
	p := L.Newuserdatauv(unsafe.Sizeof(int(0)), 1)
	*(*int)(p) = 42
	L.Setglobal(c.Str("data"))
	value := 24
	L.Pushlightuserdata(unsafe.Pointer(&value))
	L.Setglobal(c.Str("lightdata"))
	L.Getglobal(c.Str("data"))
	if L.Isuserdata(-1) == 0 || *(*int)(L.Touserdata(-1)) != 42 {
		t.Fatal("userdata mismatch")
	}
	pop(L, 1)
	L.Getglobal(c.Str("lightdata"))
	if L.Type(-1) != lua.TLIGHTUSERDATA || *(*int)(L.Touserdata(-1)) != 24 {
		t.Fatal("light userdata mismatch")
	}
}

func TestThread(t *testing.T) {
	L := newState(t)
	if L.Pushthread() == 0 {
		t.Fatal("main thread not reported as main")
	}
	pop(L, 1)
	co := L.Newthread()
	if co.Pushthread() != 0 {
		t.Fatal("new thread reported as main")
	}
	if co.Tothread(-1) != co {
		t.Fatal("Tothread returned a different state")
	}
	if co.Status() != lua.OK {
		t.Fatalf("new thread status = %d", co.Status())
	}
	if L.Closethread(co) != lua.OK {
		t.Fatal("Closethread failed")
	}
}

func TestCoroutine(t *testing.T) {
	L := newState(t)
	if status := doString(L, `function coro_func() for i = 1, 5 do coroutine.yield(i) end end`); status != lua.OK {
		t.Fatalf("define: %s", luaString(L, -1))
	}
	L.Getglobal(c.Str("coro_func"))
	co := L.Newthread()
	L.Pushvalue(-2)
	L.Xmove(co, 1)
	var nres c.Int
	for want := int64(1); ; want++ {
		status := co.Resume(nil, 0, &nres)
		if status == lua.YIELD {
			if luaInteger(co, -1) != lua.Integer(want) {
				t.Fatalf("yield = %d, want %d", luaInteger(co, -1), want)
			}
			pop(co, 1)
			continue
		}
		if status != lua.OK {
			t.Fatalf("resume status = %d, error=%s", status, luaString(co, -1))
		}
		break
	}
	if co.Status() != lua.OK {
		t.Fatalf("final status = %d", co.Status())
	}
}

func coroutineCFunction(L *lua.State) c.Int {
	L.Yieldk(0, 0, nil)
	return 0
}

func TestCoroutineCFunction(t *testing.T) {
	L := newState(t)
	co := L.Newthread()
	co.Pushcclosure(coroutineCFunction, 0)
	var nres c.Int
	if status := co.Resume(nil, 0, &nres); status != lua.YIELD {
		t.Fatalf("first resume = %d", status)
	}
	if status := co.Resume(nil, 0, &nres); status != lua.OK {
		t.Fatalf("second resume = %d", status)
	}
}

func countdownContinue(L *lua.State, status c.Int, ctx lua.KContext) c.Int { return countdown(L) }

func countdown(L *lua.State) c.Int {
	arg := luaInteger(L, upvalueIndex(1))
	if arg > 0 {
		L.Pushinteger(arg - 1)
		L.Copy(-1, upvalueIndex(1))
		pop(L, 1)
		L.Pushinteger(arg)
		return L.Yieldk(1, 0, countdownContinue)
	}
	return 0
}

func createCountdown(L *lua.State) c.Int {
	max := L.Checkinteger(1)
	L.Pushinteger(max)
	L.Pushcclosure(countdown, 1)
	return 1
}

func TestCoroutineContinue(t *testing.T) {
	L := newState(t)
	setGlobalFunction(L, "create_countdown", createCountdown)
	code := `local countdown = create_countdown(5)
local co = coroutine.create(countdown)
local values = {}
while true do
  local success, value = coroutine.resume(co)
  if not success then return value end
  if value == nil then break end
  values[#values + 1] = value
end
return table.concat(values, ",")`
	if status := doString(L, code); status != lua.OK {
		t.Fatalf("run: %s", luaString(L, -1))
	}
	if got := luaString(L, -1); got != "5,4,3,2,1" {
		t.Fatalf("countdown = %q", got)
	}
}

func triggerError(L *lua.State) c.Int {
	L.Pushstring(c.Str("This is an error triggered"))
	return L.Error()
}

func triggerFormatError(L *lua.State) c.Int {
	return L.LError(c.Str("This is an error code:(%d)"), c.Int(42))
}

func TestCustomPanic(t *testing.T) {
	L := newState(t)
	setGlobalFunction(L, "trigger_error", triggerError)
	setGlobalFunction(L, "trigger_format_error", triggerFormatError)
	for _, name := range []string{"trigger_error", "trigger_format_error"} {
		L.Getglobal(c.AllocaCStr(name))
		if status := pcall(L, 0, 0); status != lua.ERRRUN {
			t.Fatalf("%s status = %d", name, status)
		}
		if luaString(L, -1) == "" {
			t.Fatalf("%s error is empty", name)
		}
		pop(L, 1)
	}
}

var debugHookCalls c.Int

func debugHook(L *lua.State, ar *lua.Debug) {
	debugHookCalls++
	L.Getinfo(c.Str("nSl"), ar)
}

func TestDebugHook(t *testing.T) {
	L := newState(t)
	debugHookCalls = 0
	L.Sethook(debugHook, lua.MASKLINE, 0)
	if status := doString(L, `local x = 1
x = x + 1`); status != lua.OK {
		t.Fatalf("run: %s", luaString(L, -1))
	}
	L.Sethook(nil, 0, 0)
	if debugHookCalls == 0 {
		t.Fatal("debug hook was not called")
	}
}

func extraspace(L *lua.State) *int {
	return (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(L)) - unsafe.Sizeof(uintptr(0))))
}

func extraSpaceFunction(L *lua.State) c.Int {
	L.Pushfstring(c.Str("Stored integer is: %d"), *extraspace(L))
	return 1
}

func TestExtraspace(t *testing.T) {
	L := newState(t)
	*extraspace(L) = 42
	setGlobalFunction(L, "GetData", extraSpaceFunction)
	if status := doString(L, `return GetData()`); status != lua.OK {
		t.Fatalf("run: %s", luaString(L, -1))
	}
	if got := luaString(L, -1); got != "Stored integer is: 42" {
		t.Fatalf("extra space = %q", got)
	}
}

func testAlloc(ud unsafe.Pointer, ptr unsafe.Pointer, osize c.SizeT, nsize c.SizeT) unsafe.Pointer {
	if nsize == 0 {
		c.Free(ptr)
		return nil
	}
	return c.Realloc(ptr, nsize)
}

func TestStateAlloc(t *testing.T) {
	L := lua.Newstate(testAlloc, nil, 0)
	if L == nil {
		t.Fatal("lua_newstate returned nil")
	}
	t.Cleanup(L.Close)
	L.Openselectedlibs(c.Int(-1), c.Int(0))
	if status := doString(L, `return 1 + 1`); status != lua.OK || luaInteger(L, -1) != 2 {
		t.Fatalf("custom allocator result: %s", luaString(L, -1))
	}
	allocf := L.Getallocf(nil)
	L.Setallocf(allocf, nil)
	L.Settop(0)
	if status := doString(L, `return "set newstate success"`); status != lua.OK || luaString(L, -1) != "set newstate success" {
		t.Fatal("Setallocf failed")
	}
}

var dumpedChunk []byte
var dumpOffset int

func dumpWriter(L *lua.State, p unsafe.Pointer, sz c.SizeT, ud unsafe.Pointer) c.Int {
	dumpedChunk = append(dumpedChunk, unsafe.Slice((*byte)(p), int(sz))...)
	return lua.OK
}

func dumpReader(L *lua.State, ud unsafe.Pointer, size *c.SizeT) *c.Char {
	if dumpOffset >= len(dumpedChunk) {
		*size = 0
		return nil
	}
	part := dumpedChunk[dumpOffset:]
	dumpOffset = len(dumpedChunk)
	*size = c.SizeT(len(part))
	return (*c.Char)(unsafe.Pointer(unsafe.SliceData(part)))
}

func makeDump(t *testing.T) []byte {
	L := newState(t)
	if status := L.Loadstring(c.Str(`return function(name) return "Hello, " .. name .. "!" end`)); status != lua.OK {
		t.Fatalf("load: %s", luaString(L, -1))
	}
	if status := pcall(L, 0, 1); status != lua.OK || L.Type(-1) != lua.TFUNCTION {
		t.Fatalf("function creation failed: %s", luaString(L, -1))
	}
	dumpedChunk = nil
	if status := L.Dump(dumpWriter, nil, 0); status != lua.OK || len(dumpedChunk) == 0 {
		t.Fatalf("dump failed: status=%d bytes=%d", status, len(dumpedChunk))
	}
	return append([]byte(nil), dumpedChunk...)
}

func TestDumpSave(t *testing.T) {
	if len(makeDump(t)) == 0 {
		t.Fatal("dump produced no bytes")
	}
}

func TestDumpLoad(t *testing.T) {
	L := newState(t)
	dumpedChunk = makeDump(t)
	dumpOffset = 0
	if status := L.Load(dumpReader, nil, c.Str("greet"), nil); status != lua.OK {
		t.Fatalf("load dumped chunk: %s", luaString(L, -1))
	}
	L.Pushstring(c.Str("World"))
	if status := pcall(L, 1, 1); status != lua.OK {
		t.Fatalf("call dumped function: %s", luaString(L, -1))
	}
	if got := luaString(L, -1); got != "Hello, World!" {
		t.Fatalf("dumped result = %q", got)
	}
}

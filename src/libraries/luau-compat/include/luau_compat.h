/**
 * Compatibility shims so LÖVE's Lua 5.1-style C code works against Luau.
 * Include AFTER lua.h / lualib.h / luacode.h.
 **/
#ifndef LOVE_LUAU_COMPAT_H
#define LOVE_LUAU_COMPAT_H

#include <stddef.h>
#include <string.h>

/* Ensure Lua headers were included first. */
#ifndef LUA_REGISTRYINDEX
#error "luau_compat.h requires Luau lua.h (include lua.h/lualib.h/luacode.h before this header)"
#endif

#ifndef LUA_VERSION_NUM
#define LUA_VERSION_NUM 501
#endif

#define LOVE_LUAU 1

#ifndef lua_assert
#include <assert.h>
#define lua_assert(x) assert(x)
#endif

/* Luau's lua_pushcfunction / lua_pushcclosure require a debug name. */
#undef lua_pushcfunction
#undef lua_pushcclosure
#define lua_pushcfunction(L, fn) lua_pushcclosurek((L), (fn), NULL, 0, NULL)
#define lua_pushcclosure(L, fn, nup) lua_pushcclosurek((L), (fn), NULL, (nup), NULL)

/* Match Lua 5.1 reference sentinel values for table-based luaL_ref. */
#undef LUA_NOREF
#undef LUA_REFNIL
#define LUA_NOREF (-2)
#define LUA_REFNIL (-1)

/* Make noreturn error helpers usable in `return luaL_error(...)` / `return lua_error(L)`.
 * Luau declares these as void; Lua 5.1 returned int. */
static inline int love_compat_lua_error(lua_State *L)
{
	lua_error(L);
	return 0;
}
#define lua_error(L) love_compat_lua_error(L)

#undef luaL_error
#undef luaL_typeerror
#undef luaL_argerror
#define luaL_error(L, fmt, ...) (luaL_errorL((L), (fmt), ##__VA_ARGS__), 0)
#define luaL_typeerror(L, narg, tname) (luaL_typeerrorL((L), (narg), (tname)), 0)
#define luaL_argerror(L, narg, extramsg) (luaL_argerrorL((L), (narg), (extramsg)), 0)

#define luaL_checkint(L, i) ((int)luaL_checkinteger((L), (i)))
#define luaL_optint(L, i, d) ((int)luaL_optinteger((L), (i), (d)))
#define luaL_checklong(L, i) ((long)luaL_checkinteger((L), (i)))
#define luaL_optlong(L, i, d) ((long)luaL_optinteger((L), (i), (d)))

#define luaL_getn(L, i) ((int)lua_objlen((L), (i)))
#define luaL_setn(L, i, j) ((void)0)

#ifndef LUAL_BUFFERSIZE
#define LUAL_BUFFERSIZE LUA_BUFFERSIZE
#endif

#ifndef luaL_addsize
#define luaL_addsize(B, s) ((void)((B)->p += (s)))
#endif

#ifndef luaL_prepbuffer
#define luaL_prepbuffer(B) luaL_prepbuffsize((B), LUAL_BUFFERSIZE)
#endif

#define luaL_reg luaL_Reg

#ifndef lua_register
#define lua_register(L, n, f) (lua_pushcfunction((L), (f)), lua_setglobal((L), (n)))
#endif

/* Provide a classic lauxlib.h include path for third-party C libraries. */
#ifndef LOVE_LUAU_LAUXLIB_GUARD
#define LOVE_LUAU_LAUXLIB_GUARD
#endif

#ifdef __cplusplus
extern "C" {
#endif

#if defined(__GNUC__) || defined(__clang__)
#define LOVE_LUAU_EXPORT __attribute__((visibility("default")))
#else
#define LOVE_LUAU_EXPORT
#endif

LOVE_LUAU_EXPORT int luaL_ref(lua_State *L, int t);
LOVE_LUAU_EXPORT void luaL_unref(lua_State *L, int t, int ref);

/* Compile Luau source and load bytecode. Returns LUA_OK or an error code. */
LOVE_LUAU_EXPORT int luaL_loadbuffer(lua_State *L, const char *buff, size_t sz, const char *name);
LOVE_LUAU_EXPORT int luaL_loadbufferx(lua_State *L, const char *buff, size_t sz, const char *name, const char *mode);
LOVE_LUAU_EXPORT int luaL_loadstring(lua_State *L, const char *s);
LOVE_LUAU_EXPORT int luaL_loadfile(lua_State *L, const char *filename);

/* Minimal package / require (Luau does not ship package). */
LOVE_LUAU_EXPORT int luaopen_love_package(lua_State *L);
LOVE_LUAU_EXPORT void love_open_package(lua_State *L);

#ifdef __cplusplus
}
#endif

#endif /* LOVE_LUAU_COMPAT_H */

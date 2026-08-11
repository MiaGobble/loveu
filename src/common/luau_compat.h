/**
 * Compatibility shims so LÖVE's Lua 5.1-style C code works against Luau.
 * Include after lua.h / lualib.h / luacode.h.
 **/
#ifndef LOVE_LUAU_COMPAT_H
#define LOVE_LUAU_COMPAT_H

#include <stddef.h>
#include <string.h>

#ifndef LUA_VERSION_NUM
#define LUA_VERSION_NUM 501
#endif

#define LOVE_LUAU 1

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

/* Make noreturn error helpers usable in `return luaL_error(...)`. */
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

/* Provide a classic lauxlib.h include path for third-party C libraries. */
#ifndef LOVE_LUAU_LAUXLIB_GUARD
#define LOVE_LUAU_LAUXLIB_GUARD
#endif

#ifdef __cplusplus
extern "C" {
#endif

int luaL_ref(lua_State *L, int t);
void luaL_unref(lua_State *L, int t, int ref);

/* Compile Luau source and load bytecode. Returns LUA_OK or an error code. */
int luaL_loadbuffer(lua_State *L, const char *buff, size_t sz, const char *name);
int luaL_loadbufferx(lua_State *L, const char *buff, size_t sz, const char *name, const char *mode);
int luaL_loadstring(lua_State *L, const char *s);
int luaL_loadfile(lua_State *L, const char *filename);

/* Minimal package / require (Luau does not ship package). */
int luaopen_love_package(lua_State *L);
void love_open_package(lua_State *L);

#ifdef __cplusplus
}
#endif

#endif /* LOVE_LUAU_COMPAT_H */

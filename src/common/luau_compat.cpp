/**
 * Luau compatibility helpers: loadbuffer, luaL_ref, and package/require.
 **/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#ifdef __cplusplus
extern "C" {
#endif
#include <lua.h>
#include <lualib.h>
#include <luacode.h>
#ifdef __cplusplus
}
#endif
#include "luau_compat.h"

#ifndef LUA_ERRFILE
#define LUA_ERRFILE LUA_ERRRUN
#endif

#ifndef LUA_PATHSEP
#define LUA_PATHSEP ";"
#endif
#ifndef LUA_PATH_MARK
#define LUA_PATH_MARK "?"
#endif
#ifndef LUA_DIRSEP
#define LUA_DIRSEP "/"
#endif

extern "C" {

/* ---- luaL_ref / luaL_unref (Lua 5.1 table freelist semantics) ---- */

#define FREELIST_REF 0

static int absindex(lua_State *L, int idx)
{
	if (idx > 0 || idx <= LUA_REGISTRYINDEX)
		return idx;
	return lua_gettop(L) + idx + 1;
}

int luaL_ref(lua_State *L, int t)
{
	int ref;
	t = absindex(L, t);
	if (lua_isnil(L, -1))
	{
		lua_pop(L, 1);
		return LUA_REFNIL;
	}
	lua_rawgeti(L, t, FREELIST_REF);
	ref = (int)lua_tointeger(L, -1);
	lua_pop(L, 1);
	if (ref != 0)
	{
		lua_rawgeti(L, t, ref);
		lua_rawseti(L, t, FREELIST_REF);
	}
	else
	{
		ref = (int)lua_objlen(L, t) + 1;
	}
	lua_rawseti(L, t, ref);
	return ref;
}

void luaL_unref(lua_State *L, int t, int ref)
{
	if (ref >= 0)
	{
		t = absindex(L, t);
		lua_rawgeti(L, t, FREELIST_REF);
		lua_rawseti(L, t, ref);
		lua_pushinteger(L, ref);
		lua_rawseti(L, t, FREELIST_REF);
	}
}

/* ---- load helpers ---- */

int luaL_loadbufferx(lua_State *L, const char *buff, size_t sz, const char *name, const char *mode)
{
	(void)mode;
	if (buff == NULL)
	{
		lua_pushstring(L, "cannot load nil buffer");
		return LUA_ERRSYNTAX;
	}

	/* Reject classic Lua binary chunks; Luau only loads its own bytecode. */
	if (sz > 0 && (unsigned char)buff[0] == 0x1b)
	{
		lua_pushstring(L, "binary Lua chunks are not supported (Luau)");
		return LUA_ERRSYNTAX;
	}

	size_t bytecodeSize = 0;
	char *bytecode = luau_compile(buff, sz, NULL, &bytecodeSize);
	if (bytecode == NULL)
	{
		lua_pushstring(L, "Luau compile failed (out of memory)");
		return LUA_ERRMEM;
	}

	/* luau_compile encodes errors as a bytecode buffer starting with \0. */
	if (bytecodeSize > 0 && bytecode[0] == '\0')
	{
		lua_pushlstring(L, bytecode + 1, bytecodeSize - 1);
		free(bytecode);
		return LUA_ERRSYNTAX;
	}

	const char *chunkname = name ? name : "=(load)";
	if (chunkname[0] == '@' || chunkname[0] == '=')
		chunkname++;

	int result = luau_load(L, chunkname, bytecode, bytecodeSize, 0);
	free(bytecode);
	return result == 0 ? LUA_OK : LUA_ERRSYNTAX;
}

int luaL_loadbuffer(lua_State *L, const char *buff, size_t sz, const char *name)
{
	return luaL_loadbufferx(L, buff, sz, name, NULL);
}

int luaL_loadstring(lua_State *L, const char *s)
{
	return luaL_loadbuffer(L, s, strlen(s), s);
}

int luaL_loadfile(lua_State *L, const char *filename)
{
	FILE *f = fopen(filename, "rb");
	if (!f)
	{
		lua_pushfstring(L, "cannot open %s", filename);
		return LUA_ERRFILE;
	}
	if (fseek(f, 0, SEEK_END) != 0)
	{
		fclose(f);
		lua_pushfstring(L, "cannot read %s", filename);
		return LUA_ERRFILE;
	}
	long len = ftell(f);
	fseek(f, 0, SEEK_SET);
	if (len < 0)
	{
		fclose(f);
		lua_pushfstring(L, "cannot read %s", filename);
		return LUA_ERRFILE;
	}
	char *buf = (char *)malloc((size_t)len + 1);
	if (!buf)
	{
		fclose(f);
		lua_pushstring(L, "out of memory");
		return LUA_ERRMEM;
	}
	size_t n = fread(buf, 1, (size_t)len, f);
	fclose(f);
	buf[n] = '\0';
	char chunkname[512];
	snprintf(chunkname, sizeof(chunkname), "@%s", filename);
	int status = luaL_loadbuffer(L, buf, n, chunkname);
	free(buf);
	return status;
}

/* ---- package / require ---- */

static const char *love_gsub(lua_State *L, const char *s, const char *p, const char *r)
{
	const char *wild;
	size_t l = strlen(p);
	luaL_Buffer b;
	luaL_buffinit(L, &b);
	while ((wild = strstr(s, p)) != NULL)
	{
		luaL_addlstring(&b, s, (size_t)(wild - s));
		luaL_addstring(&b, r);
		s = wild + l;
	}
	luaL_addstring(&b, s);
	luaL_pushresult(&b);
	return lua_tostring(L, -1);
}

static int readable(const char *filename)
{
	FILE *f = fopen(filename, "r");
	if (!f)
		return 0;
	fclose(f);
	return 1;
}

static const char *pushnexttemplate(lua_State *L, const char *path)
{
	const char *sep = LUA_PATHSEP;
	while (*path == *sep)
		path++;
	if (*path == '\0')
		return NULL;
	const char *end = strchr(path, *sep);
	if (end == NULL)
		end = path + strlen(path);
	lua_pushlstring(L, path, (size_t)(end - path));
	return end;
}

static const char *searchpath(lua_State *L, const char *name, const char *path)
{
	luaL_Buffer msg;
	luaL_buffinit(L, &msg);
	love_gsub(L, name, ".", LUA_DIRSEP);
	/* name with dots replaced is on stack */
	while ((path = pushnexttemplate(L, path)) != NULL)
	{
		const char *filename = love_gsub(L, lua_tostring(L, -1), LUA_PATH_MARK, lua_tostring(L, -2));
		lua_remove(L, -2); /* template */
		if (readable(filename))
		{
			lua_remove(L, -2); /* transformed name */
			return lua_tostring(L, -1);
		}
		lua_pushfstring(L, "\n\tno file '%s'", filename);
		lua_remove(L, -2); /* filename */
		luaL_addvalue(&msg);
	}
	lua_pop(L, 1); /* transformed name */
	luaL_pushresult(&msg);
	return NULL;
}

static int loader_preload(lua_State *L)
{
	const char *name = luaL_checkstring(L, 1);
	lua_getfield(L, lua_upvalueindex(1), "preload");
	if (!lua_istable(L, -1))
		luaL_error(L, "'package.preload' must be a table");
	lua_getfield(L, -1, name);
	if (lua_isnil(L, -1))
		lua_pushfstring(L, "\n\tno field package.preload['%s']", name);
	return 1;
}

static int loader_Lua(lua_State *L)
{
	const char *name = luaL_checkstring(L, 1);
	lua_getfield(L, lua_upvalueindex(1), "path");
	const char *path = lua_tostring(L, -1);
	if (path == NULL)
		luaL_error(L, "'package.path' must be a string");
	const char *filename = searchpath(L, name, path);
	if (filename == NULL)
		return 1; /* error string already on stack */
	if (luaL_loadfile(L, filename) == LUA_OK)
		return 1;
	lua_pushfstring(L, "\n\terror loading module '%s' from file '%s':\n\t%s",
					name, filename, lua_tostring(L, -1));
	return 1;
}

static int ll_require(lua_State *L)
{
	const char *name = luaL_checkstring(L, 1);
	lua_settop(L, 1);
	lua_getfield(L, LUA_REGISTRYINDEX, "_LOADED");
	lua_getfield(L, 2, name);
	if (lua_toboolean(L, -1))
		return 1;

	lua_getfield(L, lua_upvalueindex(1), "loaders");
	if (!lua_istable(L, -1))
		luaL_error(L, "'package.loaders' must be a table");
	lua_pushliteral(L, "");
	for (int i = 1;; i++)
	{
		lua_rawgeti(L, -2, i);
		if (lua_isnil(L, -1))
		{
			luaL_error(L, "module '%s' not found:%s", name, lua_tostring(L, -2));
			return 0;
		}
		lua_pushstring(L, name);
		lua_call(L, 1, 1);
		if (lua_isfunction(L, -1))
			break;
		if (lua_isstring(L, -1))
			lua_concat(L, 2);
		else
			lua_pop(L, 1);
	}
	lua_pushstring(L, name);
	lua_call(L, 1, 1);
	if (!lua_isnil(L, -1))
		lua_setfield(L, 2, name);
	lua_getfield(L, 2, name);
	if (lua_isnil(L, -1))
	{
		lua_pushboolean(L, 1);
		lua_pushvalue(L, -1);
		lua_setfield(L, 2, name);
	}
	return 1;
}

static int w_loadstring(lua_State *L)
{
	size_t len = 0;
	const char *s = luaL_optlstring(L, 1, "", &len);
	const char *chunkname = luaL_optstring(L, 2, s);
	if (luaL_loadbuffer(L, s, len, chunkname) == LUA_OK)
		return 1;
	lua_pushnil(L);
	lua_insert(L, -2);
	return 2;
}

static int w_load(lua_State *L)
{
	size_t len = 0;
	const char *s = luaL_checklstring(L, 1, &len);
	const char *chunkname = luaL_optstring(L, 2, "=(load)");
	if (luaL_loadbuffer(L, s, len, chunkname) == LUA_OK)
		return 1;
	lua_pushnil(L);
	lua_insert(L, -2);
	return 2;
}

/* Alias bit32 as bit for LuaJIT BitOp API compatibility. */
static int open_bit_alias(lua_State *L)
{
	lua_getglobal(L, "bit32");
	if (lua_istable(L, -1))
	{
		lua_pushvalue(L, -1);
		lua_setglobal(L, "bit");
	}
	lua_pop(L, 1);
	return 0;
}

int luaopen_love_package(lua_State *L)
{
	lua_getfield(L, LUA_REGISTRYINDEX, "_LOADED");
	if (!lua_istable(L, -1))
	{
		lua_pop(L, 1);
		lua_newtable(L);
		lua_pushvalue(L, -1);
		lua_setfield(L, LUA_REGISTRYINDEX, "_LOADED");
	}
	lua_pop(L, 1);

	lua_newtable(L); /* package */

	lua_pushvalue(L, -1);
	lua_pushcclosure(L, ll_require, 1);
	lua_setglobal(L, "require");

	lua_newtable(L);
	lua_setfield(L, -2, "preload");

	lua_newtable(L); /* loaders */
	lua_pushvalue(L, -2); /* package as upvalue */
	lua_pushcclosure(L, loader_preload, 1);
	lua_rawseti(L, -2, 1);
	lua_pushvalue(L, -2);
	lua_pushcclosure(L, loader_Lua, 1);
	lua_rawseti(L, -2, 2);
	lua_setfield(L, -2, "loaders");
	lua_getfield(L, -1, "loaders");
	lua_setfield(L, -2, "searchers"); /* Lua 5.2+ alias */

	lua_pushstring(L, "?.luau;?/init.luau");
	lua_setfield(L, -2, "path");

	lua_pushstring(L, "");
	lua_setfield(L, -2, "cpath");

	lua_getfield(L, LUA_REGISTRYINDEX, "_LOADED");
	lua_setfield(L, -2, "loaded");

	lua_setglobal(L, "package");

	lua_pushcfunction(L, w_loadstring);
	lua_setglobal(L, "loadstring");

	lua_pushcfunction(L, w_load);
	lua_setglobal(L, "load");

	open_bit_alias(L);
	return 0;
}

void love_open_package(lua_State *L)
{
	luaopen_love_package(L);
}

} /* extern "C" */

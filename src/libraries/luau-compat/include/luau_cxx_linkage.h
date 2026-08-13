/**
 * Force C linkage on Luau API macros when included from C++.
 * Injected via lovedep_Lua compile options for CXX only so .c files keep
 * LUA_API=extern (valid C) while still linking to Luau built with LUAU_EXTERN_C.
 *
 * Only define LUA_API / LUACODE_API; luaconf.h sets LUALIB_API from LUA_API.
 **/
#pragma once

#ifdef __cplusplus
#	ifndef LUA_API
#		define LUA_API extern "C"
#	endif
#	ifndef LUACODE_API
#		define LUACODE_API extern "C"
#	endif
#endif

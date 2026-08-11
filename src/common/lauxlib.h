/**
 * Drop-in lauxlib.h for Luau embeds in loveu.
 * Luau ships auxiliares in lualib.h; this header restores the classic include path
 * and Lua 5.1-oriented helpers used by LÖVE and third-party bindings.
 **/
#pragma once

#include "lua.h"
#include "lualib.h"
#include "luacode.h"
#include "luau_compat.h"

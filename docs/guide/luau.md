# Luau scripting

loveu embeds [Luau](https://luau.org/) instead of LuaJIT / PUC Lua. Language details live on the Luau site. This page is only what matters for games.

## Files

- Entry point: **`main.luau`** (not `main.lua`).
- Modules: `require` looks for `?.luau` and `?/init.luau` under `code_root`.
- Engine-internal scripts still compile through Luau at load time.

Callbacks are the same as LÖVE: `love.load`, `love.update`, `love.draw`, and the rest. See [love.callbacks](https://love2d.org/wiki/love#Callbacks) on the wiki.

```luau
function love.load()
end

function love.update(dt)
end

function love.draw()
	love.graphics.print("Hello from Luau!", 100, 100)
end
```

## What is missing vs LuaJIT LÖVE

- No FFI (`require("ffi")`).
- No `jit.*`.
- No `io` library (`io.open` is nil). Use [love.filesystem](https://love2d.org/wiki/love.filesystem).
- Prefer `bit32` for bitwise ops. loveu also aliases it as `bit`.

## Types

Luau can typecheck, but **`love` type definitions for this fork are not implemented yet**. You can still write untyped `.luau` like the examples.

## Requires

See [Requires](/guide/requires) for `./`, `../`, and bare module names.

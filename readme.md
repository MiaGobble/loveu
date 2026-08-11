# loveu
[![CI](https://github.com/MiaGobble/loveu/actions/workflows/main.yml/badge.svg)](https://github.com/MiaGobble/loveu/actions/workflows/main.yml)

loveu is a fork of LÖVE to support different things that I want out of the engine.

## Divergence
Refer to the below table to see what changed.

|Feature|Description|Status|
|-------|-----------|------|
|Relative requires|Make requiring relative to the current file path, not the root|Not implemented|
|Luau support|Usage of Luau instead of Lua|Implemented|
|CLI tooling|Helpers to build, run, initialize, and more|Not implemented|
|Improved mobile building|Easily build for mobile, instead of using Android or iOS repositories|Not implemented|
|Native UI support|Adds support for UI|Not implemented|
|Whale package manager|A package manager for LÖVE and loveu|Not implemented|
|Audio improvements|Support for audio buses, global mixer channels, and dynamic sound effects|Not implemented|

## Luau scripting

loveu embeds [Luau](https://luau.org/) instead of LuaJIT / PUC Lua.

- Game entry point must be **`main.luau`** (not `main.lua`).
- Config file must be **`conf.luau`** when used.
- `require` resolves `?.luau` and `?/init.luau` only.
- LuaJIT FFI and `jit.*` are not available; prefer `bit32` (also aliased as `bit`).
- Engine-internal scripts still compile through Luau at load time.

Example `main.luau`:

```luau
function love.draw()
	love.graphics.print("Hello from Luau!", 100, 100)
end
```

Luau sources live in `src/libraries/luau` when vendored. If that directory is missing, CMake fetches Luau via FetchContent at configure time.

## AI Notice
This fork, at least so far, is completely done with AI. This means you can expect things to not be exceptionally implemented, or alternatively, not working the way it should.

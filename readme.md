# loveu
[![CI](https://github.com/MiaGobble/loveu/actions/workflows/main.yml/badge.svg)](https://github.com/MiaGobble/loveu/actions/workflows/main.yml)

loveu is a fork of LÖVE to support different things that I want out of the engine.

## Roadmap
I will be doing the following divergences in order:
* Foundation
* Development Features
* Content Features
* Shipping & LiveOps Features

### Foundation
|Feature|Description|Status|
|-------|-----------|------|
|Relative requires|Make requiring relative to the current file path, not the root|Implemented|
|Luau support|Usage of Luau instead of Lua|Implemented|
|Project format|Structures the project identity from a `loveu.toml` manifest|Implemented|
|CLI tooling|Helpers to build to any platform, run, initialize, switch versions, and more|Not implemented|
|`love` type definitions|Fully featured typechecking for `love`|Not implemented|
|Documentation|Documentation of this fork to note divergences, with a link back to the Love2D documentation|Not implemented|

### Development Features
|Feature|Description|Status|
|-------|-----------|------|
|Scenes|In-engine features for scenes, with transition hooks|Not implemented|
|Action-based input|Allows creation of actions that can be connected to based on an input map|Not implemented|
|Hot reload|When you save a script or texture, the running game updates without a full restart|Not implemented|
|Resource lifetimes|Safely manage assets without extra boilerplate|Not implemented|
|Improved crashing and logging|Structured logs with crash dumps and optional opt-in telemetry hooks|Not implemented|
|In-engine debugger|A toggleable overlay for showing FPS, frame time, draw calls, texture memory, current scene, and more|Not implemented|
|Fixed timestep|Ability for fixed timestepping|Not implemented|

### Content Features
|Feature|Description|Status|
|-------|-----------|------|
|Game UI support|Adds support for UI tooling like layout, focus, text, and more|Not implemented|
|Audio improvements|Support for audio buses, global mixer channels, and dynamic sound effects|Not implemented|
|Tilemaps|An easy way to create 2D worlds|Not implemented|
|Save API|Versioned, atomic saving for progress and settings|Not implemented|

### Shipping & LiveOps Features
|Feature|Description|Status|
|-------|-----------|------|
|Platform services|Services related to different platforms like Steam, Google Play, and more|Not implemented|
|Analytics|Support for third-party analytics SDKs|Not implemented|


## Luau scripting
loveu embeds [Luau](https://luau.org/) instead of LuaJIT / PUC Lua.

- Game entry point must be **`main.luau`** (not `main.lua`).
- Config file must be **`conf.luau`** when used.
- Every game must include a root **`loveu.toml`** project manifest (missing file is a boot error).
- `require` resolves `?.luau` and `?/init.luau` only (under `code_root` when set).
- Paths starting with `./` or `../` resolve relative to the requiring file (e.g. `require("../../lib/util")` from `scenes/level/main.luau` loads `lib/util.luau`). Bare names like `require("lib.util")` remain root-relative within `code_root`.
- LuaJIT FFI and `jit.*` are not available; prefer `bit32` (also aliased as `bit`).
- Engine-internal scripts still compile through Luau at load time.

### Project format (`loveu.toml`)

Required at the game root (next to or above `code_root`):

```toml
name = "mygame"
version = "0.1.0"
engine_version = "0.1.0"
code_root = "src"
```

| Field | Meaning |
|-------|---------|
| `name` | Project identity (save directory name; default window title) |
| `version` | Game version (`love.project.version`) |
| `engine_version` | Must exactly match `love._loveu_version` or boot fails |
| `code_root` | Scripts directory relative to the manifest (use `"."` when `main.luau` is beside the toml) |

`conf.luau` still configures window/modules/runtime options. Identity comes from `loveu.toml`, not `t.identity`. Upstream LÖVE API compat (`love.conf` `t.version`, `love.getVersion`) still uses `love._version` (currently `12.0`).

Exposed at runtime as `love.project = { name, version, engine_version, code_root }`.

### Versioning

| Identifier | Meaning | Current |
|------------|---------|---------|
| `love._loveu_version` | loveu fork semver; pin in `loveu.toml` | `0.1.0` |
| `love._version` | Upstream LÖVE API this build is based on | `12.0` |

Bump `LOVEU_VERSION_*` in `src/common/version.h` when loveu changes behavior games depend on. Bump `LOVE_VERSION_*` when rebasing onto a new upstream LÖVE. `love --version` prints both.

Example `main.luau`:

```luau
function love.draw()
	love.graphics.print("Hello from Luau!", 100, 100)
end
```

Luau sources live in `src/libraries/luau` when vendored. If that directory is missing, CMake fetches Luau via FetchContent at configure time.

## AI Notice
This fork, at least so far, is completely done with AI. This means you can expect things to not be exceptionally implemented, or alternatively, not working the way it should.

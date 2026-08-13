# Divergences from LÖVE

Everything not listed here is intended to match [LÖVE 12](https://love2d.org/wiki/love). When in doubt, start with the [LÖVE wiki](https://love2d.org/wiki/Main_Page).

## Language and files

| LÖVE | loveu |
|------|--------|
| LuaJIT / PUC Lua, `*.lua` | [Luau](https://luau.org/), `*.luau` |
| `main.lua` | `main.luau` |
| Optional `conf.lua` for identity, window, modules | Required `loveu.toml`. **No `conf.luau`.** |
| `require` uses `?.lua` / `?/init.lua` | `?.luau` / `?/init.luau` only |
| `require` is root-relative | `./` and `../` are relative to the **calling file** |
| LuaJIT FFI, `jit.*` | Not available. Use `bit32` (also aliased as `bit`) |

## Boot and identity

| LÖVE | loveu |
|------|--------|
| Identity from folder name / `t.identity` | Identity from `loveu.toml` `name` |
| Optional `t.version` in `love.conf` | `engine_version` must **exactly** match `love._loveu_version` or boot fails |
| Missing `main.lua` is “no code” | Missing `loveu.toml` is a **hard boot error** |
| `love --version` prints LÖVE | Prints `loveu … (LÖVE …)` |

## Runtime extras

| API | Notes |
|-----|--------|
| [`love.project`](/api/project) | Manifest fields plus parsed `config` |
| [`love._loveu_version`](/api/version) | Fork semver; pin this in `loveu.toml` |
| `love._version` / `love.getVersion()` | Still the upstream LÖVE 12 API version |

## Not in loveu (yet)

From the project roadmap: CLI tooling, `love` type definitions, scenes, action-based input, hot reload, and the later content / live-ops items. Those are not documented here because they are not implemented.

## Standard library

Luau does **not** ship `io`. Use [`love.filesystem`](https://love2d.org/wiki/love.filesystem) (`read`, `write`, `openFile`, `mountFullPath`, …). Luau `os` is limited (no `os.rename` / `os.remove`).

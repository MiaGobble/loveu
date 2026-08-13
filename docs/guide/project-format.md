# Project format

Every game must have a **`loveu.toml`** at the project root (next to `code_root`, not inside it unless `code_root` is `"."`).

Missing or invalid TOML is a hard boot error. `engine_version` must exactly equal `love._loveu_version`.

There is **no `conf.luau`**. Window, modules, and the other old `love.conf` fields belong in this file.

## Required fields

```toml
name = "mygame"
version = "0.1.0"
engine_version = "0.1.0"
code_root = "src"
```

| Field | Meaning |
|-------|---------|
| `name` | Save-directory identity and default window title |
| `version` | Game version (`love.project.version`) |
| `engine_version` | Must match running loveu (`love._loveu_version`) |
| `code_root` | Scripts directory relative to the manifest. Use `"."` when `main.luau` sits beside the toml |

`code_root` must be a relative path. `""` is treated as `"."`. `..` and absolute paths are rejected.

After boot, see [`love.project`](/api/project).

## Optional runtime config

These keys are merged into the engine’s default config (same defaults LÖVE used in `love.conf`).

### Top-level

| Key | Type | Default | Notes |
|-----|------|---------|--------|
| `title` | string | `name` | Window title if not `"Untitled"` |
| `console` | bool | `false` | Windows console |
| `appendidentity` | bool | `false` | |
| `externalstorage` | bool | `false` | Android |
| `highdpi` | bool | `false` | |
| `trackpadtouch` | bool | `false` | |

### `[window]`

| Key | Default |
|-----|---------|
| `width` | `800` |
| `height` | `600` |
| `x` / `y` | unset |
| `minwidth` / `minheight` | `1` |
| `fullscreen` | `false` |
| `fullscreentype` | `"desktop"` |
| `displayindex` | `1` |
| `vsync` | `1` |
| `msaa` | `0` |
| `borderless` | `false` |
| `resizable` | `false` |
| `centered` | `true` |
| `usedpiscale` | `true` |
| `stencil` / `depth` | unset |
| `icon` | unset (path under the game source) |

Set `[modules] window = false` (and usually `graphics = false`) for a headless run. That replaces the old `t.window = false` pattern.

### `[modules]`

All default to `true`: `data`, `event`, `keyboard`, `mouse`, `timer`, `joystick`, `touch`, `image`, `graphics`, `audio`, `math`, `physics`, `sensor`, `sound`, `system`, `font`, `thread`, `window`, `video`.

### `[graphics]`

| Key | Default |
|-----|---------|
| `gammacorrect` | `false` |
| `lowpower` | `false` |
| `renderers` | engine default |
| `excluderenderers` | unset |

`--renderers` / `--excluderenderers` on the command line still override this. See [love.conf](https://love2d.org/wiki/Config_Files) on the wiki for the LÖVE meanings; the keys are the same, the file is not.

### `[audio]`

| Key | Default |
|-----|---------|
| `mixwithsystem` | `true` |
| `mic` | `false` |

## Example

```toml
name = "mygame"
version = "0.1.0"
engine_version = "0.1.0"
code_root = "src"
console = false

[window]
width = 1280
height = 720
resizable = true
vsync = 1

[modules]
physics = false
video = false
```

## Errors

| Situation | Result |
|-----------|--------|
| No `loveu.toml` | Boot error: missing required manifest |
| Invalid TOML | Boot error |
| Missing `name` / `version` / `engine_version` / `code_root` | Boot error |
| `engine_version` ≠ `love._loveu_version` | Boot error |
| `code_root` escapes the project (`..`) or is absolute | Loader error |
| No `main.luau` under `code_root` | “No code to run” |

Identity is always `name` from the manifest. A `love.conf` function is only used by the built-in no-game screen, not by games.

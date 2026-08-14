# Version fields

Set on the `love` table when the module loads.

## loveu

| Field | Type | Current |
|-------|------|---------|
| `love._loveu_version` | string | `"0.1.1"` |
| `love._loveu_version_major` | number | `0` |
| `love._loveu_version_minor` | number | `1` |
| `love._loveu_version_revision` | number | `0` |

Pin `loveu.toml` `engine_version` to `_loveu_version`.

## Upstream LÖVE

Same as [love.getVersion](https://love2d.org/wiki/love.getVersion) / [Config Files](https://love2d.org/wiki/Config_Files) `t.version` in stock LÖVE:

| Field | Type | Current |
|-------|------|---------|
| `love._version` | string | `"12.0"` |
| `love._version_major` | number | `12` |
| `love._version_minor` | number | `0` |
| `love._version_revision` | number | `0` |
| `love._version_codename` | string | `"Bestest Friend"` |
| `love._version_compat` | table | compatible LÖVE version strings |

```luau
local major, minor, revision, codename = love.getVersion()
```

`love --version` prints both: `loveu 0.1.1 (LÖVE 12.0 "Bestest Friend")`.

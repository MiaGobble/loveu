# love.project

Read-only table set at boot from `loveu.toml`. Missing unless a game actually loaded (the no-game screen has no project).

```luau
print(love.project.name)
print(love.project.version)
print(love.project.engine_version)
print(love.project.code_root)
```

| Field | Type | Source |
|-------|------|--------|
| `name` | string | required `name` |
| `version` | string | required `version` (the **game**, not the engine) |
| `engine_version` | string | required `engine_version` (must equal `love._loveu_version`) |
| `code_root` | string | required `code_root`, normalized (`""` → `"."`) |
| `config` | table or nil | optional window/modules/runtime keys; `nil` if none were present |

`love.project.config` is the parsed optional block (`console`, `title`, `window`, `modules`, `graphics`, `audio`, …), **not** the fully merged engine defaults. Omitted keys are not listed.

Reload via `love.filesystem.loadProjectManifest()` (same parser; mainly for tests). Boot already called it once.

Related LÖVE APIs: [love.filesystem](https://love2d.org/wiki/love.filesystem) (`getIdentity`, `getSource`, `setRequirePath`).

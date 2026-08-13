# Requires

`package.path` / the filesystem require path is **`?.luau;?/init.luau`** (prefixed with `code_root/` when that is not `"."`).

`.lua` modules are not searched.

## Bare names (root-relative)

```luau
require("lib.util")      -- code_root/lib/util.luau
require("scenes.level")  -- code_root/scenes/level.luau
                         -- or code_root/scenes/level/init.luau
```

These stay relative to **`code_root`**, even when called from a nested file.

## Relative paths

Paths starting with `./` or `../` resolve against the **calling file**, then map to a dotted module name for `package.loaded`.

```
src/
  util.luau
  scenes/
    level/
      main.luau
```

From `src/scenes/level/main.luau`:

```luau
require("../../util")   -- src/util.luau
require("./enemy")      -- src/scenes/level/enemy.luau
```

Escaping the game root (`require("../../../../outside")`) errors.

Relative and bare names that resolve to the same file share `package.loaded`.

## `code_root`

If `code_root = "src"`, requires do **not** use a doubled `src/src/…` prefix. Module names are as if `src` were the script root.

See also [love.filesystem](https://love2d.org/wiki/love.filesystem) on the wiki (`getRequirePath`, `load`, `read`).

# Getting started

loveu is a fork of [LÖVE](https://love2d.org/). This site documents **what loveu changes**. For drawing, audio, physics, input, and the rest of the engine API, use the [LÖVE wiki](https://love2d.org/wiki/Main_Page).

Current versions:

| Identifier | Meaning | Value |
|------------|---------|-------|
| `love._loveu_version` | loveu fork | `0.1.0` |
| `love._version` | Upstream LÖVE API | `12.0` |

## What you need

The **`loveu` CLI** is the usual front door (see [CLI](/guide/cli)). Install a release binary, then:

```
loveu init mygame
cd mygame
loveu run
```

`loveu run` downloads the engine pinned by `engine_version`. You can still call `love` / `lovec` directly if you have a runtime:

```
love --version
```

prints both versions, for example:

```
loveu 0.1.0 (LÖVE 12.0 "Bestest Friend")
```

Building the engine from source is covered in [Building from source](/guide/building).
## Minimal game

A game directory must contain `loveu.toml` at the **project root** and `main.luau` (at the root, or under `code_root`).

`loveu.toml`:

```toml
name = "hello"
version = "0.0.1"
engine_version = "0.1.0"
code_root = "."

[window]
width = 800
height = 600
```

`main.luau`:

```luau
function love.draw()
	love.graphics.print("Hello from Luau!", 100, 100)
end
```

Run it:

```
loveu run
```

or with a runtime directly:

```
love path/to/gamedir
```

On Windows, `lovec.exe` is the console build (stdout/stderr). `love.exe` is the GUI build.

## Layout with a `src` folder

```
mygame/
  loveu.toml
  src/
    main.luau
    util.luau
```

```toml
name = "mygame"
version = "0.1.0"
engine_version = "0.1.0"
code_root = "src"
```

`loveu.toml` stays at the project root. Scripts live under `src/`. See [Project format](/guide/project-format).

## Next

- [CLI](/guide/cli)
- [Divergences from LÖVE](/guide/divergences) — checklist of what is different
- [Building from source](/guide/building)
- [Luau scripting](/guide/luau)
- [LÖVE wiki: love](https://love2d.org/wiki/love)

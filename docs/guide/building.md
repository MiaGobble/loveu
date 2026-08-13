# Building from source

loveu does not ship prebuilt binaries yet. Build it the same way as [LÖVE 12](https://love2d.org/wiki/Building_LÖVE), using **this repository** as the engine source.

Windows typically uses [megasource](https://github.com/love2d/megasource): point its `libs/love` at this checkout, then build the `love` / `lovec` targets.

```
love --version
# loveu 0.1.0 (LÖVE 12.0 "Bestest Friend")
```

`lovec` is the console binary (stdout/stderr). Use that for headless games (`[modules] window = false`).

CI builds are in [`.github/workflows/main.yml`](https://github.com/MiaGobble/loveu/actions). CLI helpers to fetch or switch versions are not implemented yet.

# Versioning

loveu keeps **two** version strings. Do not pin `loveu.toml` to the LÖVE 12 number.

| Identifier | Meaning | Current |
|------------|---------|---------|
| `love._loveu_version` | Fork semver. **`engine_version` in `loveu.toml` must match this exactly.** | `0.1.1` |
| `love._version` | Upstream LÖVE API this build is based on | `12.0` |

```
love --version
# loveu 0.1.1 (LÖVE 12.0 "Bestest Friend")
```

## When to bump

Defined in `src/common/version.h`:

- **`LOVEU_VERSION_*`** — loveu behavior games depend on (Luau, manifest, requires, later roadmap items).
- **`LOVE_VERSION_*`** — rebasing onto a new upstream LÖVE.

`love.getVersion()` still returns the LÖVE API triple and codename. See [love.getVersion](https://love2d.org/wiki/love.getVersion).

## Compatibility check

Boot does **not** use LÖVE’s loose `love.isVersionCompatible` for the manifest. `engine_version` is an exact string match against `love._loveu_version`.

Mismatches fail immediately:

```
loveu.toml engine_version '0.0.0' does not match running loveu '0.1.1'
```

# loveu
[![CI](https://github.com/MiaGobble/loveu/actions/workflows/main.yml/badge.svg)](https://github.com/MiaGobble/loveu/actions/workflows/main.yml)

loveu is a fork of LÖVE to support different things that I want out of the engine.

## Docs
Published at [miagobble.github.io/loveu](https://miagobble.github.io/loveu/).

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
|CLI tooling|Helpers to build to any platform, run, initialize, switch versions, and more|Implemented|
|`love` type definitions|Fully featured typechecking for `love`|Not implemented|
|Documentation|VitePress docs for fork divergences, linking back to the LÖVE wiki|Implemented|

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

## CLI

Install a release installer (Windows setup.exe, macOS `.pkg`, or Linux `.deb`) from [GitHub Releases](https://github.com/MiaGobble/loveu/releases), or build from source:

```
cd cli
go build -o loveu ./cmd/loveu
```

```
loveu init mygame
cd mygame
loveu run
loveu build windows
loveu version
```

`engine_version` in `loveu.toml` selects the engine. `loveu run` / `loveu build` download matching packs from GitHub Releases into a local cache (`LOVEU_HOME`). Build targets: `love`, `windows`, `macos`, `linux`, `android`, `ios` (iOS requires macOS).

## AI Notice
This fork, at least so far, is completely done with AI. This means you can expect things to not be exceptionally implemented, or alternatively, not working the way it should.

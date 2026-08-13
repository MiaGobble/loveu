---
layout: home

hero:
  name: loveu
  text: A LÖVE fork with Luau
  tagline: Same 2D engine, different language and project format. Use the LÖVE wiki for graphics, audio, and the rest of the API.
  actions:
    - theme: brand
      text: Get started
      link: /guide/getting-started
    - theme: alt
      text: Divergences
      link: /guide/divergences
    - theme: alt
      text: LÖVE wiki
      link: https://love2d.org/wiki/Main_Page

features:
  - title: Luau instead of Lua
    details: Games are written in .luau. Entry point is main.luau. No LuaJIT FFI or jit.*.
  - title: loveu.toml
    details: Required project manifest for identity, engine pin, code root, window, and modules. There is no conf.luau.
  - title: Relative requires
    details: require("./x") and require("../x") resolve from the calling file. Bare names stay root-relative inside code_root.
  - title: Upstream LÖVE 12 API
    details: love.graphics, love.audio, callbacks, and modules still follow LÖVE 12. Read the wiki for those.
---

#!/usr/bin/env python3
"""Patch loveu Xcode projects to use Luau instead of Lua.framework."""
from __future__ import annotations

import pathlib
import re

ROOT = pathlib.Path(__file__).resolve().parents[1]
LIBLOVE = ROOT / "liblove.xcodeproj" / "project.pbxproj"
LOVE = ROOT / "love.xcodeproj" / "project.pbxproj"

# Stable-ish IDs (24 hex) for new entries.
IDS = {
    "compat_file": "B0A11C01C01C01C01C01C001",
    "compat_build_macos": "B0A11C01C01C01C01C01C002",
    "compat_build_ios": "B0A11C01C01C01C01C01C003",
    "compat_hdr_file": "B0A11C01C01C01C01C01C004",
    "compat_hdr_build": "B0A11C01C01C01C01C01C005",
    "lauxlib_file": "B0A11C01C01C01C01C01C006",
    "lauxlib_build": "B0A11C01C01C01C01C01C007",
    "run_script_macos": "B0A11C01C01C01C01C01C010",
    "run_script_ios": "B0A11C01C01C01C01C01C011",
}

LUAU_HEADERS_LINES = [
    '"$(PROJECT_DIR)/../../src/libraries/luau/VM/include",',
    '"$(PROJECT_DIR)/../../src/libraries/luau/Compiler/include",',
    '"$(PROJECT_DIR)/../../src/common",',
    '"$(PROJECT_DIR)/build/luau-$(PLATFORM_NAME)",',
]

LUAU_HEADERS = "\n".join("\t\t\t\t\t" + line for line in LUAU_HEADERS_LINES)

LUAU_LDFLAGS_LINES = [
    '"-L$(PROJECT_DIR)/build/luau-$(PLATFORM_NAME)",',
    '"-lLuau.Compiler",',
    '"-lLuau.VM",',
    '"-lLuau.Ast",',
    '"-lLuau.Bytecode",',
    '"-lLuau.Common",',
    '"-lc++",',
]

LUAU_LDFLAGS = "\n".join("\t\t\t\t\t" + line for line in LUAU_LDFLAGS_LINES)

LUAU_LIB_SEARCH = (
    '"$(inherited)",\n'
    '\t\t\t\t\t"$(PROJECT_DIR)/build/luau-$(PLATFORM_NAME)",'
)

IOS_HEADER_BLOCK = (
    "(\n"
    '\t\t\t\t\t"$(inherited)",\n'
    + "\n".join("\t\t\t\t\t" + line for line in LUAU_HEADERS_LINES)
    + "\n\t\t\t\t)"
)

IOS_LDFLAGS_BLOCK = (
    "(\n"
    '\t\t\t\t\t"-ObjC",\n'
    + "\n".join("\t\t\t\t\t" + line for line in LUAU_LDFLAGS_LINES)
    + "\n\t\t\t\t)"
)

RUN_SCRIPT = r'''
		{rid} /* Build Luau */ = {{
			isa = PBXShellScriptBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			inputPaths = (
			);
			name = "Build Luau";
			outputPaths = (
			);
			runOnlyForDeploymentPostprocessing = 0;
			shellPath = /bin/bash;
			shellScript = "export PLATFORM_NAME=\"${{PLATFORM_NAME}}\"\nexport CONFIGURATION=\"${{CONFIGURATION}}\"\nexport SDKROOT=\"${{SDKROOT}}\"\nexport ARCHS=\"${{ARCHS}}\"\nexport CURRENT_ARCH=\"${{CURRENT_ARCH}}\"\nexport MACOSX_DEPLOYMENT_TARGET=\"${{MACOSX_DEPLOYMENT_TARGET}}\"\nexport IPHONEOS_DEPLOYMENT_TARGET=\"${{IPHONEOS_DEPLOYMENT_TARGET}}\"\nbash \"${{PROJECT_DIR}}/scripts/build-luau.sh\"\n";
		}};
'''


def replace_lua_headers(text: str) -> str:
    text = text.replace(
        '"$(PROJECT_DIR)/macosx/Frameworks/Lua.framework/Headers",',
        LUAU_HEADERS,
    )
    return text


def inject_ldflags(text: str) -> str:
    """Ensure Luau link flags exist in OTHER_LDFLAGS list blocks."""

    def add_flags(match: re.Match[str]) -> str:
        block = match.group(0)
        if "-lLuau.VM" in block:
            return block
        return block.replace(
            "OTHER_LDFLAGS = (",
            "OTHER_LDFLAGS = (\n" + LUAU_LDFLAGS,
            1,
        )

    return re.sub(
        r"OTHER_LDFLAGS = \((?:[^\)]|\n)*?\);",
        add_flags,
        text,
        flags=re.M,
    )


def inject_ios_liblove_settings(text: str) -> str:
    """liblove-ios uses string HEADER/LIBRARY/LDFLAGS; expand them for Luau."""

    # Only rewrite iOS target configs that still use the inherited string form.
    text = re.sub(
        r'(SDKROOT = iphoneos;\n(?:[^\n]+\n)*?\t\t\t\t)HEADER_SEARCH_PATHS = "\$\(inherited\)";',
        lambda m: m.group(1) + "HEADER_SEARCH_PATHS = " + IOS_HEADER_BLOCK + ";",
        text,
    )
    # Order varies — also match HEADER before SDKROOT.
    text = text.replace(
        'HEADER_SEARCH_PATHS = "$(inherited)";\n\t\t\t\tLIBRARY_SEARCH_PATHS = "$(inherited)";\n\t\t\t\tMTL_ENABLE_DEBUG_INFO',
        "HEADER_SEARCH_PATHS = "
        + IOS_HEADER_BLOCK
        + ";\n\t\t\t\tLIBRARY_SEARCH_PATHS = (\n\t\t\t\t\t"
        + LUAU_LIB_SEARCH
        + "\n\t\t\t\t);\n\t\t\t\tMTL_ENABLE_DEBUG_INFO",
    )
    text = text.replace(
        'OTHER_LDFLAGS = "-ObjC";\n\t\t\t\tPRODUCT_NAME = love;\n\t\t\t\tSDKROOT = iphoneos;',
        "OTHER_LDFLAGS = "
        + IOS_LDFLAGS_BLOCK
        + ";\n\t\t\t\tPRODUCT_NAME = love;\n\t\t\t\tSDKROOT = iphoneos;",
    )
    return text


def inject_ios_app_ldflags(text: str) -> str:
    """love-ios is the final link of a static liblove.a; it must pull in Luau."""
    if "-lLuau.VM" in text and 'SDKROOT = iphoneos' in text:
        # May already be partially patched; still ensure love-ios configs have flags.
        pass

    def patch_ios_config(match: re.Match[str]) -> str:
        block = match.group(0)
        if "-lLuau.VM" in block:
            return block
        if "OTHER_LDFLAGS" in block:
            return block
        # Insert before SDKROOT = iphoneos
        return block.replace(
            "SDKROOT = iphoneos;",
            "LIBRARY_SEARCH_PATHS = (\n\t\t\t\t\t"
            + LUAU_LIB_SEARCH
            + "\n\t\t\t\t);\n\t\t\t\tOTHER_LDFLAGS = "
            + IOS_LDFLAGS_BLOCK
            + ";\n\t\t\t\tSDKROOT = iphoneos;",
            1,
        )

    # love-ios target build configurations contain PRODUCT_BUNDLE_IDENTIFIER = org.love2d.love
    return re.sub(
        r"PRODUCT_BUNDLE_IDENTIFIER = org\.love2d\.love;\n(?:[^\n]+\n)*?\t\t\t\tSDKROOT = iphoneos;",
        patch_ios_config,
        text,
        flags=re.M,
    )


def add_luau_compat_sources(text: str) -> str:
    if "luau_compat.cpp" in text:
        return text

    file_ref = (
        f'\t\t{IDS["compat_file"]} /* luau_compat.cpp */ = '
        '{isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.cpp.cpp; '
        'path = luau_compat.cpp; sourceTree = "<group>"; };\n'
        f'\t\t{IDS["compat_hdr_file"]} /* luau_compat.h */ = '
        '{isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.c.h; '
        'path = luau_compat.h; sourceTree = "<group>"; };\n'
        f'\t\t{IDS["lauxlib_file"]} /* lauxlib.h */ = '
        '{isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.c.h; '
        'path = lauxlib.h; sourceTree = "<group>"; };\n'
    )
    build_files = (
        f'\t\t{IDS["compat_build_macos"]} /* luau_compat.cpp in Sources */ = '
        f'{{isa = PBXBuildFile; fileRef = {IDS["compat_file"]} /* luau_compat.cpp */; }};\n'
        f'\t\t{IDS["compat_build_ios"]} /* luau_compat.cpp in Sources */ = '
        f'{{isa = PBXBuildFile; fileRef = {IDS["compat_file"]} /* luau_compat.cpp */; }};\n'
        f'\t\t{IDS["compat_hdr_build"]} /* luau_compat.h in Headers */ = '
        f'{{isa = PBXBuildFile; fileRef = {IDS["compat_hdr_file"]} /* luau_compat.h */; }};\n'
        f'\t\t{IDS["lauxlib_build"]} /* lauxlib.h in Headers */ = '
        f'{{isa = PBXBuildFile; fileRef = {IDS["lauxlib_file"]} /* lauxlib.h */; }};\n'
    )

    text = text.replace(
        "FA0B790E1A958E3B000E1D17 /* runtime.cpp */ = {isa = PBXFileReference;",
        file_ref + "\t\tFA0B790E1A958E3B000E1D17 /* runtime.cpp */ = {isa = PBXFileReference;",
        1,
    )
    text = text.replace(
        "FA0B793B1A958E3B000E1D17 /* runtime.cpp in Sources */ = {isa = PBXBuildFile;",
        build_files + "\t\tFA0B793B1A958E3B000E1D17 /* runtime.cpp in Sources */ = {isa = PBXBuildFile;",
        1,
    )

    text = text.replace(
        "\t\t\t\tFA0B790E1A958E3B000E1D17 /* runtime.cpp */,\n",
        "\t\t\t\tFA0B790E1A958E3B000E1D17 /* runtime.cpp */,\n"
        f'\t\t\t\t{IDS["compat_file"]} /* luau_compat.cpp */,\n'
        f'\t\t\t\t{IDS["compat_hdr_file"]} /* luau_compat.h */,\n'
        f'\t\t\t\t{IDS["lauxlib_file"]} /* lauxlib.h */,\n',
        1,
    )

    text = text.replace(
        "\t\t\t\tFA0B793C1A958E3B000E1D17 /* runtime.cpp in Sources */,\n",
        "\t\t\t\tFA0B793C1A958E3B000E1D17 /* runtime.cpp in Sources */,\n"
        f'\t\t\t\t{IDS["compat_build_ios"]} /* luau_compat.cpp in Sources */,\n',
    )
    text = text.replace(
        "\t\t\t\tFA0B793B1A958E3B000E1D17 /* runtime.cpp in Sources */,\n",
        "\t\t\t\tFA0B793B1A958E3B000E1D17 /* runtime.cpp in Sources */,\n"
        f'\t\t\t\t{IDS["compat_build_macos"]} /* luau_compat.cpp in Sources */,\n',
    )
    return text


def add_run_script_phases(text: str) -> str:
    if "Build Luau" in text:
        return text

    scripts = RUN_SCRIPT.format(rid=IDS["run_script_macos"]) + RUN_SCRIPT.format(
        rid=IDS["run_script_ios"]
    )
    text = text.replace(
        "/* Begin PBXSourcesBuildPhase section */",
        "/* Begin PBXShellScriptBuildPhase section */\n"
        + scripts
        + "/* End PBXShellScriptBuildPhase section */\n\n"
        + "/* Begin PBXSourcesBuildPhase section */",
        1,
    )

    def prepend_phase(match: re.Match[str], rid: str) -> str:
        block = match.group(0)
        if rid in block:
            return block
        return block.replace(
            "buildPhases = (\n",
            f"buildPhases = (\n\t\t\t\t{rid} /* Build Luau */,\n",
            1,
        )

    text = re.sub(
        r"/\* liblove-ios \*/ = \{.*?buildPhases = \(.*?\);",
        lambda m: prepend_phase(m, IDS["run_script_ios"]),
        text,
        count=1,
        flags=re.S,
    )
    text = re.sub(
        r"/\* liblove-macosx \*/ = \{.*?buildPhases = \(.*?\);",
        lambda m: prepend_phase(m, IDS["run_script_macos"]),
        text,
        count=1,
        flags=re.S,
    )
    return text


def remove_lua_xcframework_link(text: str) -> str:
    return text.replace(
        "\t\t\t\tFACFB753276D7F860089F78D /* Lua.xcframework in Frameworks */,\n",
        "",
    )


def patch_liblove() -> None:
    text = LIBLOVE.read_text(encoding="utf-8")
    text = replace_lua_headers(text)
    text = inject_ldflags(text)
    text = inject_ios_liblove_settings(text)
    text = add_luau_compat_sources(text)
    text = add_run_script_phases(text)
    text = remove_lua_xcframework_link(text)
    LIBLOVE.write_text(text, encoding="utf-8", newline="\n")
    print(f"patched {LIBLOVE}")


def patch_love_app() -> None:
    text = LOVE.read_text(encoding="utf-8")
    text = replace_lua_headers(text)
    text = text.replace(
        "\t\t\t\tA93E6E5510420B57007D418B /* Lua.framework in Frameworks */,\n",
        "",
    )
    text = text.replace(
        "\t\t\t\tA9255DD31043183600BA1496 /* Lua.framework in Copy Frameworks */,\n",
        "",
    )
    text = inject_ios_app_ldflags(text)
    LOVE.write_text(text, encoding="utf-8", newline="\n")
    print(f"patched {LOVE}")


if __name__ == "__main__":
    patch_liblove()
    patch_love_app()

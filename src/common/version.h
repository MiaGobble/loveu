/**
 * Copyright (c) 2006-2026 LOVE Development Team
 *
 * This software is provided 'as-is', without any express or implied
 * warranty.  In no event will the authors be held liable for any damages
 * arising from the use of this software.
 *
 * Permission is granted to anyone to use this software for any purpose,
 * including commercial applications, and to alter it and redistribute it
 * freely, subject to the following restrictions:
 *
 * 1. The origin of this software must not be misrepresented; you must not
 *    claim that you wrote the original software. If you use this software
 *    in a product, an acknowledgment in the product documentation would be
 *    appreciated but is not required.
 * 2. Altered source versions must be plainly marked as such, and must not be
 *    misrepresented as being the original software.
 * 3. This notice may not be removed or altered from any source distribution.
 **/

#ifndef LOVE_VERSION_H
#define LOVE_VERSION_H

namespace love
{

// Upstream LÖVE API version this fork is based on (love.getVersion / conf t.version).
#define LOVE_VERSION_STRING "12.0"
static const int VERSION_MAJOR = 12;
static const int VERSION_MINOR = 0;
static const int VERSION_REV = 0;
static const char *VERSION = LOVE_VERSION_STRING;
static const char *VERSION_COMPATIBILITY[] =  { VERSION, 0 };
static const char *VERSION_CODENAME = "Bestest Friend";

// loveu fork semver — loveu.toml engine_version must match this exactly.
// Bump when loveu ships breaking/behavioral changes; bump LOVE_* when rebasing upstream.
#define LOVEU_VERSION_STRING "0.1.1"
static const int LOVEU_VERSION_MAJOR = 0;
static const int LOVEU_VERSION_MINOR = 1;
static const int LOVEU_VERSION_REV = 1;
static const char *LOVEU_VERSION = LOVEU_VERSION_STRING;

} // love

#endif // LOVE_VERSION_H

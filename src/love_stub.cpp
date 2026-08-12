/**
 * Thin executable entry: version check + love_main() from liblove.
 * Keeps Luau entirely inside liblove so the love binary does not need lua_*.
 **/

#include "common/config.h"
#include "common/version.h"
#include "modules/love/love.h"

#include <stdio.h>
#include <string.h>

#include <SDL3/SDL_main.h>

#ifdef LOVE_WINDOWS
#define WIN32_LEAN_AND_MEAN
#include <windows.h>

extern "C"
{
// Prefer the higher performance GPU on Windows systems that use nvidia Optimus.
LOVE_EXPORT DWORD NvOptimusEnablement = 1;
// Same with AMD GPUs.
LOVE_EXPORT DWORD AmdPowerXpressRequestHighPerformance = 1;
}

__declspec(dllimport) int love_main(int argc, char **argv);
#else
int love_main(int argc, char **argv);
#endif

int main(int argc, char **argv)
{
	if (strcmp(LOVE_VERSION_STRING, love_version()) != 0)
	{
		printf("Version mismatch detected!\nLOVE binary is version %s\n"
			   "LOVE library is version %s\n", LOVE_VERSION_STRING, love_version());
		return 1;
	}

	return love_main(argc, argv);
}

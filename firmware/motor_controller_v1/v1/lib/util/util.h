#ifndef UTIL_H_
#include <string>

// Major version definition
#ifndef VERSION_MAJOR
#define VERSION_MAJOR 0
#endif

// Minor version definition
#ifndef VERSION_MINOR
#define VERSION_MINOR 0
#endif

// Path version definition
#ifndef VERSION_PATCH
#define VERSION_PATCH 0
#endif

typedef struct
{
    uint8_t major = VERSION_MAJOR;
    uint8_t minor = VERSION_MINOR;
    uint8_t patch = VERSION_PATCH;
} Version;

inline std::string GetVersionString()
{
    char version[12];
    sprintf(version, "%s.%s.%s", VERSION_MAJOR, VERSION_MINOR, VERSION_PATCH);
    return std::string(version);
}

inline Version GetVersion()
{
    return Version{};
}

#endif
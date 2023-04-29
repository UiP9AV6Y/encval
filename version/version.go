package version

// this is a dummy package intended with no actual functionality.
//
// runtime/debug is unable to provide module information about the
// module calling it. this module lives inside the main
// application/library and therefor has the same version information,
// which can then be extracted via its debug.BuildInfo

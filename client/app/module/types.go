package module

type VersionedModule struct {
	AppModuleBasic

	Module AppModule
	// FromVersion and ToVersion indicate the continuous range of app versions
	// that this particular module is part of. The range is inclusive.
	// FromVersion should not be smaller than ToVersion. 0 is not a valid app
	// version.
	FromVersion, ToVersion uint64
}

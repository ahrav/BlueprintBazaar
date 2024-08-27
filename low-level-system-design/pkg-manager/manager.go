package pkgmanager

import "errors"

type Package struct {
	Name         string
	Version      string
	Dependencies []*Package
	metadata     map[string]string
}

type Installer interface {
	Install(pkg *Package) error
}

type PackageManager struct {
	version    string
	packages   map[string]*Package
	installing map[string]bool
	installed  map[string]bool
	installer  Installer
}

func NewPackageManager(installer Installer) *PackageManager {
	return &PackageManager{
		installing: make(map[string]bool),
		installed:  make(map[string]bool),
		installer:  installer,
	}
}

var ErrDuplicatePackage = errors.New("duplicate package")

func (pm *PackageManager) AddPackage(name string, dependencies ...*Package) (*Package, error) {
	if _, exists := pm.packages[name]; exists {
		return nil, ErrDuplicatePackage
	}

	pkg := &Package{Name: name, Dependencies: dependencies}
	pm.packages[name] = pkg

	return pkg, nil
}

var (
	ErrCyclicDependency = errors.New("cyclic dependency detected")
	ErrPkgMissing       = errors.New("package not present in package manager")
)

func (pm *PackageManager) Install(name string) error {
	if pm.installed[name] {
		return nil
	}

	if pm.installing[name] {
		return ErrCyclicDependency
	}

	pm.installing[name] = true
	defer delete(pm.installing, name)

	pkg, exists := pm.packages[name]
	if !exists {
		return ErrPkgMissing
	}

	for _, dep := range pm.packages[name].Dependencies {
		if err := pm.Install(dep.Name); err != nil {
			return err
		}
	}

	if err := pm.installer.Install(pkg); err != nil {
		return err
	}
	pm.installed[name] = true

	return nil
}

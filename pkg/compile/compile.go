package compile

import (
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/registry"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphemap-ts/pkg/mapdef"
)

// TsConverterConfig holds the configuration for TypeScript converter generation.
type TsConverterConfig struct {
	MapsPath              string
	RegistryConfig        rcfg.MorpheLoadRegistryConfig
	ExternalConfig        *rcfg.MorpheLoadRegistryConfig
	OutputPath            string
	SourceTypesImportPath string
	TargetTypesImportPath string
}

// MorpheMapToTypeScript generates TypeScript converter functions from MorpheMap definitions.
func MorpheMapToTypeScript(config TsConverterConfig) error {
	// Load maps
	maps, err := mapdef.LoadMapsFromDirectory(config.MapsPath)
	if err != nil {
		return fmt.Errorf("failed to load maps: %w", err)
	}

	if len(maps) == 0 {
		return fmt.Errorf("no .map files found in %s", config.MapsPath)
	}

	// Load local Morphe registry
	localRegistry, err := registry.LoadMorpheRegistry(
		registry.LoadMorpheRegistryHooks{},
		config.RegistryConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to load local Morphe registry: %w", err)
	}

	// Load external registry if configured
	var externalRegistry *registry.Registry
	if config.ExternalConfig != nil {
		extReg, err := registry.LoadMorpheRegistry(
			registry.LoadMorpheRegistryHooks{},
			*config.ExternalConfig,
		)
		if err != nil {
			return fmt.Errorf("failed to load external Morphe registry: %w", err)
		}
		externalRegistry = extReg
	}

	// Generate a TS file for each map
	for _, m := range maps {
		mapType := m.InferMapType()

		switch mapType {
		case mapdef.MapTypeField:
			if err := compileFieldMap(&m, localRegistry, externalRegistry, config); err != nil {
				return fmt.Errorf("failed to compile field map %q: %w", m.Name, err)
			}

		case mapdef.MapTypeEnum:
			if err := compileEnumMap(&m, config); err != nil {
				return fmt.Errorf("failed to compile enum map %q: %w", m.Name, err)
			}
		}
	}

	return nil
}

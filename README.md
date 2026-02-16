# plugin-morphemap-ts

Kalo plugin that generates TypeScript converter functions from MorpheMap (`.map`) definitions.

## Overview

Given MorpheMap files and Morphe type registries, this plugin generates TypeScript code that converts data between structurally different source and target types. Each `.map` file produces a corresponding `.ts` file with:

- **Field maps**: Converter functions that take the source type and return the target type, handling direct renames, path traversal, type coercion, enum translation, conditionals, and constants
- **Enum maps**: Enum translation functions with lookup maps and error handling

This plugin is designed for **cross-domain structural mapping** -- converting between external API types and local domain models where field-level mapping decisions are captured in `.map` files.

> **Note:** For internal serialization casing bridges (e.g., `snake_case` wire format to `camelCase` application types), use the dedicated `plugin-morphe-ts-casing-bridge` instead. See [ADR-001](https://github.com/kalo-build/kalo-plugin-registry/blob/main/docs/decisions/001-casing-bridge-architecture.md) for rationale.

## Input

- **MorpheMap files** (`KA:MM1:YAML1`): Transformation definitions
- **Local Morphe Registry** (`KA:MO1:YAML1`): Local project Morphe schema files
- **External Morphe Registry** (`KA:MO1:YAML1`, optional): Third-party API type definitions

## Output

- **TypeScript files** (`KA:MM1:TS1`): Converter functions

## Configuration

```yaml
config:
  "@kalo-build/plugin-morphemap-ts":
    sourceTypesImportPath: "@/types/external"
    targetTypesImportPath: "@/types/models"
```

### Config Options

| Option | Type | Required | Default | Description |
|--------|------|----------|---------|-------------|
| `sourceTypesImportPath` | string | Yes | - | Import path for source TS types |
| `targetTypesImportPath` | string | Yes | - | Import path for target TS types |

## Generated Output Example

### Field Map Converter

```typescript
import type { RealEstateHouseBuy } from '@/types/external';
import type { RealEstateListing } from '@/types/models';

export function is24HouseBuyToRealEstateListing(
  source: RealEstateHouseBuy
): RealEstateListing {
  const target: Partial<RealEstateListing> = {};

  target.title = source.title;
  target.streetName = source.address.street;
  target.type = 'Haus';
  target.numberOfRooms = Number(source.numRooms);

  return target as RealEstateListing;
}
```

## Build

```bash
cd scripts && bash build.sh
```

## License

MIT

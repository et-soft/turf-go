# Transformation Package Migration

## Overview
The `intersect` method has been successfully moved from the root package to a dedicated `transformation` package, following the structure of the original Turf.js library.

## Changes Made

### 1. Directory Structure
- ✅ Created `transformation/` directory
- ✅ Moved `intersect.go` to `transformation/intersect.go`
- ✅ Moved `intersect_test.go` to `transformation/intersect_test.go`
- ✅ Created `transformation/README.md` with package documentation

### 2. Code Updates
- ✅ Updated package declaration to `package transformation`
- ✅ Updated imports in test files
- ✅ Updated example usage to use `transformation.Intersect()`
- ✅ Removed old files from root directory

### 3. Documentation
- ✅ Updated main README.md to indicate intersect is in `transformation` package
- ✅ Created transformation package documentation
- ✅ Updated example imports

## New Usage

### Before (Root Package)
```go
import "github.com/tomchavakis/turf-go"

result, err := turf.Intersect(poly1, poly2)
```

### After (Transformation Package)
```go
import "github.com/tomchavakis/turf-go/transformation"

result, err := transformation.Intersect(poly1, poly2)
```

## File Structure
```
turf-go/
├── transformation/
│   ├── intersect.go
│   ├── intersect_test.go
│   └── README.md
├── examples/
│   └── example_intersect.go
└── README.md (updated)
```

## Testing
- ✅ All tests pass in transformation package
- ✅ Example works correctly with new import path
- ✅ No breaking changes to existing functionality

## Status
✅ **COMPLETED** - The intersect method has been successfully migrated to the transformation package, maintaining full functionality while following the Turf.js project structure.

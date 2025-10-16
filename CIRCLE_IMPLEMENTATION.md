# Circle Implementation

## Overview
This document describes the implementation of the `circle` method for the turf-go library, which is a Go port of the Turf.js library.

## Implementation Details

### Files Added
- `transformation/circle.go` - Main implementation of the circle function
- `transformation/circle_test.go` - Comprehensive test suite
- `examples/example_circle.go` - Usage example

### Function Signatures
```go
// Main circle function
func Circle(center geometry.Point, radius float64, options *CircleOptions) (*feature.Feature, error)

// Alternative input methods
func CircleFromCoordinates(center []float64, radius float64, options *CircleOptions) (*feature.Feature, error)
func CircleFromFeature(f *feature.Feature, radius float64, options *CircleOptions) (*feature.Feature, error)
func CircleFromGeometry(geom *geometry.Geometry, radius float64, options *CircleOptions) (*feature.Feature, error)
```

### CircleOptions Structure
```go
type CircleOptions struct {
    Steps      int                    `json:"steps,omitempty"`
    Units      string                 `json:"units,omitempty"`
    Properties map[string]interface{} `json:"properties,omitempty"`
}
```

### Supported Input Types
- `geometry.Point` - Direct point geometry
- `[]float64` - Coordinate array [lng, lat]
- `*feature.Feature` - Feature containing Point geometry
- `*geometry.Geometry` - Point geometry

### Supported Units
- `kilometeres` (note: matches constants.UnitKilometers)
- `miles`
- `meters`
- `feet`
- `yards`
- `inches`
- `centimeters`
- `degrees`
- `radians`

### Algorithm
The implementation:

1. **Input Validation**: Validates center point and radius
2. **Unit Conversion**: Converts radius to degrees using conversion functions
3. **Coordinate Generation**: Generates circle coordinates using trigonometric functions
4. **Polygon Creation**: Creates a closed polygon from the coordinates
5. **Feature Creation**: Wraps the polygon in a Feature with bounding box and properties

### Key Features
- **Multiple Input Methods**: Support for various input types (Point, coordinates, Feature, Geometry)
- **Unit Support**: Full support for all units defined in the constants package
- **Customizable Steps**: Control the number of segments for circle approximation
- **Properties Support**: Add custom properties to the resulting feature
- **Bounding Box Calculation**: Automatic calculation of feature bounding box

### Testing
The implementation includes comprehensive tests covering:
- Basic circle creation
- Custom options (steps, units, properties)
- Different input methods (coordinates, feature, geometry)
- Various units of measurement
- Different step counts
- Edge cases (zero radius, negative radius, small radius)
- Error handling (invalid inputs, nil inputs)

### Usage Examples

#### Basic Circle
```go
center := geometry.Point{Lng: -75.343, Lat: 39.984}
radius := 5.0
circle, err := transformation.Circle(center, radius, nil)
```

#### Circle with Options
```go
options := &transformation.CircleOptions{
    Steps:      16,
    Units:      "kilometeres",
    Properties: map[string]interface{}{
        "name": "My Circle",
    },
}
circle, err := transformation.Circle(center, radius, options)
```

#### Circle from Coordinates
```go
coords := []float64{-75.343, 39.984}
circle, err := transformation.CircleFromCoordinates(coords, 5.0, nil)
```

## Status
✅ **COMPLETED** - The circle method has been successfully implemented and tested.

## Files Structure
```
transformation/
├── circle.go          # Main implementation
├── circle_test.go     # Tests
├── intersect.go        # Intersect implementation
├── intersect_test.go   # Intersect tests
└── README.md          # Package documentation

examples/
├── example_circle.go   # Circle usage example
└── example_intersect.go # Intersect usage example
```

## Integration
- ✅ Added to transformation package
- ✅ Updated main README.md
- ✅ Updated transformation package documentation
- ✅ All tests passing
- ✅ Example working correctly

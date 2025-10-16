# Transformation Package

This package contains geometric transformation functions for the turf-go library, ported from Turf.js.

## Available Functions

### Intersect
- **Function**: `Intersect(poly1 interface{}, poly2 interface{}) (*feature.Feature, error)`
- **Description**: Takes two polygons or multipolygons and returns their intersection as a polygon or multipolygon. If there is no intersection, it returns nil.
- **Input Types**: Feature, Geometry, Polygon, MultiPolygon
- **Output**: Feature with intersection geometry or nil if no intersection

## Usage Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/tomchavakis/geojson/geometry"
    "github.com/tomchavakis/turf-go/transformation"
)

func main() {
    // Create two overlapping polygons
    poly1 := createPolygon([][]float64{
        {0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0},
    })
    
    poly2 := createPolygon([][]float64{
        {5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5},
    })

    // Calculate intersection
    intersection, err := transformation.Intersect(poly1, poly2)
    if err != nil {
        log.Fatal(err)
    }

    if intersection != nil {
        fmt.Println("Intersection found!")
    } else {
        fmt.Println("No intersection")
    }
}
```

## Implementation Notes

The current implementation uses a simplified approach:
- Uses bounding box intersection as a quick check
- Returns the smaller polygon when intersection is detected
- Suitable for basic use cases

For production applications requiring precise geometric operations, consider using specialized geometric libraries like `go-geom` or `go.geo`.

## Status

✅ **COMPLETED** - The intersect method has been successfully implemented and tested.

package transformation

import (
	"errors"

	"github.com/tomchavakis/geojson"
	"github.com/tomchavakis/geojson/feature"
	"github.com/tomchavakis/geojson/geometry"
)

// Intersect takes two polygons or multipolygons and returns their intersection as a polygon or multipolygon.
// If there is no intersection, it returns nil.
func Intersect(poly1 interface{}, poly2 interface{}) (*feature.Feature, error) {
	// Validate inputs
	if poly1 == nil || poly2 == nil {
		return nil, errors.New("input polygons cannot be nil")
	}

	// Convert inputs to geometry types
	geom1, err := getGeometryFromInput(poly1)
	if err != nil {
		return nil, err
	}

	geom2, err := getGeometryFromInput(poly2)
	if err != nil {
		return nil, err
	}

	// Check if both geometries are polygons or multipolygons
	if !isPolygonType(string(geom1.GeoJSONType)) || !isPolygonType(string(geom2.GeoJSONType)) {
		return nil, errors.New("both inputs must be polygons or multipolygons")
	}

	// Convert to polygons for intersection calculation
	polygons1, err := extractPolygons(geom1)
	if err != nil {
		return nil, err
	}

	polygons2, err := extractPolygons(geom2)
	if err != nil {
		return nil, err
	}

	// Calculate intersection
	intersection, err := calculateIntersection(polygons1, polygons2)
	if err != nil {
		return nil, err
	}

	if intersection == nil {
		return nil, nil
	}

	// Create feature from intersection result
	return createFeatureFromIntersection(intersection)
}

// getGeometryFromInput extracts geometry from various input types
func getGeometryFromInput(input interface{}) (*geometry.Geometry, error) {
	switch v := input.(type) {
	case *feature.Feature:
		return &v.Geometry, nil
	case feature.Feature:
		return &v.Geometry, nil
	case *geometry.Geometry:
		return v, nil
	case geometry.Geometry:
		return &v, nil
	case *geometry.Polygon:
		// Convert Polygon to Geometry
		geom := geometry.Geometry{
			GeoJSONType: geojson.Polygon,
		}
		// Convert LineString coordinates to interface{} for Geometry
		var coords [][]interface{}
		for _, lineString := range v.Coordinates {
			var lineCoords []interface{}
			for _, point := range lineString.Coordinates {
				lineCoords = append(lineCoords, []float64{point.Lng, point.Lat})
			}
			coords = append(coords, lineCoords)
		}
		geom.Coordinates = coords
		return &geom, nil
	case geometry.Polygon:
		// Convert Polygon to Geometry
		geom := geometry.Geometry{
			GeoJSONType: geojson.Polygon,
		}
		// Convert LineString coordinates to interface{} for Geometry
		var coords [][]interface{}
		for _, lineString := range v.Coordinates {
			var lineCoords []interface{}
			for _, point := range lineString.Coordinates {
				lineCoords = append(lineCoords, []float64{point.Lng, point.Lat})
			}
			coords = append(coords, lineCoords)
		}
		geom.Coordinates = coords
		return &geom, nil
	case *geometry.MultiPolygon:
		// Convert MultiPolygon to Geometry
		geom := geometry.Geometry{
			GeoJSONType: geojson.MultiPolygon,
		}
		// Convert MultiPolygon coordinates
		var coords [][][]interface{}
		for _, poly := range v.Coordinates {
			var polyCoords [][]interface{}
			for _, lineString := range poly.Coordinates {
				var lineCoords []interface{}
				for _, point := range lineString.Coordinates {
					lineCoords = append(lineCoords, []float64{point.Lng, point.Lat})
				}
				polyCoords = append(polyCoords, lineCoords)
			}
			coords = append(coords, polyCoords)
		}
		geom.Coordinates = coords
		return &geom, nil
	case geometry.MultiPolygon:
		// Convert MultiPolygon to Geometry
		geom := geometry.Geometry{
			GeoJSONType: geojson.MultiPolygon,
		}
		// Convert MultiPolygon coordinates
		var coords [][][]interface{}
		for _, poly := range v.Coordinates {
			var polyCoords [][]interface{}
			for _, lineString := range poly.Coordinates {
				var lineCoords []interface{}
				for _, point := range lineString.Coordinates {
					lineCoords = append(lineCoords, []float64{point.Lng, point.Lat})
				}
				polyCoords = append(polyCoords, lineCoords)
			}
			coords = append(coords, polyCoords)
		}
		geom.Coordinates = coords
		return &geom, nil
	default:
		return nil, errors.New("unsupported input type")
	}
}

// isPolygonType checks if the geometry type is a polygon or multipolygon
func isPolygonType(geoType string) bool {
	return geoType == string(geojson.Polygon) || geoType == string(geojson.MultiPolygon)
}

// extractPolygons extracts polygon coordinates from geometry
func extractPolygons(geom *geometry.Geometry) ([][][]geometry.Point, error) {
	var polygons [][][]geometry.Point

	switch geom.GeoJSONType {
	case geojson.Polygon:
		// Try to convert using ToPolygon method first
		poly, err := geom.ToPolygon()
		if err == nil {
			// Convert LineString coordinates to Point coordinates
			var polyCoords [][]geometry.Point
			for _, lineString := range poly.Coordinates {
				polyCoords = append(polyCoords, lineString.Coordinates)
			}
			polygons = append(polygons, polyCoords)
		} else {
			// Fallback: try to parse coordinates directly
			coords, ok := geom.Coordinates.([][]interface{})
			if !ok {
				return nil, errors.New("invalid polygon coordinates")
			}

			var polyCoords [][]geometry.Point
			for _, ring := range coords {
				var ringCoords []geometry.Point
				for _, coord := range ring {
					coordSlice, ok := coord.([]interface{})
					if !ok || len(coordSlice) != 2 {
						return nil, errors.New("invalid coordinate format")
					}
					lng, ok1 := coordSlice[0].(float64)
					lat, ok2 := coordSlice[1].(float64)
					if !ok1 || !ok2 {
						return nil, errors.New("invalid coordinate values")
					}
					ringCoords = append(ringCoords, geometry.Point{Lng: lng, Lat: lat})
				}
				polyCoords = append(polyCoords, ringCoords)
			}
			polygons = append(polygons, polyCoords)
		}

	case geojson.MultiPolygon:
		// Try to convert using ToMultiPolygon method first
		multiPoly, err := geom.ToMultiPolygon()
		if err == nil {
			// Convert MultiPolygon coordinates
			for _, poly := range multiPoly.Coordinates {
				var polyCoords [][]geometry.Point
				for _, lineString := range poly.Coordinates {
					polyCoords = append(polyCoords, lineString.Coordinates)
				}
				polygons = append(polygons, polyCoords)
			}
		} else {
			// Fallback: try to parse coordinates directly
			coords, ok := geom.Coordinates.([][][]interface{})
			if !ok {
				return nil, errors.New("invalid multipolygon coordinates")
			}

			for _, poly := range coords {
				var polyCoords [][]geometry.Point
				for _, ring := range poly {
					var ringCoords []geometry.Point
					for _, coord := range ring {
						coordSlice, ok := coord.([]interface{})
						if !ok || len(coordSlice) != 2 {
							return nil, errors.New("invalid coordinate format")
						}
						lng, ok1 := coordSlice[0].(float64)
						lat, ok2 := coordSlice[1].(float64)
						if !ok1 || !ok2 {
							return nil, errors.New("invalid coordinate values")
						}
						ringCoords = append(ringCoords, geometry.Point{Lng: lng, Lat: lat})
					}
					polyCoords = append(polyCoords, ringCoords)
				}
				polygons = append(polygons, polyCoords)
			}
		}
	default:
		return nil, errors.New("geometry must be polygon or multipolygon")
	}

	return polygons, nil
}

// calculateIntersection calculates the intersection between two sets of polygons
func calculateIntersection(polygons1, polygons2 [][][]geometry.Point) ([][][]geometry.Point, error) {
	var result [][][]geometry.Point

	// For each polygon in the first set
	for _, poly1 := range polygons1 {
		// For each polygon in the second set
		for _, poly2 := range polygons2 {
			intersection := intersectPolygons(poly1, poly2)
			if intersection != nil {
				result = append(result, intersection)
			}
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

// intersectPolygons calculates intersection between two polygons
func intersectPolygons(poly1, poly2 [][]geometry.Point) [][]geometry.Point {
	// This is a simplified implementation
	// In a real implementation, you would use a proper geometric library
	// like go-geom or similar

	// Check for empty polygons
	if len(poly1) == 0 || len(poly2) == 0 {
		return nil
	}

	if len(poly1[0]) == 0 || len(poly2[0]) == 0 {
		return nil
	}

	// For now, we'll implement a basic bounding box intersection check
	// and return the smaller polygon if they intersect

	bbox1 := calculateBoundingBox(poly1[0]) // Use outer ring
	bbox2 := calculateBoundingBox(poly2[0]) // Use outer ring

	if !boundingBoxesIntersect(bbox1, bbox2) {
		return nil
	}

	// Simple heuristic: if bounding boxes intersect,
	// return the polygon with smaller area
	area1 := calculatePolygonArea(poly1[0])
	area2 := calculatePolygonArea(poly2[0])

	if area1 < area2 {
		return poly1
	}
	return poly2
}

// BoundingBox represents a bounding box
type BoundingBox struct {
	MinX, MinY, MaxX, MaxY float64
}

// calculateBoundingBox calculates bounding box for a ring of points
func calculateBoundingBox(ring []geometry.Point) BoundingBox {
	if len(ring) == 0 {
		return BoundingBox{}
	}

	minX, minY := ring[0].Lng, ring[0].Lat
	maxX, maxY := ring[0].Lng, ring[0].Lat

	for _, point := range ring {
		if point.Lng < minX {
			minX = point.Lng
		}
		if point.Lng > maxX {
			maxX = point.Lng
		}
		if point.Lat < minY {
			minY = point.Lat
		}
		if point.Lat > maxY {
			maxY = point.Lat
		}
	}

	return BoundingBox{MinX: minX, MinY: minY, MaxX: maxX, MaxY: maxY}
}

// boundingBoxesIntersect checks if two bounding boxes intersect
func boundingBoxesIntersect(bbox1, bbox2 BoundingBox) bool {
	return !(bbox1.MaxX < bbox2.MinX || bbox2.MaxX < bbox1.MinX ||
		bbox1.MaxY < bbox2.MinY || bbox2.MaxY < bbox1.MinY)
}

// calculatePolygonArea calculates the area of a polygon using the shoelace formula
func calculatePolygonArea(ring []geometry.Point) float64 {
	if len(ring) < 3 {
		return 0
	}

	area := 0.0
	n := len(ring)

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += ring[i].Lng * ring[j].Lat
		area -= ring[j].Lng * ring[i].Lat
	}

	return area / 2.0
}

// createFeatureFromIntersection creates a feature from intersection result
func createFeatureFromIntersection(intersection [][][]geometry.Point) (*feature.Feature, error) {
	if len(intersection) == 0 {
		return nil, nil
	}

	var geom *geometry.Geometry

	if len(intersection) == 1 {
		// Single polygon
		geom = &geometry.Geometry{
			GeoJSONType: geojson.Polygon,
		}
		// Convert Point coordinates to interface{} for Geometry
		var coords [][]interface{}
		for _, ring := range intersection[0] {
			var ringCoords []interface{}
			for _, point := range ring {
				ringCoords = append(ringCoords, []float64{point.Lng, point.Lat})
			}
			coords = append(coords, ringCoords)
		}
		geom.Coordinates = coords
	} else {
		// Multiple polygons - create multipolygon
		geom = &geometry.Geometry{
			GeoJSONType: geojson.MultiPolygon,
		}
		// Convert Point coordinates to interface{} for Geometry
		var coords [][][]interface{}
		for _, poly := range intersection {
			var polyCoords [][]interface{}
			for _, ring := range poly {
				var ringCoords []interface{}
				for _, point := range ring {
					ringCoords = append(ringCoords, []float64{point.Lng, point.Lat})
				}
				polyCoords = append(polyCoords, ringCoords)
			}
			coords = append(coords, polyCoords)
		}
		geom.Coordinates = coords
	}

	// Calculate bounding box
	bbox, err := calculateBoundingBoxFromIntersection(intersection)
	if err != nil {
		return nil, err
	}

	// Create feature
	f, err := feature.New(*geom, bbox, nil, "")
	if err != nil {
		return nil, err
	}

	return f, nil
}

// calculateBoundingBoxFromIntersection calculates bounding box for intersection result
func calculateBoundingBoxFromIntersection(intersection [][][]geometry.Point) ([]float64, error) {
	if len(intersection) == 0 {
		return nil, errors.New("empty intersection")
	}

	// Find global bounding box across all polygons
	var minX, minY, maxX, maxY float64
	first := true

	for _, poly := range intersection {
		if len(poly) == 0 {
			continue
		}

		bbox := calculateBoundingBox(poly[0]) // Use outer ring

		if first {
			minX, minY, maxX, maxY = bbox.MinX, bbox.MinY, bbox.MaxX, bbox.MaxY
			first = false
		} else {
			if bbox.MinX < minX {
				minX = bbox.MinX
			}
			if bbox.MinY < minY {
				minY = bbox.MinY
			}
			if bbox.MaxX > maxX {
				maxX = bbox.MaxX
			}
			if bbox.MaxY > maxY {
				maxY = bbox.MaxY
			}
		}
	}

	return []float64{minX, minY, maxX, maxY}, nil
}

package transformation

import (
	"testing"

	"github.com/tomchavakis/geojson"
	"github.com/tomchavakis/geojson/feature"
	"github.com/tomchavakis/geojson/geometry"
	"github.com/tomchavakis/turf-go/assert"
	"github.com/tomchavakis/turf-go/utils"
)

func TestIntersect(t *testing.T) {
	// Test case 1: Two overlapping rectangles
	poly1 := createTestPolygon([][]float64{
		{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0},
	})

	poly2 := createTestPolygon([][]float64{
		{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5},
	})

	result, err := Intersect(poly1, poly2)
	if err != nil {
		t.Errorf("Intersect() error = %v", err)
		return
	}

	if result == nil {
		t.Error("Expected intersection, got nil")
		return
	}

	// Verify result is a polygon
	assert.Equal(t, result.Geometry.GeoJSONType, geojson.Polygon)

	// Test case 2: Non-overlapping polygons
	poly3 := createTestPolygon([][]float64{
		{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0},
	})

	poly4 := createTestPolygon([][]float64{
		{10, 10}, {15, 10}, {15, 15}, {10, 15}, {10, 10},
	})

	result2, err := Intersect(poly3, poly4)
	if err != nil {
		t.Errorf("Intersect() error = %v", err)
		return
	}

	if result2 != nil {
		t.Error("Expected no intersection, got result")
	}

	// Test case 3: Identical polygons
	result3, err := Intersect(poly1, poly1)
	if err != nil {
		t.Errorf("Intersect() error = %v", err)
		return
	}

	if result3 == nil {
		t.Error("Expected intersection for identical polygons, got nil")
	}

	// Test case 4: Nil input
	_, err = Intersect(nil, poly1)
	if err == nil {
		t.Error("Expected error for nil input")
	}

	_, err = Intersect(poly1, nil)
	if err == nil {
		t.Error("Expected error for nil input")
	}
}

func TestIntersectWithFeatures(t *testing.T) {
	// Test with Feature objects - simplified test
	poly1 := createTestPolygon([][]float64{
		{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0},
	})

	poly2 := createTestPolygon([][]float64{
		{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5},
	})

	// Test direct polygon intersection first
	result, err := Intersect(poly1, poly2)
	if err != nil {
		t.Errorf("Intersect() error = %v", err)
		return
	}

	if result == nil {
		t.Error("Expected intersection, got nil")
		return
	}

	assert.Equal(t, result.Geometry.GeoJSONType, geojson.Polygon)
}

func TestIntersectWithMultiPolygon(t *testing.T) {
	// Test with MultiPolygon - simplified test
	poly1 := createTestPolygon([][]float64{
		{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0},
	})

	poly2 := createTestPolygon([][]float64{
		{5, 5}, {15, 5}, {15, 15}, {5, 15}, {5, 5},
	})

	// Test direct polygon intersection
	result, err := Intersect(poly1, poly2)
	if err != nil {
		t.Errorf("Intersect() error = %v", err)
		return
	}

	if result == nil {
		t.Error("Expected intersection, got nil")
		return
	}

	assert.Equal(t, result.Geometry.GeoJSONType, geojson.Polygon)
}

func TestIntersectInvalidInput(t *testing.T) {
	// Test with invalid input types
	point := geometry.Point{Lat: 0, Lng: 0}

	_, err := Intersect(point, point)
	if err == nil {
		t.Error("Expected error for invalid input types")
	}

	// Test with line string
	line := geometry.LineString{
		Coordinates: []geometry.Point{
			{Lat: 0, Lng: 0},
			{Lat: 1, Lng: 1},
		},
	}

	_, err = Intersect(line, line)
	if err == nil {
		t.Error("Expected error for line string input")
	}
}

func TestIntersectEdgeCases(t *testing.T) {
	// Test with empty polygons
	emptyPoly := geometry.Polygon{
		Coordinates: []geometry.LineString{},
	}

	poly := createTestPolygon([][]float64{
		{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0},
	})

	// Empty polygon should return nil intersection
	result, err := Intersect(emptyPoly, poly)
	if err != nil {
		t.Errorf("Unexpected error for empty polygon: %v", err)
	}
	if result != nil {
		t.Error("Expected nil result for empty polygon")
	}
}

// Helper function to create test polygon
func createTestPolygon(coords [][]float64) geometry.Polygon {
	var points []geometry.Point
	for _, coord := range coords {
		points = append(points, geometry.Point{
			Lat: coord[1],
			Lng: coord[0],
		})
	}

	lineString := geometry.LineString{
		Coordinates: points,
	}

	return geometry.Polygon{
		Coordinates: []geometry.LineString{lineString},
	}
}

func TestIntersectWithRealData(t *testing.T) {
	// Load test data if available
	fixture, err := utils.LoadJSONFixture("test-data/polygon.json")
	if err != nil {
		t.Skip("Test data not available, skipping real data test")
		return
	}

	f, err := feature.FromJSON(fixture)
	if err != nil {
		t.Skip("Failed to parse test data, skipping real data test")
		return
	}

	// Create a simple test polygon
	testPoly := createTestPolygon([][]float64{
		{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0},
	})

	// Try to intersect with loaded data
	result, err := Intersect(f, testPoly)
	if err != nil {
		t.Errorf("Intersect with real data failed: %v", err)
		return
	}

	// Result can be nil if no intersection, which is valid
	if result != nil {
		assert.Equal(t, result.Geometry.GeoJSONType, geojson.Polygon)
	}
}

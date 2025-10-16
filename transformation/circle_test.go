package transformation

import (
	"testing"

	"github.com/et-soft/turf-go/assert"
	"github.com/tomchavakis/geojson"
	"github.com/tomchavakis/geojson/feature"
	"github.com/tomchavakis/geojson/geometry"
)

func TestCircle(t *testing.T) {
	// Test basic circle creation
	center := geometry.Point{
		Lng: -75.343,
		Lat: 39.984,
	}
	radius := 5.0

	circle, err := Circle(center, radius, nil)
	if err != nil {
		t.Errorf("Circle() error = %v", err)
		return
	}

	if circle == nil {
		t.Error("Expected circle, got nil")
		return
	}

	// Verify geometry type
	assert.Equal(t, circle.Geometry.GeoJSONType, geojson.Polygon)

	// Verify bounding box
	if len(circle.Bbox) != 4 {
		t.Errorf("Expected bbox length 4, got %d", len(circle.Bbox))
	}
}

func TestCircleWithOptions(t *testing.T) {
	center := geometry.Point{
		Lng: 0,
		Lat: 0,
	}
	radius := 1.0

	options := &CircleOptions{
		Steps:      10,
		Units:      "kilometeres",
		Properties: map[string]interface{}{"name": "test circle"},
	}

	circle, err := Circle(center, radius, options)
	if err != nil {
		t.Errorf("Circle() error = %v", err)
		return
	}

	if circle == nil {
		t.Error("Expected circle, got nil")
		return
	}

	// Verify properties
	if circle.Properties["name"] != "test circle" {
		t.Errorf("Expected property 'name' to be 'test circle', got %v", circle.Properties["name"])
	}

	// Verify geometry type
	assert.Equal(t, circle.Geometry.GeoJSONType, geojson.Polygon)
}

func TestCircleFromCoordinates(t *testing.T) {
	center := []float64{-75.343, 39.984}
	radius := 5.0

	circle, err := CircleFromCoordinates(center, radius, nil)
	if err != nil {
		t.Errorf("CircleFromCoordinates() error = %v", err)
		return
	}

	if circle == nil {
		t.Error("Expected circle, got nil")
		return
	}

	assert.Equal(t, circle.Geometry.GeoJSONType, geojson.Polygon)
}

func TestCircleFromCoordinatesInvalid(t *testing.T) {
	// Test with invalid coordinates
	center := []float64{-75.343} // Only one coordinate
	radius := 5.0

	_, err := CircleFromCoordinates(center, radius, nil)
	if err == nil {
		t.Error("Expected error for invalid coordinates")
	}
}

func TestCircleFromFeature(t *testing.T) {
	// Create a point feature
	point := geometry.Point{
		Lng: -75.343,
		Lat: 39.984,
	}

	geom := geometry.Geometry{
		GeoJSONType: geojson.Point,
		Coordinates: []float64{point.Lng, point.Lat},
	}

	f, err := feature.New(geom, []float64{point.Lng, point.Lat, point.Lng, point.Lat}, nil, "")
	if err != nil {
		t.Errorf("Failed to create feature: %v", err)
		return
	}

	radius := 5.0
	circle, err := CircleFromFeature(f, radius, nil)
	if err != nil {
		t.Errorf("CircleFromFeature() error = %v", err)
		return
	}

	if circle == nil {
		t.Error("Expected circle, got nil")
		return
	}

	assert.Equal(t, circle.Geometry.GeoJSONType, geojson.Polygon)
}

func TestCircleFromFeatureInvalid(t *testing.T) {
	// Test with non-point feature
	_ = geometry.LineString{
		Coordinates: []geometry.Point{
			{Lng: 0, Lat: 0},
			{Lng: 1, Lat: 1},
		},
	}

	geom := geometry.Geometry{
		GeoJSONType: geojson.LineString,
		Coordinates: [][]interface{}{
			{[]float64{0, 0}},
			{[]float64{1, 1}},
		},
	}

	f, err := feature.New(geom, []float64{0, 0, 1, 1}, nil, "")
	if err != nil {
		t.Errorf("Failed to create feature: %v", err)
		return
	}

	_, err = CircleFromFeature(f, 5.0, nil)
	if err == nil {
		t.Error("Expected error for non-point feature")
	}
}

func TestCircleFromGeometry(t *testing.T) {
	// Create a point geometry
	point := geometry.Point{
		Lng: -75.343,
		Lat: 39.984,
	}

	geom := geometry.Geometry{
		GeoJSONType: geojson.Point,
		Coordinates: []float64{point.Lng, point.Lat},
	}

	radius := 5.0
	circle, err := CircleFromGeometry(&geom, radius, nil)
	if err != nil {
		t.Errorf("CircleFromGeometry() error = %v", err)
		return
	}

	if circle == nil {
		t.Error("Expected circle, got nil")
		return
	}

	assert.Equal(t, circle.Geometry.GeoJSONType, geojson.Polygon)
}

func TestCircleFromGeometryInvalid(t *testing.T) {
	// Test with non-point geometry
	geom := geometry.Geometry{
		GeoJSONType: geojson.LineString,
		Coordinates: [][]interface{}{
			{[]float64{0, 0}},
			{[]float64{1, 1}},
		},
	}

	_, err := CircleFromGeometry(&geom, 5.0, nil)
	if err == nil {
		t.Error("Expected error for non-point geometry")
	}
}

func TestCircleDifferentUnits(t *testing.T) {
	center := geometry.Point{
		Lng: 0,
		Lat: 0,
	}
	radius := 1.0

	// Test with different units
	units := []string{"kilometeres", "miles", "meters", "feet"}

	for _, unit := range units {
		options := &CircleOptions{
			Steps: 8,
			Units: unit,
		}

		circle, err := Circle(center, radius, options)
		if err != nil {
			t.Errorf("Circle() error for unit %s: %v", unit, err)
			continue
		}

		if circle == nil {
			t.Errorf("Expected circle for unit %s, got nil", unit)
			continue
		}

		assert.Equal(t, circle.Geometry.GeoJSONType, geojson.Polygon)
	}
}

func TestCircleSteps(t *testing.T) {
	center := geometry.Point{
		Lng: 0,
		Lat: 0,
	}
	radius := 1.0

	// Test with different step counts
	steps := []int{4, 8, 16, 32, 64}

	for _, stepCount := range steps {
		options := &CircleOptions{
			Steps: stepCount,
			Units: "kilometeres",
		}

		circle, err := Circle(center, radius, options)
		if err != nil {
			t.Errorf("Circle() error for steps %d: %v", stepCount, err)
			continue
		}

		if circle == nil {
			t.Errorf("Expected circle for steps %d, got nil", stepCount)
			continue
		}

		assert.Equal(t, circle.Geometry.GeoJSONType, geojson.Polygon)
	}
}

func TestCircleEdgeCases(t *testing.T) {
	center := geometry.Point{
		Lng: 0,
		Lat: 0,
	}

	// Test with zero radius
	circle, err := Circle(center, 0, nil)
	if err != nil {
		t.Errorf("Circle() error for zero radius: %v", err)
		return
	}

	if circle == nil {
		t.Error("Expected circle for zero radius, got nil")
		return
	}

	// Test with negative radius
	circle, err = Circle(center, -1, nil)
	if err != nil {
		t.Errorf("Circle() error for negative radius: %v", err)
		return
	}

	if circle == nil {
		t.Error("Expected circle for negative radius, got nil")
		return
	}

	// Test with very small radius
	circle, err = Circle(center, 0.001, nil)
	if err != nil {
		t.Errorf("Circle() error for small radius: %v", err)
		return
	}

	if circle == nil {
		t.Error("Expected circle for small radius, got nil")
		return
	}
}

func TestCircleNilInputs(t *testing.T) {
	_ = geometry.Point{
		Lng: 0,
		Lat: 0,
	}

	// Test with nil feature
	_, err := CircleFromFeature(nil, 5.0, nil)
	if err == nil {
		t.Error("Expected error for nil feature")
	}

	// Test with nil geometry
	_, err = CircleFromGeometry(nil, 5.0, nil)
	if err == nil {
		t.Error("Expected error for nil geometry")
	}
}

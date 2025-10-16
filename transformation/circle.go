package transformation

import (
	"errors"
	"math"

	"github.com/tomchavakis/geojson"
	"github.com/tomchavakis/geojson/feature"
	"github.com/tomchavakis/geojson/geometry"
	"github.com/tomchavakis/turf-go/constants"
	"github.com/tomchavakis/turf-go/conversions"
)

// CircleOptions contains options for the Circle function
type CircleOptions struct {
	Steps      int                    `json:"steps,omitempty"`
	Units      string                 `json:"units,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Circle creates a circle polygon from a center point and radius
func Circle(center geometry.Point, radius float64, options *CircleOptions) (*feature.Feature, error) {
	if options == nil {
		options = &CircleOptions{}
	}

	// Set default values
	if options.Steps <= 0 {
		options.Steps = 64
	}
	if options.Units == "" {
		options.Units = constants.UnitKilometers
	}

	// Convert radius to degrees for calculation
	radiusInDegrees, err := conversions.LengthToDegrees(radius, options.Units)
	if err != nil {
		return nil, err
	}

	// Generate circle coordinates
	coordinates := generateCircleCoordinates(center, radiusInDegrees, options.Steps)

	// Create polygon (not used directly, but kept for reference)
	_ = geometry.Polygon{
		Coordinates: []geometry.LineString{
			{
				Coordinates: coordinates,
			},
		},
	}

	// Create geometry
	geom := geometry.Geometry{
		GeoJSONType: geojson.Polygon,
		Coordinates: [][]interface{}{
			convertPointsToInterface(coordinates),
		},
	}

	// Calculate bounding box
	bbox := calculateCircleBoundingBox(center, radiusInDegrees)

	// Create feature
	f, err := feature.New(geom, bbox, options.Properties, "")
	if err != nil {
		return nil, err
	}

	return f, nil
}

// CircleFromCoordinates creates a circle from coordinate array [lng, lat]
func CircleFromCoordinates(center []float64, radius float64, options *CircleOptions) (*feature.Feature, error) {
	if len(center) != 2 {
		return nil, errors.New("center must be [lng, lat]")
	}

	point := geometry.Point{
		Lng: center[0],
		Lat: center[1],
	}

	return Circle(point, radius, options)
}

// generateCircleCoordinates generates coordinates for a circle
func generateCircleCoordinates(center geometry.Point, radiusInDegrees float64, steps int) []geometry.Point {
	var coordinates []geometry.Point

	for i := 0; i < steps; i++ {
		angle := float64(i) * 2 * math.Pi / float64(steps)

		// Calculate offset from center
		dx := radiusInDegrees * math.Cos(angle)
		dy := radiusInDegrees * math.Sin(angle)

		// Create point
		point := geometry.Point{
			Lng: center.Lng + dx,
			Lat: center.Lat + dy,
		}

		coordinates = append(coordinates, point)
	}

	// Close the polygon by adding the first point at the end
	if len(coordinates) > 0 {
		coordinates = append(coordinates, coordinates[0])
	}

	return coordinates
}

// convertPointsToInterface converts geometry.Point slice to interface{} slice
func convertPointsToInterface(points []geometry.Point) []interface{} {
	var result []interface{}
	for _, point := range points {
		result = append(result, []float64{point.Lng, point.Lat})
	}
	return result
}

// calculateCircleBoundingBox calculates bounding box for a circle
func calculateCircleBoundingBox(center geometry.Point, radiusInDegrees float64) []float64 {
	return []float64{
		center.Lng - radiusInDegrees, // west
		center.Lat - radiusInDegrees, // south
		center.Lng + radiusInDegrees, // east
		center.Lat + radiusInDegrees, // north
	}
}

// CircleFromFeature creates a circle from a feature point
func CircleFromFeature(f *feature.Feature, radius float64, options *CircleOptions) (*feature.Feature, error) {
	if f == nil {
		return nil, errors.New("feature cannot be nil")
	}

	if f.Geometry.GeoJSONType != geojson.Point {
		return nil, errors.New("feature must be a Point")
	}

	// Extract point from feature
	point, err := f.ToPoint()
	if err != nil {
		return nil, err
	}

	return Circle(*point, radius, options)
}

// CircleFromGeometry creates a circle from a geometry point
func CircleFromGeometry(geom *geometry.Geometry, radius float64, options *CircleOptions) (*feature.Feature, error) {
	if geom == nil {
		return nil, errors.New("geometry cannot be nil")
	}

	if geom.GeoJSONType != geojson.Point {
		return nil, errors.New("geometry must be a Point")
	}

	// Extract point from geometry
	point, err := geom.ToPoint()
	if err != nil {
		return nil, err
	}

	return Circle(*point, radius, options)
}

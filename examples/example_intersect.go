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
		log.Fatalf("Error calculating intersection: %v", err)
	}

	if intersection != nil {
		fmt.Println("Intersection found!")
		fmt.Printf("Geometry type: %s\n", intersection.Geometry.GeoJSONType)
		fmt.Printf("Bounding box: %v\n", intersection.Bbox)
	} else {
		fmt.Println("No intersection found")
	}

	// Test with non-overlapping polygons
	poly3 := createPolygon([][]float64{
		{0, 0}, {5, 0}, {5, 5}, {0, 5}, {0, 0},
	})

	poly4 := createPolygon([][]float64{
		{10, 10}, {15, 10}, {15, 15}, {10, 15}, {10, 10},
	})

	intersection2, err := transformation.Intersect(poly3, poly4)
	if err != nil {
		log.Fatalf("Error calculating intersection: %v", err)
	}

	if intersection2 != nil {
		fmt.Println("Unexpected intersection found!")
	} else {
		fmt.Println("No intersection found (as expected)")
	}
}

// Helper function to create a polygon from coordinates
func createPolygon(coords [][]float64) geometry.Polygon {
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

package main

import (
	"fmt"
	"log"

	"github.com/tomchavakis/geojson"
	"github.com/tomchavakis/geojson/feature"
	"github.com/tomchavakis/geojson/geometry"
	"github.com/tomchavakis/turf-go/transformation"
)

func main() {
	// Example 1: Basic circle creation
	fmt.Println("=== Basic Circle Creation ===")
	center := geometry.Point{
		Lng: -75.343,
		Lat: 39.984,
	}
	radius := 5.0

	circle, err := transformation.Circle(center, radius, nil)
	if err != nil {
		log.Fatalf("Error creating circle: %v", err)
	}

	fmt.Printf("Circle created with center: [%.3f, %.3f]\n", center.Lng, center.Lat)
	fmt.Printf("Radius: %.1f km\n", radius)
	fmt.Printf("Geometry type: %s\n", circle.Geometry.GeoJSONType)
	fmt.Printf("Bounding box: %v\n", circle.Bbox)

	// Example 2: Circle with custom options
	fmt.Println("\n=== Circle with Custom Options ===")
	options := &transformation.CircleOptions{
		Steps: 16,
		Units: "kilometeres",
		Properties: map[string]interface{}{
			"name":        "Custom Circle",
			"description": "A circle with 16 steps",
		},
	}

	customCircle, err := transformation.Circle(center, radius, options)
	if err != nil {
		log.Fatalf("Error creating custom circle: %v", err)
	}

	fmt.Printf("Custom circle created with %d steps\n", options.Steps)
	fmt.Printf("Properties: %v\n", customCircle.Properties)

	// Example 3: Circle from coordinates
	fmt.Println("\n=== Circle from Coordinates ===")
	coords := []float64{-75.343, 39.984}
	circleFromCoords, err := transformation.CircleFromCoordinates(coords, 2.0, nil)
	if err != nil {
		log.Fatalf("Error creating circle from coordinates: %v", err)
	}

	fmt.Printf("Circle from coordinates created\n")
	fmt.Printf("Bounding box: %v\n", circleFromCoords.Bbox)

	// Example 4: Different units
	fmt.Println("\n=== Circle with Different Units ===")
	units := []string{"kilometeres", "miles", "meters", "feet"}

	for _, unit := range units {
		unitOptions := &transformation.CircleOptions{
			Steps: 8,
			Units: unit,
		}

		unitCircle, err := transformation.Circle(center, 1.0, unitOptions)
		if err != nil {
			fmt.Printf("Error creating circle with unit %s: %v\n", unit, err)
			continue
		}

		fmt.Printf("Circle with unit '%s' created successfully\n", unit)
		fmt.Printf("  Bounding box: %v\n", unitCircle.Bbox)
	}

	// Example 5: Circle from feature
	fmt.Println("\n=== Circle from Feature ===")
	point := geometry.Point{
		Lng: 0,
		Lat: 0,
	}

	geom := geometry.Geometry{
		GeoJSONType: geojson.Point,
		Coordinates: []float64{point.Lng, point.Lat},
	}

	feature, err := feature.New(geom, []float64{point.Lng, point.Lat, point.Lng, point.Lat}, nil, "")
	if err != nil {
		log.Fatalf("Error creating feature: %v", err)
	}

	featureCircle, err := transformation.CircleFromFeature(feature, 3.0, nil)
	if err != nil {
		log.Fatalf("Error creating circle from feature: %v", err)
	}

	fmt.Printf("Circle from feature created\n")
	fmt.Printf("Bounding box: %v\n", featureCircle.Bbox)

	fmt.Println("\n=== All examples completed successfully! ===")
}

package internal

import (
	"math"
)

const DEG_LAT_P_MILE = 0.0144927536231884 // degrees latitude per mile
const DEG_LON_P_MILE = 0.0181818181818182 // degrees longitude per mile

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// DistanceOfLatLng - returns the distance (in miles) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is MILES!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func DistanceOfLatLng(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 3959 // Earth radius in MILES

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

/**
 * lat -> latitude
 * lng -> longitude
 * radius -> miles
 *
 * retrns lower latitude, longitude and higher latitude longitude
 */
func LatLngOfDistance(latitude, longitude float64, radius int) (float64, float64, float64, float64) {

	distance := float64(radius)

	lat_min := latitude - DEG_LAT_P_MILE*distance
	lon_min := longitude - DEG_LON_P_MILE*distance

	lat_max := latitude + DEG_LAT_P_MILE*distance
	lon_max := longitude + DEG_LON_P_MILE*distance

	return lat_min, lon_min, lat_max, lon_max
}

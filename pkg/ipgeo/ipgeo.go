package ipgeo

type IpLocation struct {
	IpAddress   string
	CountryCode string
	CountryName string
	Region      string
	City        string
	Latitude    float64
	Longitude   float64
}

// Requests to some free internet services, returns text body with IP
func GetIp() (string, error) {
	return getIp()
}

// Requests to some free internet services, returns json body with IP and Location
func GetIpLocation() (IpLocation, error) {
	// returns structure
	return getIpLocation()
}

// Check path specified GEO database, return error when unsuccessfull.
// If path is empty - uses standart name of file.
func CheckLocationDb(path string) error {

	return nil
}

// Load last release of database by ip2location project. Save with standart name for next uses
func LoadLocationDb() error {

	return nil
}

// Read presaved GEO database. If path is empty - uses standart name of file.
func GetLocationByIp(path string) (IpLocation, error) {
	// returns structure
	return getLocationByIp(path)
}

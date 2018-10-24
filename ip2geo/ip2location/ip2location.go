package ip2location

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"strconv"
)

// APIVersion hold the API version number from
const APIVersion string = "8.0.3"
const InvalidAddress string = "Invalid IP address."
const MissingFile string = "Invalid database file."
const NotSupported string = "This parameter is unavailable for selected data file. Please upgrade the data file."

var countryPosition = [25]uint8{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var regionPosition = [25]uint8{0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}
var cityPosition = [25]uint8{0, 0, 0, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}
var ispPosition = [25]uint8{0, 0, 3, 0, 5, 0, 7, 5, 7, 0, 8, 0, 9, 0, 9, 0, 9, 0, 9, 7, 9, 0, 9, 7, 9}
var latitudePosition = [25]uint8{0, 0, 0, 0, 0, 5, 5, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5}
var longitudePosition = [25]uint8{0, 0, 0, 0, 0, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6}
var domainPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 6, 8, 0, 9, 0, 10, 0, 10, 0, 10, 0, 10, 8, 10, 0, 10, 8, 10}
var zipcodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 0, 7, 7, 7, 0, 7, 0, 7, 7, 7, 0, 7}
var timezonePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 8, 7, 8, 8, 8, 7, 8, 0, 8, 8, 8, 0, 8}
var netSpeedPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 11, 0, 11, 8, 11, 0, 11, 0, 11, 0, 11}
var iddCodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 12, 0, 12, 0, 12, 9, 12, 0, 12}
var areaCodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 13, 0, 13, 0, 13, 10, 13, 0, 13}
var weatherStationCodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 14, 0, 14, 0, 14, 0, 14}
var weatherStationNamePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 15, 0, 15, 0, 15, 0, 15}
var mccPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 16, 0, 16, 9, 16}
var mncPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 17, 0, 17, 10, 17}
var mobileBrandPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 18, 0, 18, 11, 18}
var elevationPosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 19, 0, 19}
var usageTypePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12, 20}

var maxIPV4Range = big.NewInt(4294967295)
var maxIPV6Range = big.NewInt(0)

const countryshort uint32 = 0x00001
const countrylong uint32 = 0x00002
const region uint32 = 0x00004
const city uint32 = 0x00008
const isp uint32 = 0x00010
const latitude uint32 = 0x00020
const longitude uint32 = 0x00040
const domain uint32 = 0x00080
const zipcode uint32 = 0x00100
const timezone uint32 = 0x00200
const netspeed uint32 = 0x00400
const iddcode uint32 = 0x00800
const areacode uint32 = 0x01000
const weatherstationcode uint32 = 0x02000
const weatherstationname uint32 = 0x04000
const mcc uint32 = 0x08000
const mnc uint32 = 0x10000
const mobilebrand uint32 = 0x20000
const elevation uint32 = 0x40000
const usagetype uint32 = 0x80000
const all uint32 = countryshort | countrylong | region | city | isp | latitude | longitude | domain | zipcode | timezone | netspeed | iddcode | areacode | weatherstationcode | weatherstationname | mcc | mnc | mobilebrand | elevation | usagetype

type meta struct {
	databasetype      uint8
	databasecolumn    uint8
	databaseday       uint8
	databasemonth     uint8
	databaseyear      uint8
	ipv4databasecount uint32
	ipv4databaseaddr  uint32
	ipv6databasecount uint32
	ipv6databaseaddr  uint32
	ipv4indexbaseaddr uint32
	ipv6indexbaseaddr uint32
	ipv4columnsize    uint32
	ipv6columnsize    uint32
}

type Record struct {
	CountryShort       string
	CountryLong        string
	Region             string
	City               string
	Isp                string
	Latitude           float32
	Longitude          float32
	Domain             string
	Zipcode            string
	Timezone           string
	Netspeed           string
	Iddcode            string
	Areacode           string
	Weatherstationcode string
	Weatherstationname string
	Mcc                string
	Mnc                string
	Mobilebrand        string
	Elevation          float32
	Usagetype          string
}

func (x *Record) Printrecord() {
	fmt.Printf("country_short: %s\n", x.CountryShort)
	fmt.Printf("country_long: %s\n", x.CountryLong)
	fmt.Printf("region: %s\n", x.Region)
	fmt.Printf("city: %s\n", x.City)
	fmt.Printf("isp: %s\n", x.Isp)
	fmt.Printf("latitude: %f\n", x.Latitude)
	fmt.Printf("longitude: %f\n", x.Longitude)
	fmt.Printf("domain: %s\n", x.Domain)
	fmt.Printf("zipcode: %s\n", x.Zipcode)
	fmt.Printf("timezone: %s\n", x.Timezone)
	fmt.Printf("netspeed: %s\n", x.Netspeed)
	fmt.Printf("iddcode: %s\n", x.Iddcode)
	fmt.Printf("areacode: %s\n", x.Areacode)
	fmt.Printf("weatherstationcode: %s\n", x.Weatherstationcode)
	fmt.Printf("weatherstationname: %s\n", x.Weatherstationname)
	fmt.Printf("mcc: %s\n", x.Mcc)
	fmt.Printf("mnc: %s\n", x.Mnc)
	fmt.Printf("mobilebrand: %s\n", x.Mobilebrand)
	fmt.Printf("elevation: %f\n", x.Elevation)
	fmt.Printf("usagetype: %s\n", x.Usagetype)
}

type Database struct {
	f      *bytes.Reader
	meta   meta
	metaok bool

	country_position_offset            uint32
	region_position_offset             uint32
	city_position_offset               uint32
	isp_position_offset                uint32
	domain_position_offset             uint32
	zipcode_position_offset            uint32
	latitude_position_offset           uint32
	longitude_position_offset          uint32
	timezone_position_offset           uint32
	netspeed_position_offset           uint32
	iddcode_position_offset            uint32
	areacode_position_offset           uint32
	weatherstationcode_position_offset uint32
	weatherstationname_position_offset uint32
	mcc_position_offset                uint32
	mnc_position_offset                uint32
	mobilebrand_position_offset        uint32
	elevation_position_offset          uint32
	usagetype_position_offset          uint32

	country_enabled            bool
	region_enabled             bool
	city_enabled               bool
	isp_enabled                bool
	domain_enabled             bool
	zipcode_enabled            bool
	latitude_enabled           bool
	longitude_enabled          bool
	timezone_enabled           bool
	netspeed_enabled           bool
	iddcode_enabled            bool
	areacode_enabled           bool
	weatherstationcode_enabled bool
	weatherstationname_enabled bool
	mcc_enabled                bool
	mnc_enabled                bool
	mobilebrand_enabled        bool
	elevation_enabled          bool
	usagetype_enabled          bool
}

// get IP type and calculate IP number; calculates index too if exists
func (db *Database) checkip(ip string) (iptype uint32, ipnum *big.Int, ipindex uint32) {
	iptype = 0
	ipnum = big.NewInt(0)
	ipnumtmp := big.NewInt(0)
	ipindex = 0
	ipaddress := net.ParseIP(ip)

	if ipaddress != nil {
		v4 := ipaddress.To4()

		if v4 != nil {
			iptype = 4
			ipnum.SetBytes(v4)
		} else {
			v6 := ipaddress.To16()

			if v6 != nil {
				iptype = 6
				ipnum.SetBytes(v6)
			}
		}
	}
	if iptype == 4 {
		if db.meta.ipv4indexbaseaddr > 0 {
			ipnumtmp.Rsh(ipnum, 16)
			ipnumtmp.Lsh(ipnumtmp, 3)
			ipindex = uint32(ipnumtmp.Add(ipnumtmp, big.NewInt(int64(db.meta.ipv4indexbaseaddr))).Uint64())
		}
	} else if iptype == 6 {
		if db.meta.ipv6indexbaseaddr > 0 {
			ipnumtmp.Rsh(ipnum, 112)
			ipnumtmp.Lsh(ipnumtmp, 3)
			ipindex = uint32(ipnumtmp.Add(ipnumtmp, big.NewInt(int64(db.meta.ipv6indexbaseaddr))).Uint64())
		}
	}
	return
}

// read byte
func (db *Database) readuint8(pos int64) uint8 {
	var retval uint8
	data := make([]byte, 1)
	_, err := db.f.ReadAt(data, pos-1)
	if err != nil {
		fmt.Println("File read failed:", err)
	}
	retval = data[0]
	return retval
}

// read unsigned 32-bit integer
func (db *Database) readuint32(pos uint32) uint32 {
	pos2 := int64(pos)
	var retval uint32
	data := make([]byte, 4)
	_, err := db.f.ReadAt(data, pos2-1)
	if err != nil {
		fmt.Println("File read failed:", err)
	}
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, &retval)
	if err != nil {
		fmt.Println("Binary read failed:", err)
	}
	return retval
}

// read unsigned 128-bit integer
func (db *Database) readuint128(pos uint32) *big.Int {
	pos2 := int64(pos)
	retval := big.NewInt(0)
	data := make([]byte, 16)
	_, err := db.f.ReadAt(data, pos2-1)
	if err != nil {
		fmt.Println("File read failed:", err)
	}

	// little endian to big endian
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	retval.SetBytes(data)
	return retval
}

// read string
func (db *Database) readstr(pos uint32) string {
	pos2 := int64(pos)
	var retval string
	lenbyte := make([]byte, 1)
	_, err := db.f.ReadAt(lenbyte, pos2)
	if err != nil {
		fmt.Println("File read failed:", err)
	}
	strlen := lenbyte[0]
	data := make([]byte, strlen)
	_, err = db.f.ReadAt(data, pos2+1)
	if err != nil {
		fmt.Println("File read failed:", err)
	}
	retval = string(data[:strlen])
	return retval
}

// read float
func (db *Database) readfloat(pos uint32) float32 {
	pos2 := int64(pos)
	var retval float32
	data := make([]byte, 4)
	_, err := db.f.ReadAt(data, pos2-1)
	if err != nil {
		fmt.Println("File read failed:", err)
	}
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, &retval)
	if err != nil {
		fmt.Println("Binary read failed:", err)
	}
	return retval
}

func OpenFromFile(dbpath string) (*Database, error) {
	fil, err := os.Open(dbpath)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(fil)
	if err != nil {
		return nil, err
	}
	return OpenFromBytes(b)
}

func OpenFromBytes(data []byte) (*Database, error) {
	var err error
	db := &Database{}
	db.f = bytes.NewReader(data)
	if err != nil {
		return nil, err
	}
	db.initialize()
	return db, nil
}

// initialize the component with the database path
func (db *Database) initialize() {
	maxIPV6Range.SetString("340282366920938463463374607431768211455", 10)

	db.meta.databasetype = db.readuint8(1)
	db.meta.databasecolumn = db.readuint8(2)
	db.meta.databaseyear = db.readuint8(3)
	db.meta.databasemonth = db.readuint8(4)
	db.meta.databaseday = db.readuint8(5)
	db.meta.ipv4databasecount = db.readuint32(6)
	db.meta.ipv4databaseaddr = db.readuint32(10)
	db.meta.ipv6databasecount = db.readuint32(14)
	db.meta.ipv6databaseaddr = db.readuint32(18)
	db.meta.ipv4indexbaseaddr = db.readuint32(22)
	db.meta.ipv6indexbaseaddr = db.readuint32(26)
	db.meta.ipv4columnsize = uint32(db.meta.databasecolumn << 2)              // 4 bytes each column
	db.meta.ipv6columnsize = uint32(16 + ((db.meta.databasecolumn - 1) << 2)) // 4 bytes each column, except IPFrom column which is 16 bytes

	dbt := db.meta.databasetype

	// since both IPv4 and IPv6 use 4 bytes for the below columns, can just do it once here
	if countryPosition[dbt] != 0 {
		db.country_position_offset = uint32(countryPosition[dbt]-1) << 2
		db.country_enabled = true
	}
	if regionPosition[dbt] != 0 {
		db.region_position_offset = uint32(regionPosition[dbt]-1) << 2
		db.region_enabled = true
	}
	if cityPosition[dbt] != 0 {
		db.city_position_offset = uint32(cityPosition[dbt]-1) << 2
		db.city_enabled = true
	}
	if ispPosition[dbt] != 0 {
		db.isp_position_offset = uint32(ispPosition[dbt]-1) << 2
		db.isp_enabled = true
	}
	if domainPosition[dbt] != 0 {
		db.domain_position_offset = uint32(domainPosition[dbt]-1) << 2
		db.domain_enabled = true
	}
	if zipcodePosition[dbt] != 0 {
		db.zipcode_position_offset = uint32(zipcodePosition[dbt]-1) << 2
		db.zipcode_enabled = true
	}
	if latitudePosition[dbt] != 0 {
		db.latitude_position_offset = uint32(latitudePosition[dbt]-1) << 2
		db.latitude_enabled = true
	}
	if longitudePosition[dbt] != 0 {
		db.longitude_position_offset = uint32(longitudePosition[dbt]-1) << 2
		db.longitude_enabled = true
	}
	if timezonePosition[dbt] != 0 {
		db.timezone_position_offset = uint32(timezonePosition[dbt]-1) << 2
		db.timezone_enabled = true
	}
	if netSpeedPosition[dbt] != 0 {
		db.netspeed_position_offset = uint32(netSpeedPosition[dbt]-1) << 2
		db.netspeed_enabled = true
	}
	if iddCodePosition[dbt] != 0 {
		db.iddcode_position_offset = uint32(iddCodePosition[dbt]-1) << 2
		db.iddcode_enabled = true
	}
	if areaCodePosition[dbt] != 0 {
		db.areacode_position_offset = uint32(areaCodePosition[dbt]-1) << 2
		db.areacode_enabled = true
	}
	if weatherStationCodePosition[dbt] != 0 {
		db.weatherstationcode_position_offset = uint32(weatherStationCodePosition[dbt]-1) << 2
		db.weatherstationcode_enabled = true
	}
	if weatherStationNamePosition[dbt] != 0 {
		db.weatherstationname_position_offset = uint32(weatherStationNamePosition[dbt]-1) << 2
		db.weatherstationname_enabled = true
	}
	if mccPosition[dbt] != 0 {
		db.mcc_position_offset = uint32(mccPosition[dbt]-1) << 2
		db.mcc_enabled = true
	}
	if mncPosition[dbt] != 0 {
		db.mnc_position_offset = uint32(mncPosition[dbt]-1) << 2
		db.mnc_enabled = true
	}
	if mobileBrandPosition[dbt] != 0 {
		db.mobilebrand_position_offset = uint32(mobileBrandPosition[dbt]-1) << 2
		db.mobilebrand_enabled = true
	}
	if elevationPosition[dbt] != 0 {
		db.elevation_position_offset = uint32(elevationPosition[dbt]-1) << 2
		db.elevation_enabled = true
	}
	if usageTypePosition[dbt] != 0 {
		db.usagetype_position_offset = uint32(usageTypePosition[dbt]-1) << 2
		db.usagetype_enabled = true
	}

	db.metaok = true
}

// close database file handle
func (db *Database) Close() {
	db.f = nil
}

// populate record with message
func loadmessage(mesg string) Record {
	var x Record

	x.CountryShort = mesg
	x.CountryLong = mesg
	x.Region = mesg
	x.City = mesg
	x.Isp = mesg
	x.Domain = mesg
	x.Zipcode = mesg
	x.Timezone = mesg
	x.Netspeed = mesg
	x.Iddcode = mesg
	x.Areacode = mesg
	x.Weatherstationcode = mesg
	x.Weatherstationname = mesg
	x.Mcc = mesg
	x.Mnc = mesg
	x.Mobilebrand = mesg
	x.Usagetype = mesg

	return x
}

// get all fields
func (db *Database) Get_all(ipaddress string) Record {
	return db.query(ipaddress, all)
}

// get country code
func (db *Database) Get_country_short(ipaddress string) Record {
	return db.query(ipaddress, countryshort)
}

// get country name
func (db *Database) Get_country_long(ipaddress string) Record {
	return db.query(ipaddress, countrylong)
}

// get region
func (db *Database) Get_region(ipaddress string) Record {
	return db.query(ipaddress, region)
}

// get city
func (db *Database) Get_city(ipaddress string) Record {
	return db.query(ipaddress, city)
}

// get isp
func (db *Database) Get_isp(ipaddress string) Record {
	return db.query(ipaddress, isp)
}

// get latitude
func (db *Database) Get_latitude(ipaddress string) Record {
	return db.query(ipaddress, latitude)
}

// get longitude
func (db *Database) Get_longitude(ipaddress string) Record {
	return db.query(ipaddress, longitude)
}

// get domain
func (db *Database) Get_domain(ipaddress string) Record {
	return db.query(ipaddress, domain)
}

// get zip code
func (db *Database) Get_zipcode(ipaddress string) Record {
	return db.query(ipaddress, zipcode)
}

// get time zone
func (db *Database) Get_timezone(ipaddress string) Record {
	return db.query(ipaddress, timezone)
}

// get net speed
func (db *Database) Get_netspeed(ipaddress string) Record {
	return db.query(ipaddress, netspeed)
}

// get idd code
func (db *Database) Get_iddcode(ipaddress string) Record {
	return db.query(ipaddress, iddcode)
}

// get area code
func (db *Database) Get_areacode(ipaddress string) Record {
	return db.query(ipaddress, areacode)
}

// get weather station code
func (db *Database) Get_weatherstationcode(ipaddress string) Record {
	return db.query(ipaddress, weatherstationcode)
}

// get weather station name
func (db *Database) Get_weatherstationname(ipaddress string) Record {
	return db.query(ipaddress, weatherstationname)
}

// get mobile country code
func (db *Database) Get_mcc(ipaddress string) Record {
	return db.query(ipaddress, mcc)
}

// get mobile network code
func (db *Database) Get_mnc(ipaddress string) Record {
	return db.query(ipaddress, mnc)
}

// get mobile carrier brand
func (db *Database) Get_mobilebrand(ipaddress string) Record {
	return db.query(ipaddress, mobilebrand)
}

// get elevation
func (db *Database) Get_elevation(ipaddress string) Record {
	return db.query(ipaddress, elevation)
}

// get usage type
func (db *Database) Get_usagetype(ipaddress string) Record {
	return db.query(ipaddress, usagetype)
}

// main query
func (db *Database) query(ipaddress string, mode uint32) Record {
	x := loadmessage(NotSupported) // default message

	// read metadata
	if db == nil || !db.metaok {
		x = loadmessage(MissingFile)
		return x
	}

	// check IP type and return IP number & index (if exists)
	iptype, ipno, ipindex := db.checkip(ipaddress)

	if iptype == 0 {
		x = loadmessage(InvalidAddress)
		return x
	}

	var colsize uint32
	var baseaddr uint32
	var low uint32
	var high uint32
	var mid uint32
	var rowoffset uint32
	var rowoffset2 uint32
	ipfrom := big.NewInt(0)
	ipto := big.NewInt(0)
	maxip := big.NewInt(0)

	if iptype == 4 {
		baseaddr = db.meta.ipv4databaseaddr
		high = db.meta.ipv4databasecount
		maxip = maxIPV4Range
		colsize = db.meta.ipv4columnsize
	} else {
		baseaddr = db.meta.ipv6databaseaddr
		high = db.meta.ipv6databasecount
		maxip = maxIPV6Range
		colsize = db.meta.ipv6columnsize
	}

	// reading index
	if ipindex > 0 {
		low = db.readuint32(ipindex)
		high = db.readuint32(ipindex + 4)
	}

	if ipno.Cmp(maxip) >= 0 {
		ipno = ipno.Sub(ipno, big.NewInt(1))
	}

	for low <= high {
		mid = ((low + high) >> 1)
		rowoffset = baseaddr + (mid * colsize)
		rowoffset2 = rowoffset + colsize

		if iptype == 4 {
			ipfrom = big.NewInt(int64(db.readuint32(rowoffset)))
			ipto = big.NewInt(int64(db.readuint32(rowoffset2)))
		} else {
			ipfrom = db.readuint128(rowoffset)
			ipto = db.readuint128(rowoffset2)
		}

		if ipno.Cmp(ipfrom) >= 0 && ipno.Cmp(ipto) < 0 {
			if iptype == 6 {
				rowoffset = rowoffset + 12 // coz below is assuming all columns are 4 bytes, so got 12 left to go to make 16 bytes total
			}

			if mode&countryshort == 1 && db.country_enabled {
				x.CountryShort = db.readstr(db.readuint32(rowoffset + db.country_position_offset))
			}

			if mode&countrylong != 0 && db.country_enabled {
				x.CountryLong = db.readstr(db.readuint32(rowoffset+db.country_position_offset) + 3)
			}

			if mode&region != 0 && db.region_enabled {
				x.Region = db.readstr(db.readuint32(rowoffset + db.region_position_offset))
			}

			if mode&city != 0 && db.city_enabled {
				x.City = db.readstr(db.readuint32(rowoffset + db.city_position_offset))
			}

			if mode&isp != 0 && db.isp_enabled {
				x.Isp = db.readstr(db.readuint32(rowoffset + db.isp_position_offset))
			}

			if mode&latitude != 0 && db.latitude_enabled {
				x.Latitude = db.readfloat(rowoffset + db.latitude_position_offset)
			}

			if mode&longitude != 0 && db.longitude_enabled {
				x.Longitude = db.readfloat(rowoffset + db.longitude_position_offset)
			}

			if mode&domain != 0 && db.domain_enabled {
				x.Domain = db.readstr(db.readuint32(rowoffset + db.domain_position_offset))
			}

			if mode&zipcode != 0 && db.zipcode_enabled {
				x.Zipcode = db.readstr(db.readuint32(rowoffset + db.zipcode_position_offset))
			}

			if mode&timezone != 0 && db.timezone_enabled {
				x.Timezone = db.readstr(db.readuint32(rowoffset + db.timezone_position_offset))
			}

			if mode&netspeed != 0 && db.netspeed_enabled {
				x.Netspeed = db.readstr(db.readuint32(rowoffset + db.netspeed_position_offset))
			}

			if mode&iddcode != 0 && db.iddcode_enabled {
				x.Iddcode = db.readstr(db.readuint32(rowoffset + db.iddcode_position_offset))
			}

			if mode&areacode != 0 && db.areacode_enabled {
				x.Areacode = db.readstr(db.readuint32(rowoffset + db.areacode_position_offset))
			}

			if mode&weatherstationcode != 0 && db.weatherstationcode_enabled {
				x.Weatherstationcode = db.readstr(db.readuint32(rowoffset + db.weatherstationcode_position_offset))
			}

			if mode&weatherstationname != 0 && db.weatherstationname_enabled {
				x.Weatherstationname = db.readstr(db.readuint32(rowoffset + db.weatherstationname_position_offset))
			}

			if mode&mcc != 0 && db.mcc_enabled {
				x.Mcc = db.readstr(db.readuint32(rowoffset + db.mcc_position_offset))
			}

			if mode&mnc != 0 && db.mnc_enabled {
				x.Mnc = db.readstr(db.readuint32(rowoffset + db.mnc_position_offset))
			}

			if mode&mobilebrand != 0 && db.mobilebrand_enabled {
				x.Mobilebrand = db.readstr(db.readuint32(rowoffset + db.mobilebrand_position_offset))
			}

			if mode&elevation != 0 && db.elevation_enabled {
				f, _ := strconv.ParseFloat(db.readstr(db.readuint32(rowoffset+db.elevation_position_offset)), 32)
				x.Elevation = float32(f)
			}

			if mode&usagetype != 0 && db.usagetype_enabled {
				x.Usagetype = db.readstr(db.readuint32(rowoffset + db.usagetype_position_offset))
			}

			return x
		}
		if ipno.Cmp(ipfrom) < 0 {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return x
}

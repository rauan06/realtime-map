package domain

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	SRID = 4326
)

type GeoPointDB struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

type GeoPoint struct {
	GeoPointDB
	SRID int `json:"srid"` // default 4326
}

func (p *GeoPoint) SetSRID(srid int) *GeoPoint {
	p.SRID = srid

	return p
}

func NewGeoPoint(lng, lat float64) *GeoPoint {
	return &GeoPoint{
		GeoPointDB: GeoPointDB{
			Lng: lng,
			Lat: lat,
		},
		SRID: SRID,
	}
}

func NewGeoPointWithSRID(lng, lat float64, srid int) *GeoPoint {
	return &GeoPoint{
		GeoPointDB: GeoPointDB{
			Lng: lng,
			Lat: lat,
		},
		SRID: srid,
	}
}

func (p *GeoPoint) String() string {
	/* SRID=%d;POINT(%v %v) */
	var builder strings.Builder
	builder.WriteString("SRID=")
	builder.WriteString(strconv.Itoa(p.SRID))
	builder.WriteString(";POINT(")
	builder.WriteString(strconv.FormatFloat(p.Lng, 'g', -1, 64))
	builder.WriteByte(' ')
	builder.WriteString(strconv.FormatFloat(p.Lat, 'g', -1, 64))
	builder.WriteByte(')')

	return builder.String()
}

func (p *GeoPointDB) Scan(val interface{}) error {
	if val == nil {
		return nil
	}

	var b []byte

	var err error

	switch v := val.(type) {
	case []uint8:
		b, err = hex.DecodeString(string(v))
	case string:
		b, err = hex.DecodeString(v)
	default:
		return errors.New("invalid type")
	}

	if err != nil {
		return err
	}

	r := bytes.NewReader(b)

	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder

	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return errors.New("invalid byte order")
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

func (p GeoPoint) Value() (driver.Value, error) {
	return p.String(), nil
}

type OBUData struct {
	SessionID uuid.UUID `gorm:"type:uuid;primaryKey" json:"session_id"`
	Point     GeoPoint  `gorm:"type:geometry(Point,4326)" json:"point"`
	CreatedAt time.Time `gorm:"type:timestamptz" json:"created_at"`
}

func (OBUData) TableName() string { return "obu_data" }

package opensearch

import "github.com/knives85/hotel-search/internal/domain"

// hotelSearchDocument is the JSON shape of a hit's _source in the hotels
// index. Fields not relevant to the search projection are omitted (giata_id,
// internal_city_id, provider_ids).
type hotelSearchDocument struct {
	UniqueID        int64                `json:"unique_id"`
	IndexStatus     string               `json:"index_status"`
	HotelName       *string              `json:"hotel_name"`
	SellStatus      *bool                `json:"sell_status"`
	StarRating      *string              `json:"star_rating"`
	Type            *string              `json:"type"`
	Country         *geoDocument         `json:"country"`
	City            *geoDocument         `json:"city"`
	AdminRegion     *geoDocument         `json:"admin_region"`
	TouristicRegion *geoDocument         `json:"touristic_region"`
	NonAdminCity    *geoDocument         `json:"nonadmin_city"`
	Neighbourhood   *geoDocument         `json:"neighbourhood"`
	Chain           *chainDocument       `json:"chain"`
	ContentScore    *int                 `json:"content_score"`
	ReviewScore     *int                 `json:"review_score"`
	NumberOfReviews *int                 `json:"number_of_reviews"`
	LocationScore   *int                 `json:"location_score"`
	LastUpdateDate  *int64               `json:"last_update_date"`
	CreationDate    *int64               `json:"creation_date"`
	FacilityCodes   []string             `json:"facility_codes"`
	NearbyPois      []geoDocument        `json:"nearby_pois"`
	Badges          []string             `json:"badges"`
	Coordinates     *coordinatesDocument `json:"coordinates"`
}

type geoDocument struct {
	Code string  `json:"code"`
	Name *string `json:"name"`
}

type chainDocument struct {
	Code string  `json:"code"`
	Name *string `json:"name"`
}

type coordinatesDocument struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

// toHotel maps the raw search document to the domain projection. Unknown
// IndexStatus values are coerced to PARTIAL, matching the legacy behaviour.
func toHotel(d hotelSearchDocument) domain.Hotel {
	return domain.Hotel{
		UniqueID:         d.UniqueID,
		IndexStatus:      parseIndexStatus(d.IndexStatus),
		HotelName:        d.HotelName,
		SellStatus:       d.SellStatus,
		StarRating:       d.StarRating,
		Type:             d.Type,
		Country:          d.Country.toRef(),
		City:             d.City.toRef(),
		AdminRegion:      d.AdminRegion.toRef(),
		TouristicRegion:  d.TouristicRegion.toRef(),
		NonAdminCity:     d.NonAdminCity.toRef(),
		Neighbourhood:    d.Neighbourhood.toRef(),
		Chain:            d.Chain.toRef(),
		Facilities:       d.FacilityCodes,
		PointsOfInterest: toGeoRefs(d.NearbyPois),
		ContentScore:     d.ContentScore,
		ReviewScore:      d.ReviewScore,
		NumberOfReviews:  d.NumberOfReviews,
		LocationScore:    d.LocationScore,
		Badges:           d.Badges,
		CreationDate:     d.CreationDate,
		LastUpdateDate:   d.LastUpdateDate,
		Coordinates:      d.Coordinates.toCoordinates(),
	}
}

func parseIndexStatus(s string) domain.IndexStatus {
	switch domain.IndexStatus(s) {
	case domain.IndexStatusComplete:
		return domain.IndexStatusComplete
	default:
		return domain.IndexStatusPartial
	}
}

func (g *geoDocument) toRef() *domain.GeoReference {
	if g == nil {
		return nil
	}
	return &domain.GeoReference{Code: g.Code, Name: g.Name}
}

func (c *chainDocument) toRef() *domain.ChainReference {
	if c == nil {
		return nil
	}
	return &domain.ChainReference{Code: c.Code, Name: c.Name}
}

func (c *coordinatesDocument) toCoordinates() *domain.Coordinates {
	if c == nil || c.Latitude == nil || c.Longitude == nil {
		return nil
	}
	return &domain.Coordinates{Latitude: *c.Latitude, Longitude: *c.Longitude}
}

func toGeoRefs(docs []geoDocument) []domain.GeoReference {
	if len(docs) == 0 {
		return nil
	}
	out := make([]domain.GeoReference, len(docs))
	for i, d := range docs {
		out[i] = domain.GeoReference{Code: d.Code, Name: d.Name}
	}
	return out
}

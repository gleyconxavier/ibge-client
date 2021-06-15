package client

import (
	"errors"
	"fmt"
	"github.com/eucatur/go-toolbox/check"
	"github.com/parnurzeal/gorequest"
	geo "github.com/paulmach/go.geo"
	"sort"
	"strings"
)

func (c *Client) request() (superAgent *gorequest.SuperAgent) {
	if c.superAgent != nil {
		return c.superAgent
	}

	c.superAgent = gorequest.New()
	return c.superAgent
}

func (c Client) get(targetURL string, isGeocodeRequest bool) (superAgent *gorequest.SuperAgent) {
	if isGeocodeRequest {
		return c.request().Get(c.HostGeocode + targetURL)
	}
	return c.request().Get(c.Host + targetURL)
}

func (c *Client) GetStates() (states []State, err error) {
	states = []State{}

	_, _, errs := c.get("/localidades/estados", false).EndStruct(&states, logFile)
	if len(errs) > 0 {
		err = errs[0]
		return
	}

	return
}

func (c *Client) GetCounties() (counties []County, err error) {
	counties = []County{}

	_, _, errs := c.get("/localidades/municipios?orderBy=nome", false).EndStruct(&counties, logFile)
	if len(errs) > 0 {
		err = errs[0]
		return
	}

	return
}

func (c *Client) GetCountiesByUF(codState int64) (counties []County, err error) {
	counties = []County{}

	_, _, errs := c.get(fmt.Sprintf("/localidades/estados/%d/municipios?orderBy=nome", int(codState)), false).EndStruct(&counties, logFile)
	if len(errs) > 0 {
		err = errs[0]
		return
	}

	return
}

func (c *Client) GetCountyByIbgeCode(ibgeCode int64) (county County, err error) {
	if ibgeCode == 0 {
		err = errors.New("it is mandatory to inform the ibge code to perform this query")
		return
	}

	_, _, errs := c.get(fmt.Sprintf("/localidades/municipios/%d", ibgeCode), false).EndStruct(&county, logFile)
	if len(errs) > 0 {
		err = errs[0]
		return
	}

	return
}

func (c *Client) GetCountyByAcronymStateAndNameCounty(acronymState, name string, ibge_code int64) (county County, equivalencePercentagem int64, err error) {

	var codState int64
	county = County{}
	equivalencePercentagem = 0

	if ibge_code > 0 {
		county, err = c.GetCountyByIbgeCode(ibge_code)
		if err != nil {
			return
		}

		if county.ID > 0 {
			equivalencePercentagem = 100
		}

		if c.HostGeocode != "" && c.KeyGeocode != "" {
			county.Point, err = c.GetGeocode(county.Name, county.MicroRegion.MesorRegion.State.Acronym)
		}

		return
	}

	states, err := c.GetStates()
	if err != nil {
		return
	}

	for _, state := range states {
		if state.Acronym == acronymState {
			codState = state.ID
		}
	}

	if codState < 1 {
		err = errors.New("it was not possible to find the state by the abbreviation informed")
		return
	}

	counties, err := c.GetCountiesByUF(codState)

	if err != nil {
		return
	}

	if len(counties) == 0 {
		err = errors.New("it was not possible to get municipalities through the state code informed")
		return
	}

	for _, coun := range counties {
		if coun.Name == name {
			county = coun
			if c.HostGeocode != "" && c.KeyGeocode != "" {
				county.Point, err = c.GetGeocode(county.Name, county.MicroRegion.MesorRegion.State.Acronym)
			}
			return
		}

		nameCounty := strings.Split(coun.Name, " ")
		sort.Strings(nameCounty)

		nameParam := strings.Split(name, " ")
		sort.Strings(nameParam)

		mirrorNameSlice := check.If(len(nameCounty) > len(nameParam), len(nameParam), len(nameCounty)).(int)

		var countEquals = 0

		for i := 0; i < mirrorNameSlice; i++ {
			if len(nameCounty[i]) < len(nameParam[i]) {
				if strings.Contains(nameParam[i], nameCounty[i]) {
					countEquals++
				}
			} else {
				if strings.Contains(nameCounty[i], nameParam[i]) {
					countEquals++
				}
			}
		}

		percecentagem := (countEquals / len(coun.Name)) * 100

		if percecentagem > 90 {
			county = coun
			equivalencePercentagem = int64(percecentagem)
			county.Point, err = c.GetGeocode(county.Name, county.MicroRegion.MesorRegion.State.Acronym)
			return
		}
	}
	err = errors.New("it was not possible to find any municipality with the information provided")
	return
}

func (c Client) GetGeocode(acronymState, name string) (point *geo.Point, err error) {
	tempStruct := struct {
		Info struct {
			Statuscode int `json:"statuscode"`
			Copyright  struct {
				Text         string `json:"text"`
				ImageURL     string `json:"imageUrl"`
				ImageAltText string `json:"imageAltText"`
			} `json:"copyright"`
			Messages []interface{} `json:"messages"`
		} `json:"info"`
		Options struct {
			MaxResults        int  `json:"maxResults"`
			ThumbMaps         bool `json:"thumbMaps"`
			IgnoreLatLngInput bool `json:"ignoreLatLngInput"`
		} `json:"options"`
		Results []struct {
			ProvidedLocation struct {
				Location string `json:"location"`
			} `json:"providedLocation"`
			Locations []struct {
				Street             string `json:"street"`
				AdminArea6         string `json:"adminArea6"`
				AdminArea6Type     string `json:"adminArea6Type"`
				AdminArea5         string `json:"adminArea5"`
				AdminArea5Type     string `json:"adminArea5Type"`
				AdminArea4         string `json:"adminArea4"`
				AdminArea4Type     string `json:"adminArea4Type"`
				AdminArea3         string `json:"adminArea3"`
				AdminArea3Type     string `json:"adminArea3Type"`
				AdminArea1         string `json:"adminArea1"`
				AdminArea1Type     string `json:"adminArea1Type"`
				PostalCode         string `json:"postalCode"`
				GeocodeQualityCode string `json:"geocodeQualityCode"`
				GeocodeQuality     string `json:"geocodeQuality"`
				DragPoint          bool   `json:"dragPoint"`
				SideOfStreet       string `json:"sideOfStreet"`
				LinkID             string `json:"linkId"`
				UnknownInput       string `json:"unknownInput"`
				Type               string `json:"type"`
				LatLng             struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"latLng"`
				DisplayLatLng struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"displayLatLng"`
				MapURL string `json:"mapUrl"`
			} `json:"locations"`
		} `json:"results"`
	}{}

	_, _, errs := c.get(fmt.Sprintf("address?key=%s&location=%s,%s", c.KeyGeocode, name, acronymState), true).EndStruct(&tempStruct, logFile)
	if len(errs) > 0 {
		err = errs[0]
		return
	}

	if len(tempStruct.Results) < 1 || len(tempStruct.Results[0].Locations) < 1 {
		err = errors.New("it was not possible to obtain the location through the data informed")
	}

	point = geo.NewPoint(tempStruct.Results[0].Locations[0].LatLng.Lat, tempStruct.Results[0].Locations[0].LatLng.Lng)

	return
}

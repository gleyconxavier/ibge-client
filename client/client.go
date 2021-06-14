package client

import (
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

func (c *Client) request() (superAgent *gorequest.SuperAgent) {
	if c.superAgent != nil {
		return c.superAgent
	}

	c.superAgent = gorequest.New()
	return c.superAgent
}

func (c Client) get(targetURL string) (superAgent *gorequest.SuperAgent) {
	return c.request().Get(c.Host + targetURL)
}

func (c *Client) GetStates() (states []State, err error) {
	states = []State{}

	_, _, errs := c.get("/localidades/estados").EndStruct(&states, logFile)
	if len(errs) > 0 {
		err = errs[0]
		return
	}

	return
}

func (c *Client) GetCounties() (counties []County, err error) {
	counties = []County{}

	_, _, errs := c.get("/localidades/municipios?orderBy=nome").EndStruct(&counties, logFile)
	if len(errs) > 0 {
		err = errs[0]
		return
	}

	return
}

func (c *Client) GetCountiesByUF(codState int64) (counties []County, err error) {
	counties = []County{}

	_, _, errs := c.get(fmt.Sprintf("/localidades/estados/%d/municipios?orderBy=nome", int(codState))).EndStruct(&counties, logFile)
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

	_, _, errs := c.get(fmt.Sprintf("/localidades/municipios/%d", ibgeCode)).EndStruct(&county, logFile)
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
			return
		}
		if len(coun.Name) < len(name) {
			var countEquals = 0
			for i := 0; i < len(coun.Name); i++ {
				if coun.Name[i] == name[i] {
					countEquals++
				}
			}
			percecentagem := (countEquals / len(coun.Name)) * 100

			if percecentagem > 90 {
				county = coun
				equivalencePercentagem = int64(percecentagem)
				return
			}
		} else {
			var countEquals = 0
			for i := 0; i < len(name); i++ {
				if coun.Name[i] == name[i] {
					countEquals++
				}
			}
			percecentagem := (countEquals / len(name)) * 100

			if percecentagem > 90 {
				county = coun
				equivalencePercentagem = int64(percecentagem)
				return
			}
		}
	}
	err = errors.New("it was not possible to find any municipality with the information provided")
	return
}

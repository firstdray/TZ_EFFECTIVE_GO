package enrichment

import (
	"effective/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Enricher interface {
	Enrich(person *model.Person) error
}

type ApiEnricher struct{}

func NewEnricher() Enricher {
	return &ApiEnricher{}
}

func (e *ApiEnricher) Enrich(person *model.Person) error {
	if err := e.enrichAge(person); err != nil {
		return err
	}

	if err := e.enrichGender(person); err != nil {
		return err
	}

	if err := e.enrichNatioanlity(person); err != nil {
		return err
	}
	return nil
}

func (e *ApiEnricher) enrichAge(person *model.Person) error {
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", person.Name))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var result struct {
		Age int `json:"age"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	person.Age = result.Age
	return nil
}

func (e *ApiEnricher) enrichGender(person *model.Person) error {
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", person.Name))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var result struct {
		Gender string `json:"gender"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	person.Gender = result.Gender
	return nil
}

func (e *ApiEnricher) enrichNatioanlity(person *model.Person) error {
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", person.Name))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var result struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if len(result.Country) > 0 {
		person.National = result.Country[0].CountryID
	} else {
		person.National = "unknown"
	}
	return nil
}

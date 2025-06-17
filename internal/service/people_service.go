package service

import (
	"effective/internal/enrichment"
	"effective/internal/model"
	"effective/internal/repository"
	"fmt"
	"log"
	"strconv"
)

type PeopleService interface {
	AddPerson(input model.InputPerson) (*model.Person, error)
	GetPeople(filters map[string]string, page, limit int) ([]model.Person, error)
	UpdatePerson(id uint, updates map[string]interface{}) (*model.Person, error)
	DeletePerson(id uint) error
}

type PeopleServiceImpl struct {
	repo     repository.PeopleRepository
	enricher enrichment.Enricher
}

func (p PeopleServiceImpl) AddPerson(input model.InputPerson) (*model.Person, error) {
	person := &model.Person{
		Name:       input.Name,
		Surname:    input.Surname,
		Patronymic: input.Patronymic,
	}

	if err := p.enricher.Enrich(person); err != nil {
		log.Printf("error enriching person: %v", err)
		return nil, err
	}

	if err := p.repo.Create(person); err != nil {
		return nil, err
	}
	return person, nil
}

func (p PeopleServiceImpl) GetPeople(filters map[string]string, page, limit int) ([]model.Person, error) {
	people, err := p.repo.FindAll(filters, page, limit)
	if err != nil {
		return nil, err
	}

	log.Printf("Fetched %d people", len(people))
	return people, nil
}

func (p PeopleServiceImpl) UpdatePerson(id uint, updates map[string]interface{}) (*model.Person, error) {
	person, err := p.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	for k, v := range updates {
		switch k {
		case "name":
			person.Name = v.(string)
		case "surname":
			person.Surname = v.(string)
		case "patronymic":
			person.Patronymic = v.(string)
		case "age":
			switch val := v.(type) {
			case string:
				ageInt, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("invalid age value: %v", err)
				}
				person.Age = ageInt
			case float64:
				person.Age = int(val)
			case int:
				person.Age = val
			default:
				return nil, fmt.Errorf("unexpected type for age: %T", val)
			}
		case "gender":
			person.Gender = v.(string)
		case "nationality":
			person.National = v.(string)
		}
	}

	if err := p.repo.Update(person); err != nil {
		return nil, err
	}

	log.Printf("Updated person: ID=%d", id)
	return person, nil
}

func (p PeopleServiceImpl) DeletePerson(id uint) error {
	if err := p.repo.Delete(id); err != nil {
		return err
	}

	log.Printf("Deleted person: ID=%d", id)
	return nil
}

func NewPeopleService(repo repository.PeopleRepository, enricher enrichment.Enricher) PeopleService {
	return &PeopleServiceImpl{
		repo:     repo,
		enricher: enricher,
	}
}

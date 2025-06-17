package repository

import (
	"effective/internal/config"
	"effective/internal/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PeopleRepository interface {
	Create(person *model.Person) error
	FindAll(filters map[string]string, page, limit int) ([]model.Person, error)
	FindByID(id uint) (*model.Person, error)
	Update(person *model.Person) error
	Delete(id uint) error
}

type PostgresPeopleRepository struct {
	db *gorm.DB
}

func (p PostgresPeopleRepository) Create(person *model.Person) error {
	return p.db.Create(person).Error
}

func (p PostgresPeopleRepository) FindAll(filters map[string]string, page, limit int) ([]model.Person, error) {
	var people []model.Person
	query := p.db.Model(&model.Person{})

	for k, v := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	offset := limit * (page - 1)
	result := query.Offset(offset).Limit(limit).Find(&people)

	return people, result.Error
}

func (p PostgresPeopleRepository) FindByID(id uint) (*model.Person, error) {
	var person model.Person
	result := p.db.First(&person, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &person, nil
}

func (p PostgresPeopleRepository) Update(person *model.Person) error {
	return p.db.Save(person).Error
}

func (p PostgresPeopleRepository) Delete(id uint) error {
	return p.db.Delete(&model.Person{}, id).Error
}

func NewPeopleRepository(dbConfig config.DBConfig) (PeopleRepository, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &PostgresPeopleRepository{db: db}, nil
}

package controller

import (
	"github.com/penguin-statistics/partial-matrix/controller/item"
	"github.com/penguin-statistics/partial-matrix/controller/matrix"
	"github.com/penguin-statistics/partial-matrix/controller/stage"
	"github.com/penguin-statistics/partial-matrix/controller/zone"
)

type Collection struct {
	Item   *item.Controller
	Matrix *matrix.Controller
	Stage  *stage.Controller
	Zone   *zone.Controller
}

func NewCollection() *Collection {
	return &Collection{
		Item:   item.New(),
		Matrix: matrix.New(),
		Stage:  stage.New(),
		Zone:   zone.New(),
	}
}

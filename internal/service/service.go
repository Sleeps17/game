package service

import (
	"game/pkg/life"
)

type LifeService struct {
	CurrentWorld *life.World
	NextWorld    *life.World
}

func New(height, width, fill int) (*LifeService, error) {

	currentWorld, err := life.NewWorld(height, width)
	if err != nil {
		return nil, err
	}

	currentWorld.Seed(fill)

	newWorld, err := life.NewWorld(height, width)
	if err != nil {
		return nil, err
	}

	life.NextState(currentWorld, newWorld)

	return &LifeService{
		CurrentWorld: currentWorld,
		NextWorld:    newWorld,
	}, nil
}

func (ls *LifeService) NextState() *life.World {
	life.NextState(ls.CurrentWorld, ls.NextWorld)

	ls.CurrentWorld = ls.NextWorld

	return ls.CurrentWorld
}

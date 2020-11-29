package main

import (
	"math/rand"
	"time"

	"github.com/oktalz/go-sse/server"
)

type Animal struct {
	Name  string
	Sound string
}

type Animals struct {
	Event   string
	animals []Animal
}

func (r *Animals) Init(server *server.Server) error {
	// init what you need
	// add one event loop for example
	go func() {
		for {
			server.Emit(r.Event, r.GetRandom())
			time.Sleep(4 * time.Second)
		}
	}()
	r.animals = append(r.animals, Animal{"cat", "mijau"})
	r.animals = append(r.animals, Animal{"dog", "wof"})
	r.animals = append(r.animals, Animal{"mouse", "ciu"})
	r.animals = append(r.animals, Animal{"lion", "roar"})
	r.animals = append(r.animals, Animal{"gopher", "high pitch"})
	return nil
}

// you can call it directly from web also via /api/bind/Animals/GetRandom
func (r *Animals) GetRandom() Animal {
	return r.animals[rand.Intn(len(r.animals))] //nolint gosec
}

// you can call it directly from web also via /api/bind/Animals/GetRandomSound
func (r *Animals) GetRandomSound() string {
	return r.animals[rand.Intn(len(r.animals))].Sound //nolint gosec
}

// you can call it directly from web also via /api/bind/Animals/GetAnimalSound?args={name}
func (r *Animals) GetAnimalSound(name string) string {
	for _, animal := range r.animals {
		if animal.Name == name {
			return animal.Sound
		}
	}
	return ""
}

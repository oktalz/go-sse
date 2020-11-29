type Animals struct {
	animals []Animal
}

func (r *Animals) GetAnimalSound(name string) string {
	for _, animal := range r.animals {
		if animal.Name == name {
			return animal.Sound
		}
	}
	return ""
}

func main() {
	animals := &Animals{} // HL
	animals.Add(Animal{"cat", "mijau"})
	animals.Add(Animal{"dog", "wof"})
	animals.Add(Animal{"gopher", "high pitch"})

	// bind/Animals/GetAnimalSound?name={name} // HL
	err := server.Bind(animals) // HL
}

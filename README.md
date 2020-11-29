# go-sse

for examples see the [examples](examples)

presentation available at [talks.godoc.org](https://talks.godoc.org/github.com/oktalz/go-sse/presentation/present.slide#1)

## bind

``` Go
type RandomStruct struct {
    // ....
}

func (r *RandomStruct) GetRandom() int {
    return rand.Intn(100) // nolint:gosec
}

func main() {

    // ...

    random := &RandomStruct{}

    // bind/RandomStruct/GetRandom
    err := server.Bind(random)

    // ...

}
```

``` Go
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
    animals := &Animals{}
    animals.Add(Animal{"cat", "mijau"})
    animals.Add(Animal{"dog", "wof"})
    animals.Add(Animal{"gopher", "high pitch"})

    // bind/Animals/GetAnimalSound?name={name}
    err := server.Bind(animals)
}
```

## SSE

```
func main() {

    server.AddEvent("time", 5) // history of 5 messages available
    server.AddEvent("status", 1)

    go func() {
        for {
            server.Emit("time", time.Now())
            time.Sleep(42 * time.Second)
        }
    }()

    server.Emit("status", Status{
        Accepted: 98,
        Rejected: 2
    })

    // ...
}
```

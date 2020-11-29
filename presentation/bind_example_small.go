type RandomStruct struct {
	// ....
}

func (r *RandomStruct) GetRandom() int {
	return rand.Intn(100) // nolint:gosec
}

func main() {

	// ...

	random := &RandomStruct{} // HL

	// bind/RandomStruct/GetRandom
	err := server.Bind(random) // HL

	// ...

}

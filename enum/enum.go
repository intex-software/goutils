package enum

type Enum[T any] interface {
	String() string
	Value() T
}

/**/

type ExampleRegion interface {
	Enum[int]
	hidden()
}

const (
	World exampleRegion = iota
	Continent
	Country
	City
	Street
	House
)

type exampleRegion int

func (pf exampleRegion) hidden() {}

func (pf exampleRegion) String() string {
	switch pf {
	case World:
		return "World"
	case Continent:
		return "Continent"
	case Country:
		return "Country"
	case City:
		return "City"
	case Street:
		return "Street"
	case House:
		return "House"
	}

	return "Unknown"
}

func (pf exampleRegion) Value() int {
	return int(pf)
}

func NewExampleRegion(key string) ExampleRegion {
	switch key {
	case "World":
		return World
	case "Continent":
		return Continent
	case "Country":
		return Country
	case "City":
		return City
	case "Street":
		return Street
	case "House":
		return House
	}

	return nil
}

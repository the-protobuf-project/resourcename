package main

import (
	"fmt"
	"log"

	"github.com/oh-tarnished/runtime-go/ulid"
	"github.com/the-protobuf-project/resourcename"
)

// User represents a basic resource
type User struct {
	_    struct{} `resource:"//example.com/users/{id}/{name}"`
	ID   string   `resource:"id"`
	Name string   `resource:"name"`
}

// Address represents a nested struct
type Address struct {
	City string `resource:"city"`
	Zip  string `resource:"zip"`
}

// UserWithAddress demonstrates nested struct support
type UserWithAddress struct {
	_       struct{} `resource:"//example.com/users/{id}/{address.city}/{address.zip}"`
	ID      string   `resource:"id"`
	Address Address  `resource:"address"`
}

func main() {
	fmt.Println("=== Resource Name Marshaling Demo ===")

	fmt.Println("\n1. Basic Example:")
	demoBasic()

	fmt.Println("\n2. Nested Struct Example:")
	demoNested()
}

func demoBasic() {
	u := &User{ID: ulid.GenerateString(), Name: "Ria"}

	rn, err := resourcename.MarshalResource(u)
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}
	fmt.Printf("   Marshaled: %s\n", rn)

	u2 := &User{}
	err = resourcename.UnmarshalResource(rn, u2)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	fmt.Printf("   Unmarshaled: ID=%s, Name=%s\n", u2.ID, u2.Name)
}

func demoNested() {
	u := &UserWithAddress{
		ID:      ulid.GenerateString(),
		Address: Address{City: "NYC", Zip: "10001"},
	}

	rn, err := resourcename.MarshalResource(u)
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}
	fmt.Printf("   Marshaled: %s\n", rn)

	u2 := &UserWithAddress{}
	err = resourcename.UnmarshalResource(rn, u2)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	fmt.Printf("   Unmarshaled: ID=%s, City=%s, Zip=%s\n", u2.ID, u2.Address.City, u2.Address.Zip)
}

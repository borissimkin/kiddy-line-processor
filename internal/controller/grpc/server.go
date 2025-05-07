package grpc

import (
	"fmt"
	pb "kiddy-line-processor/internal/proto"
)

func awdawda() {
	p := pb.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-4321", Type: pb.PhoneType_PHONE_TYPE_HOME},
		},
	}

	fmt.Println(p.Name)
}

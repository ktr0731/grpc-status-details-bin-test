package main

import (
	"encoding/base64"
	"log"

	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/anypb"
)

// Generated by:
//
// m := &anypb.Any{}
// if err := anypb.MarshalFrom(m, &apipb.Detail{Body: "foo"}, proto.MarshalOptions{}); err != nil {
// 	log.Fatal(err)
// }
//
// log.Println(m.TypeUrl)
// log.Println(base64.StdEncoding.EncodeToString(m.Value))
var (
	typeURL = "type.googleapis.com/apipb.Detail"
	value   = "CgNmb28="
)

func main() {
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		log.Fatal(err)
	}

	src := &anypb.Any{
		TypeUrl: typeURL,
		Value:   b,
	}

	do(src)
	doWithResolver(src)
}

func do(src *anypb.Any) {
	_, err := anypb.UnmarshalNew(src, proto.UnmarshalOptions{})
	log.Println(err) // proto: not found
}

func doWithResolver(src *anypb.Any) {
	// Get the descriptor dynamically to avoid protoregistry.GlobalTypes registration.
	fds, err := protoparse.Parser{ImportPaths: []string{"."}}.ParseFiles("apipb/message.proto")
	if err != nil {
		log.Fatal(err)
	}

	resolver := new(protoregistry.Types)
	for _, fd := range fds {
		fd, err := protodesc.NewFile(fd.AsFileDescriptorProto(), nil)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < fd.Messages().Len(); i++ {
			mt := dynamicpb.NewMessageType(fd.Messages().Get(i))
			if err := resolver.RegisterMessage(mt); err != nil {
				log.Fatal(err)
			}
		}
	}

	dst, err := anypb.UnmarshalNew(src, proto.UnmarshalOptions{
		Resolver: resolver,
	})
	log.Println(protojson.Format(dst))
}
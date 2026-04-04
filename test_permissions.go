package main

import (
"fmt"
"reflect"
"im.turms/server/pkg/protocol"
)

func main() {
	var req protocol.TurmsRequest
	req.Kind = &protocol.TurmsRequest_CreateSessionRequest{}
	
	m := map[any]bool{
		&protocol.TurmsRequest_CreateSessionRequest{}: true,
	}
	
	fmt.Printf("Direct pointer map lookup: %v\n", m[req.GetKind()])
	
	m2 := map[reflect.Type]bool{
		reflect.TypeOf(&protocol.TurmsRequest_CreateSessionRequest{}): true,
	}
	fmt.Printf("Reflect type map lookup: %v\n", m2[reflect.TypeOf(req.GetKind())])
}

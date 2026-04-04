package dispatcher

// ServiceRequestMapping maps to ServiceRequestMapping.java
// In Go, since annotations are not supported, we don't have a direct equivalent
// to `@ServiceRequestMapping(TurmsRequest.KindCase value())`.
// This is typically handled by manually registering handlers into a routing map
// or using code generation tools based on comment directives like:
// // @ServiceRequestMapping(TurmsRequest_XX)
type ServiceRequestMapping struct{}

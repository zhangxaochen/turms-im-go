package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func init() {
	// Equivalent to Java's static initializer block in BaseTurmsApplication
	time.Local = time.UTC
}

func validateEnv() {
	// Equivalent to Java's JVM compatibility check (IncompatibleJvmException)
	ver := runtime.Version()
	if ver == "" || strings.HasPrefix(ver, "devel") {
		log.Fatalf("Incompatible Go runtime: %s", ver)
	}
}

// @Application(nodeType = NodeType.GATEWAY)
// @SpringBootApplication(scanBasePackages = {PackageConst.GATEWAY, PackageConst.SERVER_COMMON})
// @MappedFrom main(String[] args)
func main() {
	// Catch panics and ensure clean exit (equivalent to Java's try-catch in bootstrap)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Fatal error during startup: %v\n", r)
			os.Exit(1)
		}
	}()

	validateEnv()

	log.Println("Starting Turms Gateway...")
	// TODO: Replace with DI framework or manual wiring for Gateway components.

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Turms Gateway...")
}

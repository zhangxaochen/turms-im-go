package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func init() {
	// Equivalent to Java's static initializer block in BaseTurmsApplication
	time.Local = time.UTC
}

func validateEnv() {
	// Equivalent to Java's JVM compatibility check (IncompatibleJvmException)
	if runtime.Version() == "" {
		log.Fatalf("Incompatible Go runtime")
	}
}

// @Application(nodeType = NodeType.GATEWAY)
// @SpringBootApplication(scanBasePackages = {PackageConst.GATEWAY, PackageConst.SERVER_COMMON})
// @MappedFrom main(String[] args)
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Fatal error during startup: %v\n", r)
			// Ensure logs are flushed before exiting
			// (Mocked: in Java it waits for LoggerFactory to close)
			time.Sleep(1 * time.Second)
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

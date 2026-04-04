// Package main is the entry point for the Turms Service node.
//
// @MappedFrom TurmsServiceApplication.java
// @Application(nodeType = NodeType.SERVICE)
// @SpringBootApplication(scanBasePackages = {PackageConst.SERVICE, PackageConst.SERVER_COMMON})
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
	// Setting timezone to UTC matches Java's default behavior in containerized envs.
	time.Local = time.UTC
}

// validateEnv checks runtime compatibility, mirroring Java's IncompatibleJvmException check.
func validateEnv() {
	if runtime.Version() == "" {
		log.Fatalf("Incompatible Go runtime")
	}
}

// main is the entry point for turms-service.
// @MappedFrom TurmsServiceApplication.main(String[] args)
//
// Java equivalent:
//
//	public static void main(String[] args) {
//	    bootstrap(TurmsServiceApplication.class, args);
//	}
//
// "bootstrap" in Java triggers full Spring context: DI, component scanning,
// configuration loading, cluster node startup (NodeType.SERVICE).
//
// In Go, we perform manual dependency wiring instead of Spring DI.
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Fatal error during startup: %v\n", r)
			time.Sleep(1 * time.Second)
			os.Exit(1)
		}
	}()

	// Parse command-line arguments (equivalent to Spring's args → property source overrides).
	// TODO: integrate flag/pflag or viper for config override via --key=value args.
	validateEnv()

	log.Println("Starting Turms Service (NodeType=SERVICE)...")

	// TODO: Wire services manually (admin, user, message, group, etc.)
	// equivalent to @SpringBootApplication component scan of PackageConst.SERVICE + PackageConst.SERVER_COMMON

	// Graceful shutdown — wait for SIGINT or SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Turms Service...")
	// TODO: call cluster node stop, flush pending requests, close DB connections.
}

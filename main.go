package main

import (
	"flag"
	"fmt"
	"os"

	"YAMLtecture/internal/configuration"
	"YAMLtecture/internal/mermaid"
	"YAMLtecture/internal/validate"
)

var (
	validateFlag = flag.Bool("validate", false, "Validate the YAML architecture files")
	graphFlag    = flag.Bool("graph", false, "Generate a Mermaid diagram from the YAML architecture files")
	debugFlag    = flag.Bool("debug", false, "Enable debug output")
	dirFlag      = flag.String("dir", ".", "Directory containing YAML configuration files")
)

func main() {
	flag.Parse()

	// Ensure at least one command is provided
	if !*validateFlag && !*graphFlag {
		fmt.Println("Error: You must specify at least one command: --validate or --graph")
		flag.Usage()
		os.Exit(1)
	}

	// Load and parse YAML files
	config, err := configuration.LoadYAMLFiles(*dirFlag)
	if err != nil {
		fmt.Printf("Error loading YAML files: %v\n", err)
		os.Exit(1)
	}

	// Perform validation
	err = validate.ValidateConfig(config)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		os.Exit(1)
	}

	if *debugFlag {
		fmt.Println("Validation successful.")
	}

	// Execute commands
	if *validateFlag && !*graphFlag {
		// Only validate
		fmt.Println("Validation passed.")
	}

	if *graphFlag {
		// Generate Mermaid diagram
		mermaid, err := mermaid.GenerateMermaid(config)
		if err != nil {
			fmt.Printf("Error generating Mermaid diagram: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(mermaid)
	}
}

package main

import (
	"flag"
	"fmt"
	"os"

	"YAMLtecture/internal/configuration"
	"YAMLtecture/internal/mermaid"
	"YAMLtecture/internal/query"
)

var (
	queryFlag    = flag.String("query", "", "Query file to load")
	validateFlag = flag.Bool("validate", false, "Validate the YAML architecture files")
	graphFlag    = flag.Bool("graph", false, "Generate a Mermaid diagram from the YAML architecture files")
	debugFlag    = flag.Bool("debug", false, "Enable debug output")
	fileFlag     = flag.String("file", ".", "File path to a configuration file")
)

func main() {
	flag.Parse()

	// Ensure at least one command is provided
	if !*validateFlag && !*graphFlag {
		fmt.Println("Error: You must specify at least one command: --validate or --graph")
		flag.Usage()
		os.Exit(1)
	}

	// Load and parse YAML file
	config, err := configuration.LoadConfig(*fileFlag)
	if err != nil {
		fmt.Printf("Error loading YAML file: %v\n", err)
		os.Exit(1)
	}

	// Perform validation
	err = config.Validate()
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

	// If the queryFlag is set load in the query file
	if *queryFlag != "" {
		configQuery, err := query.LoadQuery(*queryFlag)
		if err != nil {
			fmt.Printf("Error loading query: %v\n", err)
			os.Exit(1)
		}

		// Now overwrite the config by applying the query
		// Apply the query
		configNew, err := query.ApplyQuery(configQuery, config)
		if err != nil {
			fmt.Printf("Error applying query: %v\n", err)
			os.Exit(1)
		}

		config = &configNew

		if *debugFlag {
			fmt.Println("Query applied successfully.")
		}
	}

	if *graphFlag {
		// Generate Mermaid diagram
		mermaid, err := mermaid.GenerateMermaid(config)
		if err != nil {
			fmt.Printf("Error generating Mermaid diagram: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(mermaid)
	} else {
		// Print the final configuration
		fmt.Print(config.YamlString())
	}
}

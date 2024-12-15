package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"

	c "YAMLtecture/internal/configuration"
	m "YAMLtecture/internal/mermaid"
	q "YAMLtecture/internal/query"
)

var (
	// In and Out
	inFlag  = flag.String("in", "", "Input file to load")
	outFlag = flag.String("out", "", "Output file to write")

	// Explicitely set the query and config files
	configFlag  = flag.String("configIn", "", "Input file for the Config YAML architecture file")
	queryFlag   = flag.String("queryIn", "", "Input file for the Query YAML architecture file")
	mermaidFlag = flag.String("mermaidIn", "", "Input file for the Mermaid settings")

	// The various commands to run
	validateConfigFlag  = flag.Bool("validateConfig", false, "Validate the Config YAML architecture file")
	validateQueryFlag   = flag.Bool("validateQuery", false, "Validate the Query YAML architecture file")
	validateMermaidFlag = flag.Bool("validateMermaid", false, "Validate the Mermaid settings")
	executeQueryFlag    = flag.Bool("executeQuery", false, "Execute the Query YAML architecture file")
	generateMermaidFlag = flag.Bool("generateMermaid", false, "Generate a Mermaid diagram from the Config YAML architecture file")

	// Modifiers
	debugFlag = flag.Bool("debug", false, "Enable debug output")
)

func main() {
	flag.Parse()

	// First determine what we are doing
	checkMultipleCommands(*validateConfigFlag, *validateQueryFlag, *validateMermaidFlag, *executeQueryFlag, *generateMermaidFlag)

	if *validateConfigFlag {
		// Validate the config file
		content := readFileContent(*configFlag, true, *inFlag, true, "")

		config, err := c.ParseYAML(content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
			os.Exit(1)
		}

		err = config.Validate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating configuration: %v\n", err)
			os.Exit(1)
		}

	} else if *validateQueryFlag {
		// Validate the query file
		content := readFileContent(*queryFlag, true, *inFlag, true, "")

		query, err := q.LoadQuery(content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading query: %v\n", err)
			os.Exit(1)
		}

		err = query.Validate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating query: %v\n", err)
			os.Exit(1)
		}

	} else if *validateMermaidFlag {
		// Validate the mermaid file
		content := readFileContent(*mermaidFlag, true, *inFlag, true, "")

		mermaid, err := m.ParseYAML(content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
			os.Exit(1)
		}

		err = mermaid.Validate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating mermaid: %v\n", err)
			os.Exit(1)
		}

	} else if *executeQueryFlag {
		// Execute the query file
		configContent := readFileContent(*configFlag, false, *inFlag, true, "")
		queryContent := readFileContent(*queryFlag, false, *inFlag, false, "")

		config, err := c.ParseYAML(configContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
			os.Exit(1)
		}

		err = config.Validate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating configuration: %v\n", err)
			os.Exit(1)
		}

		query, err := q.LoadQuery(queryContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading query: %v\n", err)
			os.Exit(1)
		}

		err = query.Validate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating query: %v\n", err)
			os.Exit(1)
		}

		result, err := q.ExecuteQuery(query, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing query: %v\n", err)
			os.Exit(1)
		}

		writeOutput(result.YamlString(), *outFlag)

	} else if *generateMermaidFlag {
		// Generate the Mermaid diagram
		queryContent := readFileContent(*configFlag, false, *inFlag, true, "")
		mermaidContent := readFileContent(*mermaidFlag, false, *inFlag, false, "\n")

		config, err := c.ParseYAML(queryContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
			os.Exit(1)
		}

		err = config.Validate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating configuration: %v\n", err)
			os.Exit(1)
		}

		mermaid, err := m.ParseYAML(mermaidContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
			os.Exit(1)
		}

		err = mermaid.Validate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating mermaid: %v\n", err)
			os.Exit(1)
		}

		mermaidDiagram, err := m.GenerateMermaid(config, mermaid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating Mermaid diagram: %v\n", err)
			os.Exit(1)
		}

		writeOutput(mermaidDiagram, *outFlag)

	} else {
		// Write error to error output
		fmt.Fprintf(os.Stderr, "Error: No command specified\n")
		os.Exit(1)
	}
}

// writeOutput will take the string content and if the -out flag is set, write it to the file, if not it writes to STDOUT
func writeOutput(content string, outFlag string) {
	if outFlag != "" {
		err := os.WriteFile(outFlag, []byte(content), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file %s: %v\n", outFlag, err)
			os.Exit(1)
		}
	} else {
		fmt.Println(content)
	}
}

// Read the content of a file based on the flags provided
func readFileContent(specificFlag string, allowGenericFlag bool, genericFlag string, allowStdin bool, defaultValue string) string {
	if specificFlag != "" {
		// Read from the specific flag first
		content, err := os.ReadFile(specificFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", specificFlag, err)
			os.Exit(1)
		}
		return string(content)
	} else if allowGenericFlag && genericFlag != "" {
		content, err := os.ReadFile(genericFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", genericFlag, err)
			os.Exit(1)
		}
		return string(content)
	} else if allowStdin {
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			content, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading from STDIN: %v\n", err)
				os.Exit(1)
			}
			return string(content)
		} else {
			fmt.Fprintf(os.Stderr, "Error: No input provided via STDIN\n")
			os.Exit(1)
			return ""
		}
	} else if defaultValue != "" {
		return defaultValue
	} else {
		fmt.Fprintf(os.Stderr, "Error: No input file specified and reading from STDIN is not allowed\n")
		os.Exit(1)
		return ""
	}
}

func checkMultipleCommands(commands ...bool) {
	count := 0
	for _, cmd := range commands {
		if cmd {
			count++
		}
	}
	if count > 1 {
		fmt.Fprintf(os.Stderr, "Error: Multiple commands specified\n")
		os.Exit(1)
	}
}

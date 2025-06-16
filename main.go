package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"

	"github.com/UnitVectorY-Labs/YAMLtecture/internal/common"
	c "github.com/UnitVectorY-Labs/YAMLtecture/internal/configuration"
	m "github.com/UnitVectorY-Labs/YAMLtecture/internal/mermaid"
	q "github.com/UnitVectorY-Labs/YAMLtecture/internal/query"
)

var (
	// In and Out
	inFlag  = flag.String("in", "", "Input file to load")
	outFlag = flag.String("out", "", "Output file to write")

	// Explicitly set the query and config files
	configFlag  = flag.String("configIn", "", "Input file for the Config YAML architecture file")
	queryFlag   = flag.String("queryIn", "", "Input file for the Query YAML architecture file")
	mermaidFlag = flag.String("mermaidIn", "", "Input file for the Mermaid settings")

	// The various commands to run
	validateConfigFlag  = flag.Bool("validateConfig", false, "Validate the Config YAML architecture file")
	validateQueryFlag   = flag.Bool("validateQuery", false, "Validate the Query YAML architecture file")
	validateMermaidFlag = flag.Bool("validateMermaid", false, "Validate the Mermaid settings")
	mergeConfigFlag     = flag.Bool("mergeConfig", false, "Merge the Config YAML architecture file")
	executeQueryFlag    = flag.Bool("executeQuery", false, "Execute the Query YAML architecture file")
	generateMermaidFlag = flag.Bool("generateMermaid", false, "Generate a Mermaid diagram from the Config YAML architecture file")

	// Modifiers
	debugFlag = flag.Bool("debug", false, "Enable debug output")
)

var Version = "dev" // This will be set by the build systems to the release version

func main() {

	showVersion := flag.Bool("version", false, "Print version")

	flag.Parse()

	if *showVersion {
		fmt.Println("Version:", Version)
		return
	}

	// First determine what we are doing
	checkMultipleCommands(*validateConfigFlag, *validateQueryFlag, *validateMermaidFlag, *mergeConfigFlag, *executeQueryFlag, *generateMermaidFlag)

	if *validateConfigFlag {
		// Validate the config file
		content := readFileContent(*configFlag, true, *inFlag, true, "")

		config, err := c.ParseYAML(content)
		if err != nil {
			common.PrintError("Error parsing YAML", err)
		}

		err = config.Validate()
		if err != nil {
			common.PrintError("Error validating configuration", err)
		}

	} else if *validateQueryFlag {
		// Validate the query file
		content := readFileContent(*queryFlag, true, *inFlag, true, "")

		query, err := q.ParseQuery(content)
		if err != nil {
			common.PrintError("Error loading query", err)
		}

		err = query.Validate()
		if err != nil {
			common.PrintError("Error validating query", err)
		}

	} else if *mergeConfigFlag {
		// Validate that inFlag was set and is a folder
		if *inFlag == "" {
			common.PrintError("No input folder specified", nil)
		}

		config, err := c.LoadFolder(*inFlag)
		if err != nil {
			common.PrintError("Error loading folder", err)
		}

		err = config.Validate()
		if err != nil {
			common.PrintError("Error validating configuration", err)
		}

		writeOutput(config.YamlString(), *outFlag)

	} else if *validateMermaidFlag {
		// Validate the mermaid file
		content := readFileContent(*mermaidFlag, true, *inFlag, true, "")

		mermaid, err := m.ParseYAML(content)
		if err != nil {
			common.PrintError("Error parsing YAML", err)
		}

		err = mermaid.Validate()
		if err != nil {
			common.PrintError("Error validating mermaid", err)
		}

	} else if *executeQueryFlag {
		// Execute the query file
		configContent := readFileContent(*configFlag, false, *inFlag, true, "")
		queryContent := readFileContent(*queryFlag, false, *inFlag, false, "")

		config, err := c.ParseYAML(configContent)
		if err != nil {
			common.PrintError("Error parsing YAML", err)
		}

		err = config.Validate()
		if err != nil {
			common.PrintError("Error validating configuration", err)
		}

		query, err := q.ParseQuery(queryContent)
		if err != nil {
			common.PrintError("Error loading query", err)
		}

		err = query.Validate()
		if err != nil {
			common.PrintError("Error validating query", err)
		}

		result, err := q.ExecuteQuery(query, config)
		if err != nil {
			common.PrintError("Error executing query", err)
		}

		writeOutput(result.YamlString(), *outFlag)

	} else if *generateMermaidFlag {
		// Generate the Mermaid diagram
		queryContent := readFileContent(*configFlag, false, *inFlag, true, "")
		mermaidContent := readFileContent(*mermaidFlag, false, *inFlag, false, "\n")

		config, err := c.ParseYAML(queryContent)
		if err != nil {
			common.PrintError("Error parsing YAML", err)
		}

		err = config.Validate()
		if err != nil {
			common.PrintError("Error validating configuration", err)
		}

		mermaid, err := m.ParseYAML(mermaidContent)
		if err != nil {
			common.PrintError("Error parsing YAML", err)
		}

		err = mermaid.Validate()
		if err != nil {
			common.PrintError("Error validating mermaid", err)
		}

		mermaidDiagram, err := m.GenerateMermaid(config, mermaid)
		if err != nil {
			common.PrintError("Error generating Mermaid diagram", err)
		}

		writeOutput(mermaidDiagram, *outFlag)

	} else {
		// Write error to error output
		common.PrintError("No command specified", nil)
	}
}

// writeOutput will take the string content and if the -out flag is set, write it to the file, if not it writes to STDOUT
func writeOutput(content string, outFlag string) {
	if outFlag != "" {
		err := os.WriteFile(outFlag, []byte(content), 0644)
		if err != nil {
			common.PrintError(fmt.Sprintf("Error writing to file %s", outFlag), err)
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
			common.PrintError(fmt.Sprintf("Error reading file %s", specificFlag), err)
		}
		return string(content)
	} else if allowGenericFlag && genericFlag != "" {
		content, err := os.ReadFile(genericFlag)
		if err != nil {
			common.PrintError(fmt.Sprintf("Error reading file %s", genericFlag), err)
		}
		return string(content)
	} else if allowStdin {
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			content, err := io.ReadAll(os.Stdin)
			if err != nil {
				common.PrintError("Error reading from STDIN", err)
			}
			return string(content)
		} else {
			common.PrintError("No input provided via STDIN", nil)
			return ""
		}
	} else if defaultValue != "" {
		return defaultValue
	} else {
		common.PrintError("No input file specified and reading from STDIN is not allowed", nil)
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
		common.PrintError("Multiple commands specified", nil)
	}
}

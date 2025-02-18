---
layout: default
title: Commands
nav_order: 4
permalink: /commands
---

# Commands

There are a number of different commands that YAMLtecture supports for various operations. These commands are used to interact with the YAMLtecture configuration, query the configuration, and render the configuration as a mermaid diagram.

## Validate Config

The validate config command, `--validateConfig`, takes in a configuration file and runs the validation checks on the configuration.

The inputs are used in the following order of precedence:

1. The `--configIn=<filePath>` flag
2. The `--in=<filePath>` flag
3. The STDIN

This command will output the validation errors and warnings to the console through the standard output.

## Validate Query

The validate query command, ``--validateQuery`, takes in a query file and runs the validation checks on the query.

The inputs are used in the following order of precedence:

1. The `--queryIn=<filePath>` flag
2. The `--in=<filePath>` flag
3. The STDIN

This command will output the validation errors and warnings to the console through the standard output.

## Validate Mermaid

The validate mermaid command, `--validateMermaid`, takes in a mermaid file and runs the validation checks on the mermaid file.

The inputs are used in the following order of precedence:

1. The `--mermaidIn=<filePath>` flag
2. The `--in=<filePath>` flag
3. The STDIN

This command will output the validation errors and warnings to the console through the standard output.

## Merge Config

The merge config command, `--mergeConfig`, takes in a folder path and merges all of the configuration files into a single output configuration file.

The only input for this command is the `--in=<folderPath>` flag which is used to specify a folder path that contains the configuration files unlike other commands.

The output of this command will be the result of the query which is the config YAML that is output to STDOUT or if `--out=<filePath>` is specified then the output will be written to the specified file.

## Execute Query

The execute query command, `--executeQuery`, takes in both a configuration file and a query file. The validate config and validate query checks are always performed, but the details as for the failure of these checks are not displayed.

Since this takes in multiple inputs the configuration file can be specified in the following order of precedence:

1. The `--configIn=<filePath>` flag
2. The STDIN

The query file can be specified in the following order of precedence:

1. The `--queryIn=<filePath>` flag

Note that the `--in=<filePath>` flag is not applicable for this command given it must take in both a configuration and query file.

The output of this command will be the result of the query which is the config YAML that is output to STDOUT or if `--out=<filePath>` is specified then the output will be written to the specified file.

## Generate Mermaid

The generate mermaid command, `--generateMermaid`, takes in a configuration file and renders the configuration as a mermaid flowchart.

Since this takes in multiple inputs the configuration file can be specified in the following order of precedence:

1. The `--configIn=<filePath>` flag
2. The STDIN

The mermaid settings can be specified in the following order of precedence:

1. The `--mermaidIn=<settings>` flag
2. A default set of settings is used

The output of this command will be a Mermaid flowchart that is output to STDOUT or if `--out=<filePath>` is specified then the output will be written to the specified file.

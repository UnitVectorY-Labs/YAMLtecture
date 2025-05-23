# Example

Data driven test cases are used by YAMLtecture.  Thie folder contains the data for the test cases that are used to validate the applications behavior in different circumstances, but also serves as useful references for various situations.

## Generation

Running the `generate.sh` command will compile and run YAMLtecture to validate the configuration, validate queries, run the query and generate the output, validate the mermaid, and run the mermaid to generate the output. These outputs are committed to the Git repository, allowing test cases to be quickly created and validated. The primary validation is whether the generated output matches what’s committed.

## Folder Structure

Each folder contains the following files that are hand crafted for each example / test case:

- `config.yaml`: The configuration file that defines the architecture.
- `mermaid.yaml`: The mermaid configuration file.

The following files are generated by the `generate.sh` script by running YAMLtecture:

- `mermaid.mmd`: The mermaid file that is generated by YAMLtecture.

Multiple queries can be defined for each config. These are stored in the `queries` folder. Each query is defined in its own folder with the name. Inside of that folder the following files are defined:

- `query.yaml`: The query file that defines the query.
- `mermaid.yaml`: The mermaid configuration file.

The following files are generated by the `generate.sh` script by running YAMLtecture:

- `config.yaml`: The configuration file that is generated by YAMLtecture.
- `mermaid.mmd`: The mermaid file that is generated by YAMLtecture.

### Invalid Files

A special folder is used for invalid files, `invalid` which contains a folder for each type of file that can be validated.

The `config` folder contains configuration files that are validated with the `--validateConfig` flag.

The `query` folder contains query files that are validated with the `--validateQuery` flag.

The `mermaid` folder contains mermaid files that are validated with the `--validateMermaid` flag.

Each of these folders contains a folder named for the test case. Inside of the folder there are two files.

The `input.yaml` file contains the actual input file that is used in the test case. This file is crafted to be invalid.

The `expected_error.txt` file contains the error message produced by the input file's validation.

## Testing Philosophy

The testing philosophy here is not to write individual code tests for each function of YAMLtecture. Instead a data driven approach is used by generating and validating the examples. This speeds up development while also providing a library of examples that demonstrates how YAMLtecture works.

As the application is enhanced running the `generate.sh` script will easily show with Git if one of the generated outputs changed or remains the same, thereby validating the changes if they are expected or not during the code review.

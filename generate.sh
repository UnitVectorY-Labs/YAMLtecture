#!/bin/bash

# Define color codes
RED="\033[0;31m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
YELLOW="\033[1;33m"
NC="\033[0m" # No Color

# Define log file
LOG_FILE="generate.log"

# Add failure flag
FAILURE=0

# Array to keep track of tree indicators for each depth
declare -a TREE_INDICATORS=()

# Function to log messages with timestamps
log() {
    local level=$1
    local message=$2
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $message" >> "$LOG_FILE"
}

# Logging function
# Arguments:
#   $1 - Depth level (integer)
#   $2 - Name of the file/folder
#   $3 - Status ("Valid", "Generated", "Error", "Processing", "Built", "Deleted")
#   $4 - Is last item in the current directory ("yes" or "no")
log_entry() {
  local depth=$1
  local name=$2
  local status=$3
  local is_last=$4
  local prefix=""

  # Build the prefix based on TREE_INDICATORS
  for ((i=1; i<depth; i++)); do
    if [ "${TREE_INDICATORS[i]}" == "yes" ]; then
      prefix+="    "
    else
      prefix+="│   "
    fi
  done

  # Determine the branch symbol
  if [ "$is_last" == "yes" ]; then
    prefix+="└── "
    TREE_INDICATORS[depth]="yes"
  else
    prefix+="├── "
    TREE_INDICATORS[depth]="no"
  fi

  # Determine the color based on status
  local color="$NC"
  case "$status" in
    Valid)
      color="$GREEN"
      ;;
    Generated)
      color="$BLUE"
      ;;
    Error)
      color="$RED"
      ;;
    Processing)
      color="$YELLOW"
      ;;
    Built)
      color="$GREEN"
      ;;
    Deleted)
      color="$RED"
      ;;
    Merged)
      color="$GREEN"
      ;;
    *)
      color="$NC"
      ;;
  esac

  # Print the log line
  echo -e "${prefix}${name} ${color}[${status}]${NC}"
}

# Helper function to execute commands and handle errors
# Arguments:
#   $1 - Command to execute
#   $2 - Depth level
#   $3 - Name to display in log
#   $4 - Success status message
#   $5 - Is last item ("yes" or "no")
execute_command() {
    local cmd=$1
    local depth=$2
    local name=$3
    local success_status=$4
    local is_last=$5

    # Save current stdout/stderr
    exec 3>&1 4>&2
    
    # Redirect command output to log file
    exec 1>>"$LOG_FILE" 2>&1
    
    echo "$cmd"
    set -o pipefail
    eval "$cmd"
    local status=$?
    echo ""
    
    # Restore stdout/stderr
    exec 1>&3 3>&- 2>&4 4>&-
    
    if [ $status -ne 0 ]; then
        log_entry "$depth" "$name" "Error" "$is_last"
        FAILURE=1
        return 1
    fi
    
    log_entry "$depth" "$name" "$success_status" "$is_last"
    return 0
}

# Function to compile the Go application
compile_go_app() {
    if ! execute_command "go build -o YAMLtecture" 1 "YAMLtecture" "Built" "yes"; then
        exit 1
    fi
}

# Function to process mermaid.yaml and generate mermaid.mmd
# Arguments:
#   $1 - Directory path
#   $2 - Depth level
process_mermaid() {
    local dir="${1%/}"
    local depth=$2

    [ ! -f "$dir/mermaid.yaml" ] && return 0

    if ! execute_command "./YAMLtecture --validateMermaid --mermaidIn=$dir/mermaid.yaml" "$depth" "mermaid.yaml" "Valid" "no"; then
        return 1
    fi

    if ! execute_command "./YAMLtecture --generateMermaid --configIn=$dir/config.yaml --mermaidIn=$dir/mermaid.yaml --out=$dir/mermaid.mmd" "$depth" "mermaid.mmd" "Generated" "no"; then
        return 1
    fi

    return 0
}

# Function to process queries within a directory
# Arguments:
#   $1 - Configuration directory path
#   $2 - Depth level
process_queries() {
    local config_dir=$1
    local depth=$2

    local queries_dir="$config_dir/queries"
    if [ ! -d "$queries_dir" ]; then
        log_entry "$depth" "queries" "Error" "no"
        FAILURE=1
        return 1
    fi

    log_entry "$depth" "queries" "Processing" "yes"

    local queries=()
    for query in "$queries_dir"/*/; do
        [ -d "$query" ] || continue
        queries+=("$query")
    done

    local total=${#queries[@]}
    local count=0

    for query in "${queries[@]}"; do
        count=$((count + 1))
        local current_is_last="no"
        [ $count -eq $total ] && current_is_last="yes"

        query=${query%*/}
        local queryName=$(basename "$query")
        log_entry "$((depth + 1))" "$queryName" "Processing" "$current_is_last"

        local query_file="$query/query.yaml"
        if [ ! -f "$query_file" ]; then
            log_entry "$((depth + 2))" "query.yaml" "Error" "no"
            FAILURE=1
            continue
        fi

        if ! execute_command "./YAMLtecture --validateQuery --queryIn=$query_file" "$((depth + 2))" "query.yaml" "Valid" "no"; then
            FAILURE=1
            continue
        fi

        if ! execute_command "./YAMLtecture --executeQuery --configIn=$config_dir/config.yaml --queryIn=$query_file --out=$query/config.yaml" "$((depth + 2))" "config.yaml" "Generated" "no"; then
            FAILURE=1
            continue
        fi

        if [ ! -f "$query/config.yaml" ]; then
            log_entry "$((depth + 2))" "config.yaml" "Error" "no"
            FAILURE=1
            continue
        fi

        if ! execute_command "./YAMLtecture --validateConfig --configIn=$query/config.yaml" "$((depth + 2))" "config.yaml" "Valid" "no"; then
            FAILURE=1
            continue
        fi

        process_mermaid "$query" "$((depth + 2))"
    done

    return $FAILURE
}

# Function to process a single configuration directory
# Arguments:
#   $1 - Configuration directory path
#   $2 - Depth level
#   $3 - Is last item in the parent directory ("yes" or "no")
process_config_dir() {
    local dir="${1%/}"
    local depth=$2
    local is_last=$3

    local name=$(basename "$dir")
    log_entry "$depth" "$name" "Processing" "$is_last"

    # Process the "configs" folder if it exists
    if [ -d "$dir/configs" ]; then
        if ! execute_command "./YAMLtecture --in=$dir/configs --mergeConfig --out=$dir/config.yaml" "$((depth + 1))" "configs" "Merged" "no"; then
            return 1
        fi
    fi

    if [ -f "$dir/config.yaml" ]; then
        if ! execute_command "./YAMLtecture --validateConfig --configIn=$dir/config.yaml" "$((depth + 1))" "config.yaml" "Valid" "no"; then
            return 1
        fi
    fi

    # Process mermaid if it exists
    [ -f "$dir/mermaid.yaml" ] && process_mermaid "$dir" "$((depth + 1))"

    # Process queries if they exist
    [ -d "$dir/queries" ] && process_queries "$dir" "$((depth + 1))"
}

# Function to process the tests directory
process_tests_dir() {
  local example_dir="tests"
  if [ ! -d "$example_dir" ]; then
    echo -e "${RED}Directory 'tests' does not exist.${NC}"
    FAILURE=1
    return 1
  fi
  
  # Get a list of directories (excluding hidden ones)
  local dirs=($(find "$example_dir" -maxdepth 1 -mindepth 1 -type d -not -path '*/\.*'))
  local total=${#dirs[@]}
  
  if [ $total -eq 0 ]; then
    echo -e "${RED}No test directories found.${NC}"
    FAILURE=1
    return 1
  fi
  
  local count=0
  for dir in "${dirs[@]}"; do
    count=$((count + 1))
    local is_last="no"
    [ $count -eq $total ] && is_last="yes"
    process_config_dir "$dir" 2 "$is_last"
  done
}

# Main script execution starts here
rm -f "$LOG_FILE"
FAILURE=0

# Compile the Go application
compile_go_app

# Process the example directory
process_tests_dir

# Print final status
if [ $FAILURE -eq 1 ]; then
    echo -e "\n${RED}ERROR: One or more operations failed during execution. Check the log file for details.${NC}"
    exit 1
fi

echo -e "\n${GREEN}SUCCESS: All operations completed successfully.${NC}"

# Optional cleanup
# rm -f YAMLtecture
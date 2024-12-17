#!/bin/bash

# Define color codes
RED="\033[0;31m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
YELLOW="\033[1;33m"
NC="\033[0m" # No Color

# Array to keep track of tree indicators for each depth
declare -a TREE_INDICATORS=()

# Define log file
LOG_FILE="example.log"

# Function to log messages with timestamps
log() {
    local level=$1
    local message=$2
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$level] $message" | tee -a "$LOG_FILE"
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
    *)
      color="$NC"
      ;;
  esac

  # Print the log line
  echo -e "${prefix}${name} ${color}[${status}]${NC}"
}

# Function to compile the Go application
compile_go_app() {
  go build -o YAMLtecture >>"$LOG_FILE" 2>&1
  if [ $? -ne 0 ]; then
    log_entry 1 "YAMLtecture" "Error" "yes"
    exit 1
  else
    # Log the built application as [Built]
    log_entry 1 "YAMLtecture" "Built" "yes"
  fi
}

# Function to process mermaid.yaml and generate mermaid.mmd
# Arguments:
#   $1 - Directory path
#   $2 - Depth level
process_mermaid() {
  local dir="${1%/}"
  local depth=$2

  local mermaid_in="$dir/mermaid.yaml"
  local mermaid_out="$dir/mermaid.mmd"

  if [ -f "$mermaid_in" ]; then
    echo "./YAMLtecture --validateMermaid --mermaidIn=$mermaid_in" >> "$LOG_FILE"
    ./YAMLtecture --validateMermaid --mermaidIn="$mermaid_in" >>"$LOG_FILE" 2>&1
    echo "" >> "$LOG_FILE"
    if [ $? -ne 0 ]; then
      log_entry "$depth" "mermaid.yaml" "Error" "no"
    else
      log_entry "$depth" "mermaid.yaml" "Valid" "no"
      # Generate mermaid.mmd
      echo "./YAMLtecture --generateMermaid --configIn=$dir/config.yaml --mermaidIn=$mermaid_in --out=$mermaid_out" >> "$LOG_FILE"
      ./YAMLtecture --generateMermaid --configIn="$dir/config.yaml" --mermaidIn="$mermaid_in" --out="$mermaid_out" >>"$LOG_FILE" 2>&1
      echo "" >> "$LOG_FILE"
      if [ $? -ne 0 ]; then
        log_entry "$depth" "mermaid.mmd" "Error" "no"
      else
        log_entry "$depth" "mermaid.mmd" "Generated" "no"
      fi
    fi
  else
    # No mermaid.yaml found, remove mermaid.mmd if exists
    if [ -f "$mermaid_out" ]; then
      rm -f "$mermaid_out"
      log_entry "$depth" "mermaid.mmd" "Deleted" "no"
    fi
  fi
}

# Function to process queries within a directory
# Arguments:
#   $1 - Configuration directory path
#   $2 - Depth level
process_queries() {
  local config_dir=$1
  local depth=$2

  local queries_dir="$config_dir/queries"
  if [ -d "$queries_dir" ]; then
    # Log the 'queries' directory
    local is_last_queries="yes" # Assuming 'queries' is the last item in its parent directory
    log_entry "$depth" "queries" "Processing" "$is_last_queries"

    # Collect all query directories
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
      if [ "$count" -eq "$total" ]; then
        current_is_last="yes"
      fi

      # Remove the trailing slash and get the basename
      query=${query%*/}
      local queryName=$(basename "$query")

      # Log the query directory
      log_entry "$((depth + 1))" "$queryName" "Processing" "$current_is_last"

      # Validate the query.yaml
      local query_file="$query/query.yaml"
      if [ -f "$query_file" ]; then
        echo "./YAMLtecture --validateQuery --queryIn=$query_file" >> "$LOG_FILE"
        ./YAMLtecture --validateQuery --queryIn="$query_file" >>"$LOG_FILE" 2>&1
        echo "" >> "$LOG_FILE"
        if [ $? -ne 0 ]; then
          log_entry "$((depth + 2))" "query.yaml" "Error" "no"
        else
          log_entry "$((depth + 2))" "query.yaml" "Valid" "no"
          # Execute the query
          echo "./YAMLtecture --executeQuery --configIn=$config_dir/config.yaml --queryIn=$query_file --out=$query/config.yaml" >> "$LOG_FILE"
          ./YAMLtecture --executeQuery --configIn="$config_dir/config.yaml" --queryIn="$query_file" --out="$query/config.yaml" >>"$LOG_FILE" 2>&1
          echo "" >> "$LOG_FILE"
          if [ $? -ne 0 ]; then
            log_entry "$((depth + 2))" "config.yaml" "Error" "no"
          else
            log_entry "$((depth + 2))" "config.yaml" "Generated" "no"
            # Validate the generated config.yaml
            echo "./YAMLtecture --validateConfig --configIn=$query/config.yaml" >> "$LOG_FILE"
            ./YAMLtecture --validateConfig --configIn="$query/config.yaml" >>"$LOG_FILE" 2>&1
            echo "" >> "$LOG_FILE"
            if [ $? -ne 0 ]; then
              log_entry "$((depth + 2))" "config.yaml" "Error" "no"
            else
              # Process mermaid for the query
              process_mermaid "$query" "$((depth + 2))"
            fi
          fi
        fi
      else
        log_entry "$((depth + 2))" "query.yaml" "Error" "no"
      fi
    done
  else
    log_entry "$depth" "queries" "Error" "no"
  fi
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

  # Collect all items in the configuration directory except 'queries'
  local items=()
  for item in "$dir"/*; do
    if [ -e "$item" ] && [ "$(basename "$item")" != "queries" ]; then
      items+=("$item")
    fi
  done

  # Check if 'queries' exists to determine if it's the last item
  local has_queries=0
  if [ -d "$dir/queries" ]; then
    has_queries=1
  fi

  local total=${#items[@]}
  if [ "$has_queries" -eq 1 ]; then
    total=$((total + 1))
  fi

  local count=0

  # Process each item in the configuration directory
  for item in "${items[@]}"; do
    count=$((count + 1))
    local current_is_last="no"
    if [ "$count" -eq "$total" ]; then
      current_is_last="yes"
    fi

    local item_name
    item_name=$(basename "$item")

    case "$item_name" in
      config.yaml)
        # Validate the config.yaml
        echo "./YAMLtecture --validateConfig --configIn=$item" >> "$LOG_FILE"
        ./YAMLtecture --validateConfig --configIn="$item" >>"$LOG_FILE" 2>&1
        echo "" >> "$LOG_FILE"
        if [ $? -ne 0 ]; then
          log_entry "$((depth + 1))" "config.yaml" "Error" "no"
        else
          log_entry "$((depth + 1))" "config.yaml" "Valid" "no"
        fi
        ;;
      mermaid.yaml)
        # Validate and generate mermaid.mmd
        echo "./YAMLtecture --validateMermaid --mermaidIn=$item" >> "$LOG_FILE"
        ./YAMLtecture --validateMermaid --mermaidIn="$item" >>"$LOG_FILE" 2>&1
        echo "" >> "$LOG_FILE"
        if [ $? -ne 0 ]; then
          log_entry "$((depth + 1))" "mermaid.yaml" "Error" "no"
        else
          log_entry "$((depth + 1))" "mermaid.yaml" "Valid" "no"
          echo "./YAMLtecture --generateMermaid --configIn=$dir/config.yaml --mermaidIn=$item --out=$dir/mermaid.mmd" >> "$LOG_FILE"
          ./YAMLtecture --generateMermaid --configIn="$dir/config.yaml" --mermaidIn="$item" --out="$dir/mermaid.mmd" >>"$LOG_FILE" 2>&1
            echo "" >> "$LOG_FILE"
          if [ $? -ne 0 ]; then
            log_entry "$((depth + 1))" "mermaid.mmd" "Error" "no"
          else
            log_entry "$((depth + 1))" "mermaid.mmd" "Generated" "no"
          fi
        fi
        ;;
      *)
        # Handle other files if necessary
        ;;
    esac
  done

  # If 'queries' directory exists, process it
  if [ "$has_queries" -eq 1 ]; then
    local is_last_queries="yes" # Assuming 'queries' is the last item
    process_queries "$dir" "$((depth + 1))"
  fi
}

# Function to process the example directory
process_example_dir() {
  local example_dir="example"
  if [ -d "$example_dir" ]; then
    # Collect top-level directories
    local dirs=()
    for dir in "$example_dir"/*/; do
      [ -d "$dir" ] || continue
      dirs+=("$dir")
    done

    local total=${#dirs[@]}
    local count=0

    for dir in "${dirs[@]}"; do
      count=$((count + 1))
      local is_last="no"
      if [ "$count" -eq "$total" ]; then
        is_last="yes"
      fi
      process_config_dir "$dir" 2 "$is_last"
    done
  else
    echo -e "${RED}Directory 'example' does not exist.${NC}"
    exit 1
  fi
}

# Main script execution starts here
rm -f "$LOG_FILE"

# Compile the Go application
compile_go_app

# Process the example directory
process_example_dir

# Optionally, clean up the binary after processing
# rm -f YAMLtecture
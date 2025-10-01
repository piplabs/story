#!/bin/bash

# Script to parse test_N*.json files and generate CSV
# Usage: ./generate-test-csv.sh [directory]

set -e

# Default directory if not provided
DIR="${1:-1315}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1" >&2
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" >&2
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# Check if directory exists
if [ ! -d "$DIR" ]; then
    print_error "Directory $DIR does not exist"
    exit 1
fi

# Function to extract action name from filename
extract_action_name() {
    local filename="$1"
    
    # Extract the type (schedule/execute/cancel) and action part
    local type=$(echo "$filename" | sed -E 's/^test_[0-9]+\.-(.+)-(schedule|execute|cancel)\.json$/\2/')
    local action=$(echo "$filename" | sed -E 's/^test_[0-9]+\.-(.+)-(schedule|execute|cancel)\.json$/\1/')
    
    # Replace dashes with spaces in action
    action=$(echo "$action" | tr '-' ' ')
    
    # Capitalize each word in action
    action=$(echo "$action" | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2)); print}')
    
    # Capitalize the type
    type=$(echo "$type" | awk '{print toupper(substr($1,1,1)) tolower(substr($1,2))}')
    
    # Combine type and action
    echo "$type $action"
}

# Function to extract number from filename
extract_number() {
    local filename="$1"
    echo "$filename" | sed -E 's/^test_([0-9]+)\.-.*/\1/'
}

# Function to extract type from filename
extract_type() {
    local filename="$1"
    if [[ "$filename" =~ -schedule\.json$ ]]; then
        echo "schedule"
    elif [[ "$filename" =~ -execute\.json$ ]]; then
        echo "execute"
    elif [[ "$filename" =~ -cancel\.json$ ]]; then
        echo "cancel"
    else
        echo "unknown"
    fi
}

# Function to get sort priority for type
get_type_priority() {
    local type="$1"
    case "$type" in
        "schedule") echo "1" ;;
        "execute") echo "2" ;;
        "cancel") echo "3" ;;
        *) echo "4" ;;
    esac
}

# Function to parse JSON and extract fields
parse_json() {
    local file="$1"
    
    # Check if jq is available
    if command -v jq &> /dev/null; then
        # Use jq for robust JSON parsing
        local from=$(jq -r '.[0].from' "$file" 2>/dev/null || echo "")
        local to=$(jq -r '.[0].to' "$file" 2>/dev/null || echo "")
        local data=$(jq -r '.[0].data' "$file" 2>/dev/null || echo "")
    else
        # Fallback to basic parsing without jq
        local content=$(cat "$file")
        local from=$(echo "$content" | grep -o '"from":"[^"]*"' | head -1 | cut -d'"' -f4)
        local to=$(echo "$content" | grep -o '"to":"[^"]*"' | head -1 | cut -d'"' -f4)
        local data=$(echo "$content" | grep -o '"data":"[^"]*"' | head -1 | cut -d'"' -f4)
    fi
    
    echo "$from|$to|$data"
}

# Main function
main() {
    print_info "Generating CSV from test JSON files in directory: $DIR"
    
    # Find all test_*.json files
    local test_files=($(find "$DIR" -name "test_*.json" -type f | sort))
    
    if [ ${#test_files[@]} -eq 0 ]; then
        print_error "No test_*.json files found in directory $DIR"
        exit 1
    fi
    
    print_info "Found ${#test_files[@]} test files"
    
    # Create temporary file for sorting
    local temp_file=$(mktemp)
    
    # Process each file
    for file in "${test_files[@]}"; do
        local filename=$(basename "$file")
        local number=$(extract_number "$filename")
        local type=$(extract_type "$filename")
        local type_priority=$(get_type_priority "$type")
        local action=$(extract_action_name "$filename")
        local json_data=$(parse_json "$file")
        
        # Create sort key and data line
        local sort_key=$(printf "%03d-%d" "$number" "$type_priority")
        echo "$sort_key|$action|$json_data" >> "$temp_file"
    done
    
    # Output CSV header
    echo "Action,from,to,data,Technology,Signers,Status,Result"
    
    # Sort and output data
    sort "$temp_file" | while IFS='|' read -r sort_key action from to data; do
        # Escape any commas or quotes in the data for CSV
        action=$(echo "$action" | sed 's/"/\\"/g')
        from=$(echo "$from" | sed 's/"/\\"/g')
        to=$(echo "$to" | sed 's/"/\\"/g')
        data=$(echo "$data" | sed 's/"/\\"/g')
        
        # Output CSV row with constant fields
        echo "\"$action\",\"$from\",\"$to\",\"$data\",\"MPCVault\",\"Governance\",\"Not Started\",\"\""
    done
    
    # Clean up
    rm "$temp_file"
    
    print_success "CSV generation completed successfully"
}

# Show usage if help requested
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    echo "Usage: $0 [directory]"
    echo ""
    echo "Generate CSV from test_*.json files in the specified directory"
    echo ""
    echo "Parameters:"
    echo "  directory    Directory containing test_*.json files (default: 1315)"
    echo ""
    echo "Output:"
    echo "  CSV data to stdout with columns: Action,from,to,data,Technology,Signers,Status,Result"
    echo ""
    echo "Examples:"
    echo "  $0                    # Use default directory (1315)"
    echo "  $0 1315              # Use specific directory"
    echo "  $0 1315 > output.csv # Save to file"
    exit 0
fi

# Run main function
main "$@"
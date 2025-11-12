#!/bin/bash

# Script to parse numbered *.json files and generate CSV with verification rows
# Usage: ./generate-numbered-csv.sh [directory]

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
    
    # Extract the type (schedule/execute) and action part
    local type=$(echo "$filename" | sed -E 's/^[0-9]+\.[0-9]+-(.+)-(schedule|execute)\.json$/\2/')
    local action=$(echo "$filename" | sed -E 's/^[0-9]+\.[0-9]+-(.+)-(schedule|execute)\.json$/\1/')
    
    # Replace safe-migr- prefix and dashes with spaces
    action=$(echo "$action" | sed 's/^safe-migr-//')
    action=$(echo "$action" | tr '-' ' ')
    
    # Capitalize each word in action
    action=$(echo "$action" | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2)); print}')
    
    # Capitalize the type
    type=$(echo "$type" | awk '{print toupper(substr($1,1,1)) tolower(substr($1,2))}')
    
    # Combine type and action
    echo "$type $action"
}

# Function to extract number from filename for sorting
extract_sort_number() {
    local filename="$1"
    # Extract major.minor version and convert to sortable format
    local major=$(echo "$filename" | sed -E 's/^([0-9]+)\.([0-9]+)-.*/\1/')
    local minor=$(echo "$filename" | sed -E 's/^([0-9]+)\.([0-9]+)-.*/\2/')
    printf "%02d.%02d" "$major" "$minor"
}

# Function to extract type from filename
extract_type() {
    local filename="$1"
    if [[ "$filename" =~ -schedule\.json$ ]]; then
        echo "schedule"
    elif [[ "$filename" =~ -execute\.json$ ]]; then
        echo "execute"
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
        *) echo "3" ;;
    esac
}

# Function to parse JSON and extract fields
parse_json() {
    local file="$1"
    
    # Debug: Check if file exists and has content
    if [ ! -f "$file" ]; then
        print_error "File not found: $file" >&2
        echo "||"
        return
    fi
    
    # Check if jq is available
    if command -v jq &> /dev/null; then
        # Use jq for robust JSON parsing
        local from=$(jq -r '.[0].from // ""' "$file" 2>/dev/null)
        local to=$(jq -r '.[0].to // ""' "$file" 2>/dev/null)
        local data=$(jq -r '.[0].data // ""' "$file" 2>/dev/null)
    else
        # Fallback to basic parsing without jq
        local content=$(cat "$file")
        local from=$(echo "$content" | grep -o '"from":"[^"]*"' | head -1 | cut -d'"' -f4 || echo "")
        local to=$(echo "$content" | grep -o '"to":"[^"]*"' | head -1 | cut -d'"' -f4 || echo "")
        local data=$(echo "$content" | grep -o '"data":"[^"]*"' | head -1 | cut -d'"' -f4 || echo "")
    fi
    
    # Debug output
    if [ -z "$data" ]; then
        print_error "Warning: No data found in file $file" >&2
    fi
    
    echo "$from|$to|$data"
}

# Function to output verification rows
output_verification_rows() {
    local signers=("SC - Vera" "SC - Michael" "SC - Anton" "SC - Patrick")
    
    for signer in "${signers[@]}"; do
        echo "\"Verify\",\"\",\"\",\"\",\"Manual\",\"$signer\",\"Not Started\",\"\""
    done
}

# Main function
main() {
    print_info "Generating CSV from numbered JSON files in directory: $DIR"
    
    # Find all numbered *.json files (exclude cancel files)
    local numbered_files=($(find "$DIR" -name "[0-9]*.[0-9]*-*.json" -type f | grep -v cancel | sort))
    
    if [ ${#numbered_files[@]} -eq 0 ]; then
        print_error "No numbered JSON files found in directory $DIR"
        exit 1
    fi
    
    print_info "Found ${#numbered_files[@]} numbered files (excluding cancel files)"
    
    # Create temporary file for sorting
    local temp_file=$(mktemp)
    
    # Process each file
    for file in "${numbered_files[@]}"; do
        local filename=$(basename "$file")
        local sort_number=$(extract_sort_number "$filename")
        local type=$(extract_type "$filename")
        local type_priority=$(get_type_priority "$type")
        local action=$(extract_action_name "$filename")
        local json_data=$(parse_json "$file")
        
        # Create sort key and data line
        local sort_key=$(printf "%s-%d" "$sort_number" "$type_priority")
        echo "$sort_key|$sort_number|$action|$json_data" >> "$temp_file"
    done
    
    # Output CSV header
    echo "Action,from,to,data,Technology,Signers,Status,Result"
    
    # Sort and output data with verification rows inserted after each execute
    sort "$temp_file" | while IFS='|' read -r sort_key sort_number action from to data; do
        # Extract type from sort_key (1=schedule, 2=execute)
        local type_priority=$(echo "$sort_key" | cut -d'-' -f2)
        
        # Escape any commas or quotes in the data for CSV
        action=$(echo "$action" | sed 's/"/\\"/g')
        from=$(echo "$from" | sed 's/"/\\"/g')
        to=$(echo "$to" | sed 's/"/\\"/g')
        data=$(echo "$data" | sed 's/"/\\"/g')
        
        # Output CSV row with constant fields
        echo "\"$action\",\"$from\",\"$to\",\"$data\",\"MPCVault\",\"Governance\",\"Not Started\",\"\""
        
        # Insert verification rows after each schedule operation (between schedule and execute)
        if [[ "$type_priority" == "1" ]]; then
            output_verification_rows
        fi
    done
    
    # Clean up
    rm "$temp_file"
    
    print_success "CSV generation completed successfully"
}

# Show usage if help requested
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    echo "Usage: $0 [directory]"
    echo ""
    echo "Generate CSV from numbered *.json files in the specified directory"
    echo "Ignores cancel files and inserts verification rows in the middle"
    echo ""
    echo "Parameters:"
    echo "  directory    Directory containing numbered *.json files (default: 1315)"
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
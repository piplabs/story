#!/bin/bash

# Migration script runner for Safe Migration
# This script runs all migration scripts in the correct sequence
# Usage: ./run-migration.sh <RPC_URL> [--broadcast]

# set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to print usage
usage() {
    echo "Usage: $0 <RPC_URL> [--broadcast]"
    echo ""
    echo "Parameters:"
    echo "  <RPC_URL>     Required. The RPC endpoint URL"
    echo "  --broadcast   Optional. Include this flag to broadcast transactions"
    echo ""
    echo "Examples:"
    echo "  $0 http://localhost:8545"
    echo "  $0 http://localhost:8545 --broadcast"
    echo "  $0 https://rpc.ankr.com/eth --broadcast"
    exit 1
}

# Check if at least one argument is provided
if [ $# -lt 1 ]; then
    print_error "Missing required RPC_URL parameter"
    usage
fi

# Parse arguments
RPC_URL="$1"
BROADCAST_FLAG=""

# Check for --broadcast flag
if [ $# -eq 2 ]; then
    if [ "$2" = "--broadcast" ]; then
        BROADCAST_FLAG="--broadcast"
        print_warning "Broadcasting is ENABLED - transactions will be sent to the network!"
    else
        print_error "Invalid second parameter. Use --broadcast or omit it."
        usage
    fi
else
    print_info "Dry run mode - transactions will NOT be broadcast"
fi

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONTRACT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

print_info "Script directory: $SCRIPT_DIR"
print_info "Contract root: $CONTRACT_ROOT"
print_info "RPC URL: $RPC_URL"

# Change to contract root directory
cd "$CONTRACT_ROOT"

# Define the migration scripts in order (excluding timelock deployment)
MIGRATION_SCRIPTS=(
    "2.1.TransferOwnershipProxyAdmin1.s.sol"
    "2.2.TransferOwnershipProxyAdmin2.s.sol"
    "2.3.TransferOwnershipProxyAdmin3.s.sol"
    "2.4.TransferOwnershipProxyAdmin4.s.sol"
    "3.1.TransferOwnershipUpgradesEntrypoint.s.sol"
    "3.2.ReceiveOwnershipUpgradesEntryPoint.s.sol"
    "3.3.TransferOwnershipUBIPool.s.sol"
    "3.4.ReceiveOwnershipUBIPool.s.sol"
    "3.5.TransferOwnershipIPTokenStaking.s.sol"
    "3.6.ReceiveOwnershipIPTokenStaking.s.sol"
    "4.RenounceGovernanceRoles.s.sol"
)

# Function to get the correct contract name from filename
get_contract_name() {
    local script_file="$1"
    
    case "$script_file" in
        "2.1.TransferOwnershipProxyAdmin1.s.sol")
            echo "TransferOwnershipsProxyAdmin1"
            ;;
        "2.2.TransferOwnershipProxyAdmin2.s.sol")
            echo "TransferOwnershipsProxyAdmin2"
            ;;
        "2.3.TransferOwnershipProxyAdmin3.s.sol")
            echo "TransferOwnershipsProxyAdmin3"
            ;;
        "2.4.TransferOwnershipProxyAdmin4.s.sol")
            echo "TransferOwnershipsProxyAdmin4"
            ;;
        "3.1.TransferOwnershipUpgradesEntrypoint.s.sol")
            echo "TransferOwnershipsUpgradesEntrypoint"
            ;;
        "3.2.ReceiveOwnershipUpgradesEntryPoint.s.sol")
            echo "ReceiveOwnershipUpgradesEntryPoint"
            ;;
        "3.3.TransferOwnershipUBIPool.s.sol")
            echo "TransferOwnershipUBIPool"
            ;;
        "3.4.ReceiveOwnershipUBIPool.s.sol")
            echo "ReceiveOwnershipUBIPool"
            ;;
        "3.5.TransferOwnershipIPTokenStaking.s.sol")
            echo "TransferOwnershipIPTokenStaking"
            ;;
        "3.6.ReceiveOwnershipIPTokenStaking.s.sol")
            echo "ReceiveOwnershipIPTokenStaking"
            ;;
        "4.RenounceGovernanceRoles.s.sol")
            echo "RenounceGovernanceRoles"
            ;;
        *)
            # Fallback: remove .s.sol and number prefix
            local contract_name=$(basename "$script_file" .s.sol)
            contract_name=$(echo "$contract_name" | sed 's/^[0-9]*\.[0-9]*\.//')
            echo "$contract_name"
            ;;
    esac
}

# Function to run a single migration script
run_migration_script() {
    local script_file="$1"
    local script_path="script/admin-actions/migrate-to-safe/$script_file"
    
    # Get the correct contract name using our mapping function
    local contract_name=$(get_contract_name "$script_file")
    
    print_info "=========================================="
    print_info "Running migration script: $script_file"
    print_info "Contract: $contract_name"
    print_info "=========================================="
    
    # Check if script file exists
    if [ ! -f "$script_path" ]; then
        print_error "Script file not found: $script_path"
        exit 1
    fi
    
    # Build the forge command
    local forge_cmd="forge script $script_path:$contract_name --rpc-url $RPC_URL"
    
    if [ -n "$BROADCAST_FLAG" ]; then
        forge_cmd="$forge_cmd $BROADCAST_FLAG"
    fi
    
    # Add verbose flag
    forge_cmd="$forge_cmd -vvv"
    
    print_info "Executing: $forge_cmd"
    
    # Execute the command and capture exit code
    if eval "$forge_cmd"; then
        print_success "✅ $script_file completed successfully"
    else
        local exit_code=$?
        print_error "❌ $script_file failed with exit code $exit_code"
        exit $exit_code
    fi
    
    # Wait a moment between scripts
    sleep 2
}

# Function to confirm before proceeding (only in broadcast mode)
confirm_execution() {
    if [ -n "$BROADCAST_FLAG" ]; then
        echo ""
        print_warning "⚠️  WARNING: You are about to run the migration sequence with broadcasting enabled!"
        print_warning "This will execute ${#MIGRATION_SCRIPTS[@]} scripts and send transactions to the network."
        print_warning "Note: Timelock deployment (1.DeployNewTimelock.s.sol) is excluded and must be run separately."
        echo ""
        echo "Migration sequence:"
        for i in "${!MIGRATION_SCRIPTS[@]}"; do
            printf "%2d. %s\n" $((i+1)) "${MIGRATION_SCRIPTS[$i]}"
        done
        echo ""
        read -p "Are you sure you want to proceed? (type 'yes' to continue): " -r
        if [[ ! $REPLY =~ ^yes$ ]]; then
            print_info "Migration cancelled by user"
            exit 0
        fi
        echo ""
    fi
}

# Main execution
main() {
    print_info "🚀 Starting Safe Migration Script Runner"
    print_info "Total scripts to run: ${#MIGRATION_SCRIPTS[@]} (timelock deployment excluded)"
    
    # Confirm execution if broadcasting
    confirm_execution
    
    local start_time=$(date +%s)
    local success_count=0
    
    # Run each migration script in sequence
    for script in "${MIGRATION_SCRIPTS[@]}"; do
        run_migration_script "$script"
        ((success_count++))
    done
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    print_success "🎉 All migration scripts completed successfully!"
    print_success "Scripts executed: $success_count/${#MIGRATION_SCRIPTS[@]}"
    print_success "Total execution time: ${duration} seconds"
    
    if [ -n "$BROADCAST_FLAG" ]; then
        print_success "✅ Migration has been broadcast to the network"
        print_info "Please verify all transactions on the blockchain explorer"
    else
        print_info "ℹ️  This was a dry run - no transactions were broadcast"
        print_info "Add --broadcast flag to execute transactions on the network"
    fi
}

# Check for required environment variables
check_environment() {
    # Check if forge is available
    if ! command -v forge &> /dev/null; then
        print_error "forge command not found. Please install Foundry."
        exit 1
    fi
}

# Run environment checks
check_environment

# Execute main function
main

print_info "Migration script runner completed"
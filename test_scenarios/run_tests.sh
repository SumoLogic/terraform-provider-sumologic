#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PROVIDER_DIR="/Users/ashutosh.kumbhaj/terraform-provider-sumologic"
RESULTS_FILE="$PROVIDER_DIR/test_scenarios/TEST_RESULTS.md"

echo "# Multi-Notification Test Results" > "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"
echo "Date: $(date)" >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"

# Build provider
echo -e "${YELLOW}Building provider...${NC}"
cd "$PROVIDER_DIR"
go build -o terraform-provider-sumologic
echo -e "${GREEN}Provider built successfully${NC}"

# Check env vars
if [ -z "$SUMOLOGIC_ACCESSID" ] || [ -z "$SUMOLOGIC_ACCESSKEY" ] || [ -z "$SUMOLOGIC_BASE_URL" ]; then
    echo -e "${RED}Error: Set SUMOLOGIC_ACCESSID, SUMOLOGIC_ACCESSKEY, and SUMOLOGIC_BASE_URL${NC}"
    echo "Example:"
    echo '  export SUMOLOGIC_ACCESSID="su7YFH0jPodNk4"'
    echo '  export SUMOLOGIC_ACCESSKEY="98NX6py5iU6GOy7Dk9Zev2NKuwWRfqyWhzTklSdSX3pF7I0zCGISZ1hZgKinH4P6"'
    echo '  export SUMOLOGIC_BASE_URL="https://test3ak-www.ephemeral.in/api/"'
    exit 1
fi

run_scenario() {
    local scenario=$1
    local description=$2
    local dir="$PROVIDER_DIR/test_scenarios/$scenario"

    echo ""
    echo -e "${YELLOW}================================================================${NC}"
    echo -e "${YELLOW}Running: $scenario - $description${NC}"
    echo -e "${YELLOW}================================================================${NC}"
    echo "" >> "$RESULTS_FILE"
    echo "## $scenario: $description" >> "$RESULTS_FILE"
    echo "" >> "$RESULTS_FILE"

    cd "$dir"

    # Init
    terraform init -input=false > /dev/null 2>&1 || true

    # Plan
    echo -e "${YELLOW}  -> terraform plan${NC}"
    if terraform plan -input=false -detailed-exitcode > plan_output.txt 2>&1; then
        echo -e "${GREEN}  Plan: No changes (exit 0)${NC}"
        echo "- Plan: No changes" >> "$RESULTS_FILE"
    elif [ $? -eq 2 ]; then
        echo -e "${YELLOW}  Plan: Changes detected (exit 2)${NC}"
        echo "- Plan: Changes detected" >> "$RESULTS_FILE"
    else
        echo -e "${RED}  Plan: ERROR${NC}"
        echo "- Plan: ERROR" >> "$RESULTS_FILE"
        echo '```' >> "$RESULTS_FILE"
        tail -20 plan_output.txt >> "$RESULTS_FILE"
        echo '```' >> "$RESULTS_FILE"
        return 1
    fi

    # Apply
    echo -e "${YELLOW}  -> terraform apply${NC}"
    if terraform apply -input=false -auto-approve > apply_output.txt 2>&1; then
        echo -e "${GREEN}  Apply: SUCCESS${NC}"
        echo "- Apply: SUCCESS" >> "$RESULTS_FILE"

        # Get the resource ID
        local resource_id=$(terraform output -raw log_search_id 2>/dev/null || echo "unknown")
        echo "- Resource ID: $resource_id" >> "$RESULTS_FILE"
    else
        echo -e "${RED}  Apply: FAILED${NC}"
        echo "- Apply: FAILED" >> "$RESULTS_FILE"
        echo '```' >> "$RESULTS_FILE"
        tail -20 apply_output.txt >> "$RESULTS_FILE"
        echo '```' >> "$RESULTS_FILE"
        return 1
    fi

    # Plan again (check for drift)
    echo -e "${YELLOW}  -> terraform plan (drift check)${NC}"
    if terraform plan -input=false -detailed-exitcode > plan_drift_output.txt 2>&1; then
        echo -e "${GREEN}  Drift check: No drift (PASS)${NC}"
        echo "- Drift check: **NO DRIFT (PASS)**" >> "$RESULTS_FILE"
    elif [ $? -eq 2 ]; then
        echo -e "${RED}  Drift check: DRIFT DETECTED${NC}"
        echo "- Drift check: **DRIFT DETECTED (FAIL)**" >> "$RESULTS_FILE"
        echo '```' >> "$RESULTS_FILE"
        cat plan_drift_output.txt >> "$RESULTS_FILE"
        echo '```' >> "$RESULTS_FILE"
    else
        echo -e "${RED}  Drift check: ERROR${NC}"
        echo "- Drift check: ERROR" >> "$RESULTS_FILE"
        echo '```' >> "$RESULTS_FILE"
        tail -20 plan_drift_output.txt >> "$RESULTS_FILE"
        echo '```' >> "$RESULTS_FILE"
    fi

    # Import test
    echo -e "${YELLOW}  -> terraform import test${NC}"
    local resource_id=$(terraform output -raw log_search_id 2>/dev/null || echo "")
    if [ -n "$resource_id" ] && [ "$resource_id" != "unknown" ]; then
        # Remove from state and reimport
        terraform state rm sumologic_log_search.test > /dev/null 2>&1 || true
        if terraform import sumologic_log_search.test "$resource_id" > import_output.txt 2>&1; then
            echo -e "${GREEN}  Import: SUCCESS${NC}"
            echo "- Import: SUCCESS" >> "$RESULTS_FILE"

            # Plan after import (should show no changes)
            if terraform plan -input=false -detailed-exitcode > plan_after_import.txt 2>&1; then
                echo -e "${GREEN}  Plan after import: No drift (PASS)${NC}"
                echo "- Plan after import: **NO DRIFT (PASS)**" >> "$RESULTS_FILE"
            elif [ $? -eq 2 ]; then
                echo -e "${RED}  Plan after import: DRIFT DETECTED${NC}"
                echo "- Plan after import: **DRIFT DETECTED (FAIL)**" >> "$RESULTS_FILE"
                echo '```' >> "$RESULTS_FILE"
                cat plan_after_import.txt >> "$RESULTS_FILE"
                echo '```' >> "$RESULTS_FILE"
            fi
        else
            echo -e "${RED}  Import: FAILED${NC}"
            echo "- Import: FAILED" >> "$RESULTS_FILE"
            echo '```' >> "$RESULTS_FILE"
            tail -10 import_output.txt >> "$RESULTS_FILE"
            echo '```' >> "$RESULTS_FILE"
        fi
    fi

    # Destroy
    echo -e "${YELLOW}  -> terraform destroy${NC}"
    if terraform destroy -input=false -auto-approve > destroy_output.txt 2>&1; then
        echo -e "${GREEN}  Destroy: SUCCESS${NC}"
        echo "- Destroy: SUCCESS" >> "$RESULTS_FILE"
    else
        echo -e "${RED}  Destroy: FAILED${NC}"
        echo "- Destroy: FAILED" >> "$RESULTS_FILE"
    fi

    echo "" >> "$RESULTS_FILE"
}

# Parse arguments
SCENARIO=${1:-all}

case $SCENARIO in
    a|scenario_a)
        run_scenario "scenario_a" "FF ON + Old Terraform (singular notification)"
        ;;
    b|scenario_b)
        run_scenario "scenario_b" "FF OFF + New Terraform (plural notifications)"
        ;;
    c|scenario_c)
        echo -e "${YELLOW}Scenario C is a multi-step manual test. Running Step 1 (create with FF OFF):${NC}"
        run_scenario "scenario_c" "FF OFF → create → FF ON → plan (Step 1: create)"
        echo ""
        echo -e "${YELLOW}Now turn the BE FF ON, then run:${NC}"
        echo -e "${YELLOW}  cd $PROVIDER_DIR/test_scenarios/scenario_c && terraform plan${NC}"
        echo -e "${YELLOW}Expected: No changes (no drift)${NC}"
        ;;
    d|scenario_d)
        run_scenario "scenario_d" "FF ON + New Terraform + Multiple Notifications"
        ;;
    all)
        echo -e "${YELLOW}Running all scenarios. Make sure FF state is correct for each!${NC}"
        echo ""
        echo -e "${YELLOW}Scenarios A and D require FF ON.${NC}"
        echo -e "${YELLOW}Scenario B requires FF OFF.${NC}"
        echo -e "${YELLOW}Scenario C requires FF OFF first, then ON.${NC}"
        echo ""
        read -p "Is the BE FF currently ON? (y/n): " ff_state

        if [ "$ff_state" = "y" ]; then
            echo -e "${GREEN}Running Scenario A (FF ON + singular notification)...${NC}"
            run_scenario "scenario_a" "FF ON + Old Terraform (singular notification)"

            echo -e "${GREEN}Running Scenario D (FF ON + multiple notifications)...${NC}"
            run_scenario "scenario_d" "FF ON + New Terraform + Multiple Notifications"

            echo ""
            echo -e "${YELLOW}Now turn the BE FF OFF and run:${NC}"
            echo -e "${YELLOW}  $0 b${NC}"
            echo -e "${YELLOW}  $0 c${NC}"
        else
            echo -e "${GREEN}Running Scenario B (FF OFF + plural notifications)...${NC}"
            run_scenario "scenario_b" "FF OFF + New Terraform (plural notifications)"

            echo -e "${GREEN}Running Scenario C Step 1 (FF OFF + create)...${NC}"
            run_scenario "scenario_c" "FF OFF → create → FF ON → plan (Step 1: create)"

            echo ""
            echo -e "${YELLOW}Now turn the BE FF ON and run:${NC}"
            echo -e "${YELLOW}  $0 a${NC}"
            echo -e "${YELLOW}  $0 d${NC}"
            echo -e "${YELLOW}  cd $PROVIDER_DIR/test_scenarios/scenario_c && terraform plan${NC}"
        fi
        ;;
    *)
        echo "Usage: $0 [a|b|c|d|all]"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}Results written to: $RESULTS_FILE${NC}"

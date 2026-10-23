#!/bin/bash
# API-level tests that simulate Terraform provider behavior via curl
# These test the same request/response contract the provider uses

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_URL="https://epddata.stag.sumologic.net/api"
SESSION="Wn1IPLvktP6OacGdUraBToAS3KuSXPR9mQkVN6l9tlvuEHvhP0KI8H1Hdiri2tuZvi0sGDU4uvmNJZujWvdSnRn2VDHTWDEVvt34"
PARENT_ID="00000000192CC89F"

RESULTS_FILE="/Users/ashutosh.kumbhaj/terraform-provider-sumologic/test_scenarios/TEST_RESULTS.md"

api_call() {
  local method=$1
  local endpoint=$2
  local data=$3

  if [ "$method" = "GET" ]; then
    curl -s "$BASE_URL$endpoint" \
      -H 'content-type: application/json' \
      -H "apisession: $SESSION" \
      -b "ASID=$SESSION"
  elif [ "$method" = "POST" ]; then
    curl -s "$BASE_URL$endpoint" \
      -H 'content-type: application/json' \
      -H "apisession: $SESSION" \
      -b "ASID=$SESSION" \
      --data-raw "$data"
  elif [ "$method" = "PUT" ]; then
    curl -s -X PUT "$BASE_URL$endpoint" \
      -H 'content-type: application/json' \
      -H "apisession: $SESSION" \
      -b "ASID=$SESSION" \
      --data-raw "$data"
  elif [ "$method" = "DELETE" ]; then
    curl -s -X DELETE "$BASE_URL$endpoint" \
      -H 'content-type: application/json' \
      -H "apisession: $SESSION" \
      -b "ASID=$SESSION"
  fi
}

check_drift() {
  local id=$1
  local config_field=$2  # "notification" or "notifications"
  local expected_count=$3

  local response=$(api_call GET "/v1/logSearches/$id")
  local notif=$(echo "$response" | python3 -c "import sys,json; d=json.load(sys.stdin); print(json.dumps(d.get('schedule',{}).get('notification')))" 2>/dev/null)
  local notifs=$(echo "$response" | python3 -c "import sys,json; d=json.load(sys.stdin); print(json.dumps(d.get('schedule',{}).get('notifications')))" 2>/dev/null)

  echo "    API response: notification=$notif"
  echo "    API response: notifications=$notifs"

  if [ "$config_field" = "notification" ]; then
    if [ "$notif" != "null" ] && [ "$notif" != "None" ]; then
      echo -e "    ${GREEN}NO DRIFT: API returns notification (singular) matching config${NC}"
      return 0
    else
      echo -e "    ${RED}DRIFT: Config uses 'notification' but API returned null for it${NC}"
      return 1
    fi
  elif [ "$config_field" = "notifications" ]; then
    if [ "$notifs" != "null" ] && [ "$notifs" != "None" ]; then
      local count=$(echo "$response" | python3 -c "import sys,json; d=json.load(sys.stdin); n=d.get('schedule',{}).get('notifications'); print(len(n) if n else 0)" 2>/dev/null)
      if [ "$count" = "$expected_count" ]; then
        echo -e "    ${GREEN}NO DRIFT: API returns $count notifications matching config${NC}"
        return 0
      else
        echo -e "    ${RED}DRIFT: Config has $expected_count notifications but API returns $count${NC}"
        return 1
      fi
    else
      echo -e "    ${RED}DRIFT: Config uses 'notifications' but API returned null for it${NC}"
      return 1
    fi
  fi
}

echo "# Multi-Notification API Test Results" > "$RESULTS_FILE"
echo "Date: $(date)" >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"

# ============================================================
echo -e "${BLUE}============================================================${NC}"
echo -e "${BLUE}SCENARIO A: FF ON + Old Terraform (singular notification)${NC}"
echo -e "${BLUE}============================================================${NC}"
echo "" >> "$RESULTS_FILE"
echo "## Scenario A: FF ON + Old Terraform (singular notification)" >> "$RESULTS_FILE"

echo -e "${YELLOW}  Creating SS with singular notification...${NC}"
RESPONSE_A=$(api_call POST "/v1/logSearches" '{
  "queryString": "*",
  "timeRange": {"type":"BeginBoundedTimeRange","from":{"type":"RelativeTimeRangeBoundary","relativeTime":"-900s"}},
  "runByReceiptTime": false,
  "parsingMode": "Manual",
  "name": "TF_Scenario_A_SingularNotif",
  "description": "Scenario A: singular notification",
  "parentId": "'"$PARENT_ID"'",
  "schedule": {
    "parseableTimeRange": {"type":"BeginBoundedTimeRange","from":{"type":"RelativeTimeRangeBoundary","relativeTime":"-900s"}},
    "timeZone": "Asia/Kolkata",
    "scheduleType": "15Minutes",
    "cronExpression": "0 0/15 * * * ? *",
    "notification": {"taskType":"EmailSearchNotificationSyncDefinition","toList":["scenario_a@sumologic.com"],"subjectTemplate":"Alert: {{SearchName}}","includeQuery":true,"includeResultSet":true,"includeHistogram":true,"includeCsvAttachment":false},
    "muteErrorEmails": false
  }
}')

ID_A=$(echo "$RESPONSE_A" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('id','ERROR'))" 2>/dev/null)
if [ "$ID_A" = "ERROR" ] || [ -z "$ID_A" ]; then
  echo -e "  ${RED}CREATE FAILED${NC}"
  echo "$RESPONSE_A"
  echo "- Create: FAILED" >> "$RESULTS_FILE"
  echo '```' >> "$RESULTS_FILE"
  echo "$RESPONSE_A" >> "$RESULTS_FILE"
  echo '```' >> "$RESULTS_FILE"
else
  echo -e "  ${GREEN}CREATE SUCCESS: ID=$ID_A${NC}"
  echo "- Create: SUCCESS (ID: $ID_A)" >> "$RESULTS_FILE"

  echo -e "${YELLOW}  Checking for drift (simulating terraform plan)...${NC}"
  if check_drift "$ID_A" "notification" "1"; then
    echo "- Drift check: **NO DRIFT (PASS)**" >> "$RESULTS_FILE"
  else
    echo "- Drift check: **DRIFT DETECTED (FAIL)**" >> "$RESULTS_FILE"
  fi

  # Cleanup
  api_call DELETE "/v1/logSearches/$ID_A" > /dev/null 2>&1
  echo -e "  ${GREEN}Cleaned up${NC}"
  echo "- Cleanup: Done" >> "$RESULTS_FILE"
fi

# ============================================================
echo ""
echo -e "${BLUE}============================================================${NC}"
echo -e "${BLUE}SCENARIO B: Current FF state + plural notifications${NC}"
echo -e "${BLUE}============================================================${NC}"
echo "" >> "$RESULTS_FILE"
echo "## Scenario B: Current FF state + New Terraform (plural notifications)" >> "$RESULTS_FILE"

echo -e "${YELLOW}  Creating SS with plural notifications (2 emails)...${NC}"
RESPONSE_B=$(api_call POST "/v1/logSearches" '{
  "queryString": "*",
  "timeRange": {"type":"BeginBoundedTimeRange","from":{"type":"RelativeTimeRangeBoundary","relativeTime":"-900s"}},
  "runByReceiptTime": false,
  "parsingMode": "Manual",
  "name": "TF_Scenario_B_PluralNotif",
  "description": "Scenario B: plural notifications",
  "parentId": "'"$PARENT_ID"'",
  "schedule": {
    "parseableTimeRange": {"type":"BeginBoundedTimeRange","from":{"type":"RelativeTimeRangeBoundary","relativeTime":"-900s"}},
    "timeZone": "Asia/Kolkata",
    "scheduleType": "15Minutes",
    "cronExpression": "0 0/15 * * * ? *",
    "notifications": [
      {"taskType":"EmailSearchNotificationSyncDefinition","toList":["scenario_b_1@sumologic.com"],"subjectTemplate":"Alert 1: {{SearchName}}","includeQuery":true,"includeResultSet":true,"includeHistogram":true,"includeCsvAttachment":false},
      {"taskType":"EmailSearchNotificationSyncDefinition","toList":["scenario_b_2@sumologic.com"],"subjectTemplate":"Alert 2: {{SearchName}}","includeQuery":false,"includeResultSet":true,"includeHistogram":false,"includeCsvAttachment":false}
    ],
    "muteErrorEmails": false
  }
}')

ID_B=$(echo "$RESPONSE_B" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('id','ERROR'))" 2>/dev/null)
if [ "$ID_B" = "ERROR" ] || [ -z "$ID_B" ]; then
  echo -e "  ${RED}CREATE FAILED${NC}"
  echo "$RESPONSE_B"
  echo "- Create: FAILED" >> "$RESULTS_FILE"
  echo '```' >> "$RESULTS_FILE"
  echo "$RESPONSE_B" >> "$RESULTS_FILE"
  echo '```' >> "$RESULTS_FILE"
else
  echo -e "  ${GREEN}CREATE SUCCESS: ID=$ID_B${NC}"
  echo "- Create: SUCCESS (ID: $ID_B)" >> "$RESULTS_FILE"

  echo -e "${YELLOW}  Checking for drift (simulating terraform plan)...${NC}"
  if check_drift "$ID_B" "notifications" "2"; then
    echo "- Drift check: **NO DRIFT (PASS)**" >> "$RESULTS_FILE"
  else
    echo "- Drift check: **DRIFT DETECTED (FAIL)**" >> "$RESULTS_FILE"
    # Show what API actually returned
    FULL_RESP=$(api_call GET "/v1/logSearches/$ID_B")
    echo '```json' >> "$RESULTS_FILE"
    echo "$FULL_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); s=d.get('schedule',{}); print(json.dumps({'notification':s.get('notification'),'notifications':s.get('notifications')}, indent=2))" >> "$RESULTS_FILE"
    echo '```' >> "$RESULTS_FILE"
  fi

  # Don't cleanup - keep for scenario C
  echo "  Keeping ID=$ID_B for scenario C testing"
  echo "- Kept alive for Scenario C (ID: $ID_B)" >> "$RESULTS_FILE"
fi

# ============================================================
echo ""
echo -e "${BLUE}============================================================${NC}"
echo -e "${BLUE}SCENARIO C: Read existing SS created before/after FF toggle${NC}"
echo -e "${BLUE}============================================================${NC}"
echo "" >> "$RESULTS_FILE"
echo "## Scenario C: Read existing SS (pre-existing: SearchBeforeFF)" >> "$RESULTS_FILE"

echo -e "${YELLOW}  Reading existing SS 'SearchBeforeFF' (000000001A0D0E3B)...${NC}"
RESPONSE_C=$(api_call GET "/v1/logSearches/000000001A0D0E3B")
echo "$RESPONSE_C" | python3 -c "
import sys, json
d = json.load(sys.stdin)
s = d.get('schedule', {})
print(f'  name: {d.get(\"name\")}')
print(f'  notification: {json.dumps(s.get(\"notification\"), indent=4) if s.get(\"notification\") else \"null\"}')
print(f'  notifications: {json.dumps(s.get(\"notifications\"), indent=4) if s.get(\"notifications\") else \"null\"}')
" 2>/dev/null

echo "- Existing SS 'SearchBeforeFF' read:" >> "$RESULTS_FILE"
echo "$RESPONSE_C" | python3 -c "
import sys, json
d = json.load(sys.stdin)
s = d.get('schedule', {})
print(f'  - notification: {\"present\" if s.get(\"notification\") else \"null\"}')
print(f'  - notifications: {\"present (\" + str(len(s.get(\"notifications\"))) + \" items)\" if s.get(\"notifications\") else \"null\"}')
" >> "$RESULTS_FILE" 2>/dev/null

# ============================================================
echo ""
echo -e "${BLUE}============================================================${NC}"
echo -e "${BLUE}SCENARIO D: Create with multiple notifications (happy path)${NC}"
echo -e "${BLUE}============================================================${NC}"
echo "" >> "$RESULTS_FILE"
echo "## Scenario D: FF ON + New Terraform + Multiple Notifications (Happy Path)" >> "$RESULTS_FILE"

echo -e "${YELLOW}  Creating SS with 2 different notification types...${NC}"
RESPONSE_D=$(api_call POST "/v1/logSearches" '{
  "queryString": "*",
  "timeRange": {"type":"BeginBoundedTimeRange","from":{"type":"RelativeTimeRangeBoundary","relativeTime":"-900s"}},
  "runByReceiptTime": false,
  "parsingMode": "Manual",
  "name": "TF_Scenario_D_MultiNotif_HappyPath",
  "description": "Scenario D: multi notifications happy path",
  "parentId": "'"$PARENT_ID"'",
  "schedule": {
    "parseableTimeRange": {"type":"BeginBoundedTimeRange","from":{"type":"RelativeTimeRangeBoundary","relativeTime":"-900s"}},
    "timeZone": "Asia/Kolkata",
    "scheduleType": "15Minutes",
    "cronExpression": "0 0/15 * * * ? *",
    "notifications": [
      {"taskType":"EmailSearchNotificationSyncDefinition","toList":["scenario_d_team_a@sumologic.com"],"subjectTemplate":"Alert 1: {{SearchName}}","includeQuery":true,"includeResultSet":true,"includeHistogram":true,"includeCsvAttachment":false},
      {"taskType":"EmailSearchNotificationSyncDefinition","toList":["scenario_d_team_b@sumologic.com"],"subjectTemplate":"Alert 2: {{SearchName}}","includeQuery":true,"includeResultSet":true,"includeHistogram":false,"includeCsvAttachment":true}
    ],
    "muteErrorEmails": false,
    "threshold": {"thresholdType":"group","operator":"gt","count":10}
  }
}')

ID_D=$(echo "$RESPONSE_D" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('id','ERROR'))" 2>/dev/null)
if [ "$ID_D" = "ERROR" ] || [ -z "$ID_D" ]; then
  echo -e "  ${RED}CREATE FAILED${NC}"
  echo "$RESPONSE_D"
  echo "- Create: FAILED" >> "$RESULTS_FILE"
  echo '```' >> "$RESULTS_FILE"
  echo "$RESPONSE_D" >> "$RESULTS_FILE"
  echo '```' >> "$RESULTS_FILE"
else
  echo -e "  ${GREEN}CREATE SUCCESS: ID=$ID_D${NC}"
  echo "- Create: SUCCESS (ID: $ID_D)" >> "$RESULTS_FILE"

  echo -e "${YELLOW}  Checking for drift...${NC}"
  if check_drift "$ID_D" "notifications" "2"; then
    echo "- Drift check: **NO DRIFT (PASS)**" >> "$RESULTS_FILE"
  else
    echo "- Drift check: **DRIFT DETECTED (FAIL)**" >> "$RESULTS_FILE"
    FULL_RESP=$(api_call GET "/v1/logSearches/$ID_D")
    echo '```json' >> "$RESULTS_FILE"
    echo "$FULL_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); s=d.get('schedule',{}); print(json.dumps({'notification':s.get('notification'),'notifications':s.get('notifications')}, indent=2))" >> "$RESULTS_FILE"
    echo '```' >> "$RESULTS_FILE"
  fi

  # Test UPDATE (simulating terraform apply on existing resource)
  echo -e "${YELLOW}  Testing UPDATE (add a 3rd notification)...${NC}"
  UPDATE_RESP=$(api_call PUT "/v1/logSearches/$ID_D" '{
    "queryString": "*",
    "timeRange": {"type":"BeginBoundedTimeRange","from":{"type":"RelativeTimeRangeBoundary","relativeTime":"-900s"}},
    "runByReceiptTime": false,
    "parsingMode": "Manual",
    "name": "TF_Scenario_D_MultiNotif_HappyPath",
    "description": "Scenario D: updated to 3 notifications",
    "parentId": "'"$PARENT_ID"'",
    "schedule": {
      "parseableTimeRange": {"type":"BeginBoundedTimeRange","from":{"type":"RelativeTimeRangeBoundary","relativeTime":"-900s"}},
      "timeZone": "Asia/Kolkata",
      "scheduleType": "15Minutes",
      "cronExpression": "0 0/15 * * * ? *",
      "notifications": [
        {"taskType":"EmailSearchNotificationSyncDefinition","toList":["scenario_d_team_a@sumologic.com"],"subjectTemplate":"Alert 1: {{SearchName}}","includeQuery":true,"includeResultSet":true,"includeHistogram":true,"includeCsvAttachment":false},
        {"taskType":"EmailSearchNotificationSyncDefinition","toList":["scenario_d_team_b@sumologic.com"],"subjectTemplate":"Alert 2: {{SearchName}}","includeQuery":true,"includeResultSet":true,"includeHistogram":false,"includeCsvAttachment":true},
        {"taskType":"EmailSearchNotificationSyncDefinition","toList":["scenario_d_team_c@sumologic.com"],"subjectTemplate":"Alert 3: {{SearchName}}","includeQuery":false,"includeResultSet":true,"includeHistogram":false,"includeCsvAttachment":false}
      ],
      "muteErrorEmails": false,
      "threshold": {"thresholdType":"group","operator":"gt","count":10}
    }
  }')

  UPDATE_ERR=$(echo "$UPDATE_RESP" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('errors',''))" 2>/dev/null)
  if [ -z "$UPDATE_ERR" ] || [ "$UPDATE_ERR" = "" ] || [ "$UPDATE_ERR" = "None" ]; then
    echo -e "  ${GREEN}UPDATE SUCCESS${NC}"
    echo "- Update (2→3 notifications): SUCCESS" >> "$RESULTS_FILE"
    check_drift "$ID_D" "notifications" "3"
  else
    echo -e "  ${RED}UPDATE FAILED: $UPDATE_ERR${NC}"
    echo "- Update: FAILED" >> "$RESULTS_FILE"
  fi

  # Cleanup
  api_call DELETE "/v1/logSearches/$ID_D" > /dev/null 2>&1
  echo -e "  ${GREEN}Cleaned up${NC}"
  echo "- Cleanup: Done" >> "$RESULTS_FILE"
fi

# Cleanup scenario B
if [ -n "$ID_B" ] && [ "$ID_B" != "ERROR" ]; then
  api_call DELETE "/v1/logSearches/$ID_B" > /dev/null 2>&1
  echo -e "\n  ${GREEN}Cleaned up Scenario B (ID=$ID_B)${NC}"
fi

echo ""
echo -e "${GREEN}============================================================${NC}"
echo -e "${GREEN}ALL TESTS COMPLETE. Results: $RESULTS_FILE${NC}"
echo -e "${GREEN}============================================================${NC}"

cat "$RESULTS_FILE"

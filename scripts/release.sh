#!/usr/bin/env bash
set -euo pipefail

DRY_RUN=false
if [[ "${1:-}" == "--dry-run" ]]; then
  DRY_RUN=true
fi

REPO="SumoLogic/terraform-provider-sumologic"
REGISTRY_URL="https://registry.terraform.io/v1/providers/SumoLogic/sumologic"
BOLD='\033[1m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

step_num=0
total_steps=6

print_step() {
  step_num=$((step_num + 1))
  echo ""
  echo -e "${BOLD}━━━ Step ${step_num} of ${total_steps}: $1 ━━━${NC}"
}

confirm() {
  local prompt="${1:-Continue?}"
  echo ""
  echo -en "${YELLOW}${prompt} [y/N] ${NC}"
  read -r answer
  [[ "$answer" =~ ^[Yy]$ ]] || { echo "Aborted."; exit 1; }
}

info() { echo -e "${CYAN}ℹ ${NC}$1"; }
success() { echo -e "${GREEN}✓ ${NC}$1"; }
warn() { echo -e "${YELLOW}⚠ ${NC}$1"; }
fail() { echo -e "${RED}✗ ${NC}$1"; exit 1; }

# ── Prerequisites ──────────────────────────────────────────────────────────
command -v gh >/dev/null 2>&1 || fail "gh CLI is required. Install: https://cli.github.com"
command -v git >/dev/null 2>&1 || fail "git is required."
command -v jq >/dev/null 2>&1 || fail "jq is required. Install: brew install jq"

if [[ ! -f CHANGELOG.md ]]; then
  fail "CHANGELOG.md not found. Run this script from the repo root."
fi

# Work on a copy of origin/master's CHANGELOG so the audit, version detection and the
# release PR all see the same content, regardless of the local checkout's state.
git fetch -q origin master || fail "Could not fetch origin/master."
WORK_CL=$(mktemp)
trap 'rm -f "$WORK_CL"' EXIT
git show origin/master:CHANGELOG.md > "$WORK_CL" || fail "CHANGELOG.md not found on origin/master."

# ── Verify team membership ────────────────────────────────────────────────
gh_user=$(gh api user --jq '.login' 2>/dev/null) || fail "Not authenticated with gh CLI. Run: gh auth login"
if ! gh api "orgs/SumoLogic/teams/tf-platform-team/members" --jq '.[].login' 2>/dev/null | grep -qx "$gh_user"; then
  fail "${gh_user} is not a member of @SumoLogic/tf-platform-team. Only tf-platform-team members can run releases."
fi
success "Authenticated as ${gh_user} (tf-platform-team member)"

# ── Gather context ─────────────────────────────────────────────────────────
last_tag=$(gh api "repos/${REPO}/releases/latest" --jq '.tag_name' 2>/dev/null) || fail "Could not fetch latest release."
last_version="${last_tag#v}"
last_release_date=$(gh api "repos/${REPO}/releases/latest" --jq '.published_at' 2>/dev/null)
last_release_short=$(echo "$last_release_date" | cut -c1-10)

echo ""
echo -e "${BOLD}Terraform Provider Release Tool${NC}"
if [[ "$DRY_RUN" == true ]]; then
  echo -e "${YELLOW}DRY RUN — Steps 1-2 only, no changes will be made${NC}"
fi
echo -e "Repository:     ${REPO}"
echo -e "Latest release: ${GREEN}${last_tag}${NC} (${last_release_short})"
echo ""

# ── Step 1: Audit changelog completeness ───────────────────────────────────
print_step "Audit changelog completeness"
info "Checking merged PRs since ${last_tag} against CHANGELOG.md..."

unreleased_section=$(sed -n '/^## X\.Y\.Z/,/^## [0-9]/p' "$WORK_CL" | sed '$d')

echo ""
echo -e "${BOLD}Current Unreleased section:${NC}"
echo "$unreleased_section"
echo ""

merged_prs=$(gh pr list \
  --repo "$REPO" \
  --state merged \
  --base master \
  --search "merged:>=${last_release_short}" \
  --limit 100 \
  --json number,title,mergedAt,author,labels \
  --jq '.' 2>/dev/null || echo "[]")

missing_pr_json="[]"
if [[ -n "$merged_prs" && "$merged_prs" != "[]" ]]; then
  # Filter out:
  #   - PRs labeled changelog/no-changelog
  #   - Previous release PRs (title starts with "Release " or "Releasing ")
  # Then check which remaining PRs are missing from the Unreleased section.
  # Match by PR number (#NNN) OR by checking if the PR title content appears in the section.
  missing_pr_json=$(echo "$merged_prs" | jq --arg section "$unreleased_section" '
    [.[] |
      select((.labels | map(.name) | index("changelog/no-changelog")) == null) |
      select(.title | test("^(Release|Releasing) [0-9]") | not) |
      select(
        (.number as $n | ($section | test("#" + ($n | tostring) + "\\b") | not)) and
        (.title as $t |
          # Extract meaningful keywords (3+ char words, skip ticket prefixes like BE-1234/SUMO-1234)
          [$t | scan("[a-zA-Z_][a-zA-Z0-9_]{2,}") | select(test("^(Add|Added|Fix|Fixed|Update|Updated|Support|Remove|New|Set|Run|Get|Use)$"; "i") | not)] |
          if length == 0 then true
          else any(. as $word | $section | test($word; "i")) | not
          end
        )
      )
    ]
  ')
  skipped_count=$(echo "$merged_prs" | jq '[.[] |
    select(
      ((.labels | map(.name) | index("changelog/no-changelog")) != null) or
      (.title | test("^(Release|Releasing) [0-9]"))
    )] | length')
  if [[ "$skipped_count" -gt 0 ]]; then
    info "Skipped ${skipped_count} PR(s) (changelog/no-changelog or release PRs)."
  fi
  missing_count=$(echo "$missing_pr_json" | jq 'length')

  if [[ "$missing_count" -gt 0 ]]; then
    warn "${missing_count} merged PR(s) are NOT referenced in the Unreleased section:"
    echo ""
    echo "$missing_pr_json" | jq -r '.[] | "  #\(.number) \(.title) (@\(.author.login))"'
    echo ""

    # Auto-draft changelog entries from PR metadata
    info "Generating draft changelog entries..."
    draft_entries=$(echo "$missing_pr_json" | jq -r '
      def categorize:
        (.labels | map(.name)) as $labels |
        if ($labels | index("breaking-change") != null) then "BREAKING CHANGES"
        elif ($labels | index("new-resource") != null) or ($labels | index("new-data-source") != null) then "FEATURES"
        elif ($labels | index("bug") != null) then "BUG FIXES"
        elif ($labels | index("documentation") != null) then "DOCS"
        elif ($labels | index("enhancement") != null) then "ENHANCEMENTS"
        else "ENHANCEMENTS"
        end;
      group_by(categorize) | sort_by(.[0] | categorize |
        if . == "BREAKING CHANGES" then 0
        elif . == "FEATURES" then 1
        elif . == "ENHANCEMENTS" then 2
        elif . == "BUG FIXES" then 3
        elif . == "DOCS" then 4
        else 5 end
      ) | map(
        (.[0] | categorize) as $cat |
        "\($cat):\n" + (map("* \(.title) ([#\(.number)](https://github.com/SumoLogic/terraform-provider-sumologic/pull/\(.number)))") | join("\n"))
      ) | join("\n\n")
    ')

    echo ""
    echo -e "${BOLD}Draft changelog entries:${NC}"
    echo "$draft_entries"
    echo ""
    warn "These are best-guess entries based on PR titles and labels."
    warn "Some PRs may be infra/CI and don't need entries — edit as needed."
    echo ""
    echo -en "${YELLOW}Add these entries to the release changelog? [y/n/e(dit manually)] ${NC}"
    read -r add_answer

    if [[ "$add_answer" =~ ^[Yy]$ ]]; then
      # Insert draft entries after the Unreleased header line.
      # awk + temp file (not sed -i) so this behaves the same with BSD and GNU sed.
      cl_tmp=$(mktemp)
      DRAFT="$draft_entries" awk '{print} /^## X\.Y\.Z \(Unreleased\)/{print ENVIRON["DRAFT"]}' "$WORK_CL" > "$cl_tmp"
      mv "$cl_tmp" "$WORK_CL"
      success "Draft entries added."
      unreleased_section=$(sed -n '/^## X\.Y\.Z/,/^## [0-9]/p' "$WORK_CL" | sed '$d')
      echo ""
      echo -e "${BOLD}Updated Unreleased section:${NC}"
      echo "$unreleased_section"
    elif [[ "$add_answer" =~ ^[Ee]$ ]]; then
      info "Opening the Unreleased changelog copy for manual editing (applied to the release PR)..."
      ${EDITOR:-vim} "$WORK_CL"
      unreleased_section=$(sed -n '/^## X\.Y\.Z/,/^## [0-9]/p' "$WORK_CL" | sed '$d')
    else
      info "Skipping — proceeding with existing changelog."
    fi

    confirm "Changelog looks complete enough to proceed?"
  else
    success "All merged PRs appear to be referenced in the changelog."
  fi
else
  info "No merged PRs found since last release (or couldn't query)."
fi

# ── Step 2: Determine version ──────────────────────────────────────────────
print_step "Determine version number"

# Guard: refuse to cut an empty release. Count bullet entries (any category, or headerless
# older-style bullets), ignoring the placeholder line.
entry_count=$(echo "$unreleased_section" | grep -E '^[[:space:]]*[*-][[:space:]]' | grep -vc 'Add new change notes here' || true)
if [[ "$entry_count" -eq 0 ]]; then
  info "The Unreleased section of CHANGELOG.md has no entries — nothing to release."
  exit 0
fi
info "Unreleased entries: ${entry_count}"

IFS='.' read -r major minor patch <<< "$last_version"

# Auto-detect recommended bump from Unreleased section categories
has_breaking=$(echo "$unreleased_section" | grep -ci "BREAKING CHANGE" || true)
has_features=$(echo "$unreleased_section" | grep -ci "^FEATURES:" || true)

if [[ "$has_breaking" -gt 0 ]]; then
  recommended="major"
  suggested="$((major + 1)).0.0"
  reason="BREAKING CHANGES detected"
  echo ""
  echo -e "${RED}${BOLD}⚠ BREAKING CHANGES detected in the Unreleased section.${NC}"
  echo -e "${RED}Major version bumps affect all downstream consumers and require careful coordination.${NC}"
  echo ""
  echo "Before proceeding, verify:"
  echo "  1. Breaking changes have been communicated to users via deprecation notices"
  echo "  2. Migration guide is documented"
  echo "  3. Team lead has approved the major version bump"
  echo ""
  confirm "Has this major version bump been explicitly approved by the team?"
elif [[ "$has_features" -gt 0 ]]; then
  recommended="minor"
  suggested="$major.$((minor + 1)).0"
  reason="New FEATURES (resources/data sources) detected"
else
  recommended="patch"
  suggested="$major.$minor.$((patch + 1))"
  reason="Bug fixes / enhancements only"
fi

info "Last version: ${last_version}"
echo ""
echo "  Detected sections → ${reason}"
echo ""
echo -e "  ${GREEN}→ Recommended: ${suggested} (${recommended})${NC}"
if [[ "$recommended" != "major" ]]; then
  echo "    Major bump → $((major + 1)).0.0 (breaking changes)"
fi
if [[ "$recommended" != "minor" ]]; then
  echo "    Minor bump → $major.$((minor + 1)).0 (new resources/features)"
fi
if [[ "$recommended" != "patch" ]]; then
  echo "    Patch bump → $major.$minor.$((patch + 1)) (bug fixes/enhancements)"
fi
echo ""
echo -en "${YELLOW}Enter new version (without 'v' prefix) [${suggested}]: ${NC}"
read -r new_version
new_version="${new_version:-$suggested}"

if [[ ! "$new_version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  fail "Invalid version format. Expected: X.Y.Z"
fi

# Guard: must be strictly greater than the latest release (numeric compare; portable, no sort -V)
IFS='.' read -r nv_major nv_minor nv_patch <<< "$new_version"
if (( nv_major * 1000000 + nv_minor * 1000 + nv_patch <= major * 1000000 + minor * 1000 + patch )); then
  fail "Version ${new_version} must be greater than the latest release (${last_version})."
fi

# Guard: block major bumps that weren't already approved above
IFS='.' read -r new_major _ _ <<< "$new_version"
if [[ "$new_major" -gt "$major" && "$recommended" != "major" ]]; then
  echo ""
  echo -e "${RED}${BOLD}⚠ You entered a major version bump (${last_version} → ${new_version}) but no BREAKING CHANGES were detected.${NC}"
  echo ""
  echo "Before proceeding, verify:"
  echo "  1. Breaking changes have been communicated to users via deprecation notices"
  echo "  2. Migration guide is documented"
  echo "  3. Team lead has approved the major version bump"
  echo ""
  confirm "Has this major version bump been explicitly approved by the team?"
fi

new_tag="v${new_version}"
release_date=$(date +"%B %-d, %Y")

echo ""
info "Version: ${new_version}"
info "Tag:     ${new_tag}"
info "Date:    ${release_date}"

if [[ "$DRY_RUN" == true ]]; then
  echo ""
  success "Dry run complete. Steps 1-2 validated successfully."
  echo ""
  echo "Summary:"
  echo "  Changelog audit:  done"
  echo "  Version detected: ${new_version} (${recommended} bump)"
  echo "  Would create:     branch release-${new_version}, PR 'Release ${new_version}', tag ${new_tag}"
  echo ""
  echo "Run without --dry-run to execute the full release."
  exit 0
fi

confirm "Proceed with this version?"

# ── Step 3: Update CHANGELOG.md and create PR ─────────────────────────────
print_step "Update CHANGELOG.md and create PR"

git checkout --no-track -b "release-${new_version}" origin/master
cp "$WORK_CL" CHANGELOG.md

# Rewrite CHANGELOG to match the established pattern:
#   1. Fresh Unreleased header with placeholder
#   2. New versioned section with the items that were under Unreleased
#   3. Rest of the file unchanged
#
# Real example (PR #907):
#   ## X.Y.Z (Unreleased)
#   * Add new change notes here
#
#   ## 3.3.4 (October 6, 2026)
#   ENHANCEMENTS:
#   * ...

# Extract everything between Unreleased header and next version header
# awk (not sed) so this behaves the same on BSD and GNU: the Unreleased header is line 1,
# so a sed range starting at line 1 never sees its end pattern.
unreleased_body=$(awk '/^## X\.Y\.Z/{f=1; next} /^## [0-9]/{f=0} f' CHANGELOG.md | sed '/^[[:space:]]*$/d; /^\* Add new change notes here$/d')
rest_of_file=$(sed -n '/^## [0-9]/,$p' CHANGELOG.md)

{
  echo "## X.Y.Z (Unreleased)"
  echo "* Add new change notes here"
  echo ""
  echo "## ${new_version} (${release_date})"
  if [[ -n "$unreleased_body" ]]; then
    echo "$unreleased_body"
  fi
  echo ""
  echo "$rest_of_file"
} > CHANGELOG.md

echo ""
echo -e "${BOLD}CHANGELOG.md diff:${NC}"
git diff CHANGELOG.md
echo ""

confirm "Create PR with this CHANGELOG update?"

git add CHANGELOG.md
git commit -m "Release ${new_version}"
git push -u origin "release-${new_version}:release-${new_version}"

pr_url=$(gh pr create \
  --repo "$REPO" \
  --title "Release ${new_version}" \
  --body "Updates CHANGELOG.md for release ${new_version}." \
  --base master \
  --head "release-${new_version}")

success "PR created: ${pr_url}"

# ── Step 4: Wait for PR merge, then tag ────────────────────────────────────
print_step "Wait for PR merge, then tag"

info "Waiting for the PR to be merged..."
info "Review and merge the PR, then come back here."
echo ""

while true; do
  pr_state=$(gh pr view "release-${new_version}" --repo "$REPO" --json state --jq '.state' 2>/dev/null || echo "UNKNOWN")
  if [[ "$pr_state" == "MERGED" ]]; then
    success "PR merged."
    break
  fi
  echo -en "\r  PR state: ${pr_state}. Checking again in 15s... (Ctrl+C to abort)"
  sleep 15
done

info "Pulling latest master and tagging..."
git checkout master
git pull origin master

merge_sha=$(git rev-parse HEAD)
info "Tagging ${new_tag} at commit ${merge_sha:0:12}"
confirm "Push tag ${new_tag} to origin? This triggers the release build."

git tag "$new_tag"
git push origin "$new_tag"
success "Tag ${new_tag} pushed."

# ── Step 5: Wait for goreleaser, then publish ──────────────────────────────
print_step "Wait for goreleaser, then publish release"

info "Waiting for goreleaser GitHub Action to start..."
sleep 10

run_id=""
for i in $(seq 1 12); do
  run_id=$(gh run list \
    --repo "$REPO" \
    --workflow release.yml \
    --limit 1 \
    --json databaseId,headBranch,status \
    --jq ".[] | select(.headBranch == \"${new_tag}\") | .databaseId" 2>/dev/null || echo "")
  if [[ -n "$run_id" ]]; then
    break
  fi
  echo -en "\r  Waiting for workflow run to appear... (attempt ${i}/12)"
  sleep 10
done

if [[ -z "$run_id" ]]; then
  warn "Could not find goreleaser run. Check manually:"
  info "https://github.com/${REPO}/actions/workflows/release.yml"
  confirm "Continue anyway? (Only if the run completed successfully)"
else
  info "Goreleaser run: https://github.com/${REPO}/actions/runs/${run_id}"
  info "Waiting for completion (typically <10 min)..."

  if gh run watch "$run_id" --repo "$REPO" --exit-status 2>/dev/null; then
    success "Goreleaser completed successfully."
  else
    fail "Goreleaser failed. Check: https://github.com/${REPO}/actions/runs/${run_id}"
  fi
fi

# Fetch the draft release and publish it
draft_release_id=$(gh api "repos/${REPO}/releases" --jq ".[] | select(.tag_name == \"${new_tag}\" and .draft == true) | .id" 2>/dev/null || echo "")

if [[ -z "$draft_release_id" ]]; then
  warn "No draft release found for ${new_tag}."
  info "Check: https://github.com/${REPO}/releases"
  confirm "Skip auto-publish and verify manually?"
else
  # Extract notes between version header and next version header, strip the header line and blank lines
  release_notes=$(sed -n "/^## ${new_version}/,/^## [0-9]/p" CHANGELOG.md | sed '$d' | tail -n +2 | sed '/^[[:space:]]*$/d')

  echo ""
  echo -e "${BOLD}Release notes to publish:${NC}"
  echo "$release_notes"
  echo ""

  confirm "Publish this release as ${new_tag}?"

  gh api \
    --method PATCH \
    "repos/${REPO}/releases/${draft_release_id}" \
    -f body="$release_notes" \
    -F draft=false \
    -F make_latest=true \
    --silent

  success "Release ${new_tag} published: https://github.com/${REPO}/releases/tag/${new_tag}"
fi

# ── Step 6: Verify Terraform Registry ──────────────────────────────────────
print_step "Verify Terraform Registry"

info "Checking Terraform Registry for version ${new_version}..."
info "This can take up to 30 minutes."

for i in $(seq 1 20); do
  registry_version=$(curl -sf "$REGISTRY_URL" | jq -r '.version' 2>/dev/null || echo "unknown")
  if [[ "$registry_version" == "$new_version" ]]; then
    success "Terraform Registry updated to ${new_version}!"
    break
  fi
  echo -en "\r  Registry shows: ${registry_version}. Checking again in 90s... (attempt ${i}/20)"
  sleep 90
done

if [[ "$registry_version" != "$new_version" ]]; then
  warn "Registry has not updated after polling. Try redelivering webhooks:"
  info "1. Open https://github.com/${REPO}/settings/hooks"
  info "2. Edit the Terraform Registry webhook"
  info "3. Go to Recent Deliveries tab"
  info "4. Redeliver the recent messages"
  echo ""
  info "Or just wait — it sometimes takes up to 30 min."
fi

# ── Done ───────────────────────────────────────────────────────────────────
echo ""
echo -e "${BOLD}${GREEN}━━━ Release ${new_tag} complete! ━━━${NC}"
echo ""
echo "Remaining manual steps:"
echo "  1. Post status update in #terraform-team Slack thread"
echo "  2. Close the release Jira ticket"

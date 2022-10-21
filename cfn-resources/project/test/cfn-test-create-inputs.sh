#!/usr/bin/env bash
# cfn-test-create-inputs.sh
#
# This tool generates json files in the inputs/ for `cfn test`.
#

set -o errexit
set -o nounset
set -o pipefail

function usage {
    echo "usage:$0 <project_name>"
}

if [ "$#" -ne 1 ]; then usage; fi
if [[ "$*" == help ]]; then usage; fi

rm -rf inputs
mkdir inputs

#create apikey
org_id="$ATLAS_ORG_ID"
team_name="${CFN_TEST_TAG}"

if [ -z "$API_KEY_ID" ]; then
  api_key_id=$(mongocli iam org apikey create --orgId "${ATLAS_ORG_ID}" --desc "${CFN_TEST_TAG}" --role ORG_MEMBER --output json | jq -r '.id')
else
  api_key_id="$API_KEY_ID"
fi

#create team
export user_name=$(mongocli iam project users list --output json | jq -r '.[0].emailAddress')
if [ -z "$TEAM_ID" ]; then
  team_id=$(mongocli iam team create "${team_name}" --username "${user_name}" --orgId "${org_id}" --output json | jq -r '.id')
else
  team_id="$TEAM_ID"
fi


name="${1}"
jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg org "$ATLAS_ORG_ID" \
   --arg name "$name" \
   '.OrgId?|=$org | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .Name?|=$name' \
   "$(dirname "$0")/inputs_1_create.template.json" > "inputs/inputs_1_create.json"

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg org "$ATLAS_ORG_ID" \
   --arg name "${name}- more B@d chars !@(!(@====*** ;;::" \
   '.OrgId?|=$org | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .Name?|=$name' \
   "$(dirname "$0")/inputs_1_invalid.template.json" > "inputs/inputs_1_invalid.json"

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg org "$ATLAS_ORG_ID" \
   --arg name "${name}" \
   --arg key_id "$api_key_id" \
   --arg team_id "$team_id" \
   '.OrgId?|=$org | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey | .Name?|=$name | .ProjectApiKeys[0].Key?|=$key_id | .ProjectTeams[0].TeamId?|=$team_id' \
   "$(dirname "$0")/inputs_1_update.template.json" > "inputs/inputs_1_update.json"

ls -l inputs

echo "TODO: Delete the team and api_key created above"

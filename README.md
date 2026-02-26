# Flaggy

Lightweight feature flag server written in Go. Single binary, SQLite storage, zero external dependencies at runtime.

## Features

- **Flag types** — boolean, string, number, JSON
- **Rule engine** — conditions evaluated with AND logic, priority ordering, first match wins
- **Segments** — reusable groups of conditions shared across rules
- **Rollout** — percentage-based rollout with deterministic bucketing (MurmurHash3)
- **12 operators** — `equals`, `not_equals`, `in`, `not_in`, `contains`, `starts_with`, `gt`, `gte`, `lt`, `lte`, `exists`, `regex`
- **Nested context** — dot-notation attribute resolution (`user.plan`, `user.meta.role`)
- **API key auth** — SHA-256 hashed keys with environment scoping (live/test/staging)
- **SSE streaming** — real-time flag change notifications
- **Batch evaluation** — evaluate multiple flags in a single request
- **CLI** — manage flags, segments, and evaluate from the terminal

## Quick start

### Docker

```bash
docker run -d \
  -p 8080:8080 \
  -e FLAGGY_MASTER_KEY=$(openssl rand -hex 32) \
  -v flaggy-data:/data \
  -e FLAGGY_DB_PATH=/data/flaggy.db \
  ghcr.io/alexis-marcel/flaggy:latest
```

### From source

```bash
git clone https://github.com/Alexis-Marcel/flaggy.git
cd flaggy
make build
FLAGGY_MASTER_KEY=my-secret-key ./bin/flaggyd
```

## Configuration

| Variable | Default | Description |
|---|---|---|
| `FLAGGY_PORT` | `:8080` | Listen address |
| `FLAGGY_DB_PATH` | `flaggy.db` | SQLite database path |
| `FLAGGY_MASTER_KEY` | *(empty)* | Master key for admin routes. If unset, auth is disabled (dev mode) |

## API

All admin routes require `Authorization: Bearer <MASTER_KEY>`. Client routes (evaluate, stream) accept API keys or the master key.

### Flags

```
POST   /api/v1/flags                Create a flag
GET    /api/v1/flags                List all flags
GET    /api/v1/flags/{key}          Get a flag
PUT    /api/v1/flags/{key}          Update a flag
DELETE /api/v1/flags/{key}          Delete a flag
PATCH  /api/v1/flags/{key}/toggle   Toggle enabled/disabled
```

### Rules

```
POST   /api/v1/flags/{key}/rules            Create a rule
PUT    /api/v1/flags/{key}/rules/{ruleID}    Update a rule
DELETE /api/v1/flags/{key}/rules/{ruleID}    Delete a rule
```

### Segments

```
POST   /api/v1/segments             Create a segment
GET    /api/v1/segments             List all segments
GET    /api/v1/segments/{key}       Get a segment
PUT    /api/v1/segments/{key}       Update a segment
DELETE /api/v1/segments/{key}       Delete a segment
```

### Evaluation

```
POST   /api/v1/evaluate             Evaluate a single flag
POST   /api/v1/evaluate/batch       Evaluate multiple flags
GET    /api/v1/stream               SSE stream of flag changes
```

## Usage examples

```bash
export FLAGGY=http://localhost:8080
export AUTH="Authorization: Bearer my-secret-key"

# Create a flag
curl -s -H "$AUTH" $FLAGGY/api/v1/flags -d '{
  "key": "new_checkout",
  "type": "boolean",
  "default_value": false,
  "enabled": true
}'

# Create a segment for pro users
curl -s -H "$AUTH" $FLAGGY/api/v1/segments -d '{
  "key": "pro_users",
  "description": "Users on the pro plan",
  "conditions": [{"attribute": "user.plan", "operator": "equals", "value": "\"pro\""}]
}'

# Create a rule that uses the segment
curl -s -H "$AUTH" $FLAGGY/api/v1/flags/new_checkout/rules -d '{
  "description": "Enable for pro users",
  "segment_keys": ["pro_users"],
  "conditions": [],
  "value": true,
  "priority": 1
}'

# Evaluate
curl -s -H "$AUTH" $FLAGGY/api/v1/evaluate -d '{
  "flag_key": "new_checkout",
  "context": {"user": {"plan": "pro"}}
}'
# → {"flag_key":"new_checkout","value":true,"match":true,"reason":"rule_match"}
```

## CLI

```bash
# Connect to a server
export FLAGGY_SERVER=http://localhost:8080

flaggy flag list
flaggy flag create my_flag --type boolean --default false --enabled
flaggy flag enable my_flag
flaggy flag disable my_flag

flaggy segment list
flaggy segment create pro_users --description "Pro plan users" \
  --conditions '[{"attribute":"user.plan","operator":"equals","value":"\"pro\""}]'
flaggy segment get pro_users

flaggy evaluate my_flag -c '{"user":{"plan":"pro"}}'
```

## How evaluation works

1. If the flag is **disabled** → return default value
2. Sort rules by **priority** (lower number = higher priority)
3. For each rule, evaluate **all inline conditions AND all segment conditions** (AND logic)
4. First rule where everything matches → return the rule's value
5. If a rule has a **rollout percentage**, hash `flagKey:entityID` to check if the user is in the bucket
6. No rule matched → return default value

Segments referenced by a rule that don't exist are treated as **non-matching** (fail closed).

## License

MIT

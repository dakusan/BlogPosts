#!/usr/bin/env bash

RPC="http://127.0.0.1:18081/json_rpc"
POLL=5
STALL=30
BAN_SECONDS=600
RATE_INTERVAL=30

last_height=0
last_progress=$(date +%s)
samples=()

while true; do
	json=$(
		curl -fsS "$RPC" \
			-H 'Content-Type: application/json' \
			-d '{"jsonrpc":"2.0","id":"0","method":"sync_info"}'
	) || {
		echo "Failed to query monerod"
		sleep "$POLL"
		continue
	}

	height=$(jq -r '.result.height // 0' <<<"$json")
	target=$(jq -r '.result.target_height // 0' <<<"$json")
	now=$(date +%s)

	if (( target == 0 || height >= target )); then
		echo "$(date +%T) Height: $height (100.000%) - Synchronized"
		exit 0
	fi

	samples+=("$now:$height")

	# Remove samples older than 30 seconds, retaining one baseline sample
	while (( ${#samples[@]} > 1 )); do
		IFS=: read -r sample_time sample_height <<<"${samples[1]}"
		(( sample_time <= now - RATE_INTERVAL )) || break
		samples=("${samples[@]:1}")
	done

	if (( height != last_height )); then
		percent=$(awk -v h="$height" -v t="$target" 'BEGIN { printf "%.3f", h * 100 / t }')

		remaining=""
		IFS=: read -r sample_time sample_height <<<"${samples[0]}"
		elapsed=$((now - sample_time))
		processed=$((height - sample_height))

		if (( elapsed > 0 && processed > 0 && target > height )); then
			seconds_left=$(( (target - height) * elapsed / processed ))
			printf -v remaining ' %d'\''%02d" eta' "$((seconds_left / 60))" "$((seconds_left % 60))"
		fi

		printf '%s Height: %s (%s%%)%s\n' "$(date +%T)" "$height" "$percent" "$remaining"

		last_height=$height
		last_progress=$now
		sleep "$POLL"
		continue
	fi

	stalled=$((now - last_progress))

	if (( stalled < STALL )); then
		sleep "$POLL"
		continue
	fi

	address=$(
		jq -r --argjson height "$height" '
			(.result.spans // [])
			| map(select(.start_block_height == $height))
			| first
			| .remote_address // empty
		' <<<"$json"
	)

	if [[ -z "$address" ]]; then
		echo "Height $height stalled ${stalled}s, but no head span found"
		sleep "$POLL"
		continue
	fi

	ip="${address%:*}"

	echo "Height $height stalled ${stalled}s: banning $ip for ${BAN_SECONDS}s"

	curl -fsS "$RPC" \
		-H 'Content-Type: application/json' \
		-d "$(jq -nc \
			--arg host "$ip" \
			--argjson seconds "$BAN_SECONDS" \
			'{
				jsonrpc:"2.0",
				id:"0",
				method:"set_bans",
				params:{
					bans:[{
						host:$host,
						ban:true,
						seconds:$seconds
					}]
				}
			}')" >/dev/null

	last_progress=$(date +%s)

	sleep "$POLL"
done
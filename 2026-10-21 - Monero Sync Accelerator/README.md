Original Post: [Monero Sync Accelerator](https://www.castledragmire.com/Posts/Monero_Sync_Accelerator)<br>
Date: 2026-10-21<br>
Tags: Monero, Bash, Optimization<br><br>

# Monero Sync Accelerator

A Bash script that monitors Monero blockchain synchronization, detects stalled downloads, automatically bans unresponsive peers, and provides real-time synchronization progress and estimated completion times.

## Overview

Monero's daemon (`monerod`) downloads blockchain data from multiple peers in parallel. These downloads are organized into *block spans*, each representing a sequence of consecutive blocks retrieved from a particular peer.

Occasionally, a peer may become unresponsive or excessively slow while supplying blocks needed at the current blockchain height. Even when thousands of subsequent blocks have already been downloaded from other peers, synchronization can be delayed while waiting for the missing blocks.

This script addresses that problem by monitoring synchronization progress and automatically banning peers associated with stalled downloads, encouraging `monerod` to retrieve the required blocks from alternative peers.

It also displays blockchain synchronization progress, including an estimated time remaining based on recent processing speed.

## Intended Use

This script is particularly useful when synchronizing a Monero node that has fallen significantly behind the network, where large numbers of blocks must be downloaded and processed.

Its primary purpose is to prevent potentially stalled peer downloads from unnecessarily delaying synchronization while providing a convenient, continuously updated indication of synchronization progress.

It can be run alongside an existing `monerod` instance without restarting or modifying the daemon's startup configuration.

## Features

- **Automatic stall detection:** Detects when the local blockchain height hasn't advanced for a configurable period.
- **Problematic peer identification:** Examines Monero's synchronization spans to identify the peer associated with the current blockchain height.
- **Automatic temporary banning:** Uses Monero's RPC interface to ban the identified peer, allowing synchronization to continue through other peers.
- **Real-time progress reporting:** Displays the current blockchain height, completion percentage, and estimated time remaining.
- **Rolling speed estimation:** Calculates the estimated completion time using the number of blocks processed over approximately the previous 30 seconds.
- **RPC failure handling:** Continues monitoring after temporary communication failures with the daemon.

## Example Output

```text
13:14:15 Height: 3771537 (99.535%) 5'30" eta
13:14:20 Height: 3771557 (99.536%) 5'12" eta
13:14:25 Height: 3771577 (99.536%) 4'58" eta
13:14:55 Height 3771577 stalled 30s: banning 192.0.2.10 for 600s
13:15:00 Height: 3771597 (99.537%) 5'05" eta
```

Each progress line includes:

- Current local time (`HH:MM:SS`).
- Currently processed blockchain height.
- Synchronization percentage to three decimal places.
- Estimated time remaining, expressed in minutes and seconds.

The estimated completion time is omitted until sufficient progress data is available.

## Configuration

The following variables control the script's behavior:

| Variable | Default | Description |
|---|---|---|
| `RPC` | `http://127.0.0.1:18081/json_rpc` | Monero daemon JSON-RPC endpoint |
| `POLL` | `5` | Seconds between synchronization checks |
| `STALL` | `30` | Seconds without blockchain height progress before investigating a stall |
| `BAN_SECONDS` | `600` | Duration of temporary peer bans, in seconds |
| `RATE_INTERVAL` | `30` | Approximate rolling window, in seconds, used for ETA calculations |

## Requirements

- Bash
- `curl`
- `jq`
- A running Monero daemon (`monerod`)
- Access to the daemon's JSON-RPC interface
- RPC permissions sufficient to invoke `sync_info` and `set_bans`

The script uses standard Monero RPC functionality and does not modify the blockchain database directly.

## Limitations

**Height-based detection:** A stall is only detected when the local blockchain height stops advancing. A consistently slow peer that still allows occasional progress may never trigger a ban.

**Span matching:** The script only considers spans beginning exactly at the current height. Other situations causing slow synchronization may not produce a matching span.

**Downloaded versus unresponsive spans:** A matching span does not necessarily mean the peer is currently unresponsive. The span may already be downloaded and awaiting processing. Consequently, the script can ban a peer even when the underlying delay is caused by block verification or database performance.

**Estimated completion time:** The ETA assumes the recent processing rate will continue. Changes in network performance, block complexity, database activity, or synchronization behavior may produce substantial fluctuations.

**Processing bottlenecks:** The script cannot accelerate CPU-bound blockchain verification or database operations. It is intended to mitigate peer-related stalls, not optimize Monero's internal processing.

**Temporary bans:** Peers may be banned unnecessarily when synchronization pauses for reasons unrelated to their performance. Bans expire automatically after the configured duration.

## Performance Recommendations

Monero uses LMDB for blockchain storage, which performs frequent random reads and writes. Storage latency and filesystem behavior can significantly affect synchronization performance.

- **Avoid copy-on-write (CoW) filesystems:** Filesystems such as Btrfs and ZFS can introduce substantial write amplification and fragmentation when used with LMDB. Prefer a filesystem such as ext4 or XFS, or disable CoW specifically for the blockchain database.
- **Use an SSD rather than an HDD:** Although Monero may transfer relatively little data per second, its random I/O workload benefits significantly from the lower latency and higher IOPS of SSDs. NVMe may provide additional improvements, although sequential throughput is generally not the limiting factor.

In one real-world test, moving the blockchain database from Btrfs with CoW enabled to ext4 reduced the estimated synchronization time from approximately 4.5 hours to 1 hour, a roughly **4.5× performance improvement**.

- **Recommended `monerod` startup flags for speed:** The following configuration is intended to improve synchronization performance, increase peer availability, and reduce blockchain storage requirements.

  ```bash
  --prune-blockchain --sync-pruned-blocks --fast-block-sync 1 --db-sync-mode fast --out-peers 32
  ```

  | Flag | Description |
  |---|---|
  | `--prune-blockchain` | Reduces blockchain storage requirements by approximately two-thirds, discarding historical data unnecessary for normal node operation. |
  | `--sync-pruned-blocks` | Allows downloading already-pruned blocks from peers, reducing network traffic during synchronization. |
  | `--fast-block-sync 1` | Uses embedded, known block hashes to accelerate historical blockchain synchronization. Enabled by default. |
  | `--db-sync-mode fast` | Uses faster database synchronization settings, reducing disk synchronization overhead at the expense of some crash resilience. |
  | `--out-peers 32` | Allows up to 32 outgoing peer connections, potentially increasing the availability of block download sources. |

**Database synchronization safety:** Avoid `--db-sync-mode fast` on systems prone to power outages, crashes, or unexpected shutdowns. It improves performance but can cause blockchain database corruption if the operating system crashes or loses power. Use `--db-sync-mode safe:sync` for maximum reliability, at the expense of synchronization speed. Alternatively, omit `--db-sync-mode` entirely to use Monero's default behavior, which prioritizes performance during synchronization and automatically switches to safer database writes near the blockchain tip.

**Note:** These settings are not universally optimal. Peer counts should be adjusted for the available hardware and network connection.

See the [official Monero daemon reference](https://docs.getmonero.org/interacting/monerod-reference/) for additional information.

## Disclosure

This script was largely written with the assistance of an LLM, but required considerable testing, tweaking, and refinement to get it working as intended.

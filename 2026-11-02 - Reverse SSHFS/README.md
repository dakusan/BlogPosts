# Reverse SSHFS

Mount a local directory on a remote Linux machine over SSH, without requiring an SSH connection to the local machine.

Uses `dpipe`, `bubblewrap`, and `sshfs` to securely expose a local directory through an isolated SFTP server, allowing the remote machine to access its contents as a filesystem.

## Usage

```bash
rsshfs <local_path> <user@host:/remote/path> [sshfs_option...]
```

The remote mount directory is created automatically if it does not exist.

Any trailing arguments are passed through to the remote `sshfs` command as-is (e.g. `-o uid=1000 -o gid=1000`).

If no `sshfs_option`s are given, `-o allow_other` is assumed if the remote host has it enabled (`user_allow_other` in `/etc/fuse.conf`).

Press **Ctrl+C** to unmount and terminate the connection.

**Example:**

```bash
rsshfs ~/projects user@example.com:/home/user/mnt/projects
```

This mounts the local `~/projects` directory at `/home/user/mnt/projects` on `example.com`.

## Features

- **Reverse SSHFS:** Mounts a local directory on a remote machine using an outbound SSH connection.
- **No local SSH server required:** Uses OpenSSH's `sftp-server` directly.
- **Sandboxed access:** Runs the SFTP server inside a Bubblewrap sandbox, exposing only the specified local directory.
- **Exits on unmount:** Terminates automatically once the remote mount is unmounted, whether via Ctrl+C or externally.
- **Automatic cleanup:** Attempts to unmount the remote filesystem when interrupted with Ctrl+C or SIGTERM.
- **Read/write access:** Files can be read and modified from the remote machine, subject to local filesystem permissions.

## Requirements

**Local machine:**
- Linux
- Bash
- OpenSSH client and `sftp-server`
- `dpipe` (from the `vde2` package)
- `bubblewrap` (`bwrap`)
- `realpath`

**Remote machine:**
- Linux
- OpenSSH server
- `sshfs`
- `fusermount3`

SSH authentication to the remote machine must already be configured.

## How It Works

1. Resolves and validates the local directory.
2. Creates a Bubblewrap sandbox with the local directory mounted at `/share`.
3. Runs OpenSSH's `sftp-server` inside the sandbox.
4. Uses `dpipe` to connect the SFTP server's standard input/output to a remote SSH session.
5. Runs `sshfs` on the remote machine in passive mode, using the SSH connection as its SFTP transport.
6. Keeps the connection active until interrupted, then attempts to unmount the remote filesystem.

## Security

The SFTP server runs inside a separate mount and namespace environment using `bwrap --unshare-all`. The local directory is exposed at `/share`, while system binaries and essential configuration files are mounted read-only.

The remote machine does not receive general access to the local filesystem or a shell on the local machine.

**Important:** The shared directory is writable. The remote host can modify or delete files accessible through the mount. The `allow_other` SSHFS option also permits other remote users to access the mount, subject to filesystem permissions and FUSE configuration.

This is a filesystem isolation mechanism, not a guarantee against every possible sandbox escape or remote-host attack.

## Troubleshooting

**<ins>`allow_other` permission denied</ins>**

Enable `user_allow_other` in `/etc/fuse.conf` on the remote machine.

**<ins>Remote mount is unresponsive</ins>**

Verify that the SSH connection is still active. An interrupted transport may require manually unmounting the remote filesystem:

```bash
fusermount3 -u /remote/path
```

**<ins>Mount remains busy</ins>**

Close applications accessing the mounted directory, then retry the unmount.

## Blog post mirror

[Reverse SSHFS](https://www.castledragmire.com/Posts/Reverse_SSHFS)<br>
Date: 2026-11-02<br>
Tags: Command Line Utility, Servers, SSH, SSH Tunnels, SSHFS

SSH is still surprising me after 25+ years of use.
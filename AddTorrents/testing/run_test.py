#!/usr/bin/env python3
"""
End-to-end test for AddTorrent.

Builds AddTorrent and the fake qBittorrent server (fakeqbit.go), runs them both against a temporary
workspace, drops several fake .torrent files (each containing a random 8-character token as its
content), and simulates a browser by polling /Updates throughout the run (AddTorrent only logs
"Finished downloading"/cleans up a tag when something actually polls /Updates - there's no background
timer for that on its own).

At the end it checks:
  - The expected output files exist in the right places with the right content (original .torrent gone
    from root, present in doneDir, and a resulting .mp4 present under the mirrored newRoot path
    containing the same token).
  - The captured stdout/stderr log contains every expected line, IN ORDER (extra/noisy lines - like the
    fsnotify "Event: ..." lines - are allowed in between; this is a forward, in-order substring search,
    not an exact match).

Prints PASS/FAIL and exits 0/1 accordingly. The web UI stays up for the whole run, so you can open
http://127.0.0.1:<ADDTORRENT_PORT>/ in a browser while it's running to watch it live.
"""

import base64
import json
import os
import re
import shutil
import subprocess
import sys
import threading
import time
import urllib.error
import urllib.request

ScriptDir = os.path.dirname(os.path.abspath(__file__))
RepoRoot = os.path.dirname(ScriptDir)

QbitPort = 18080
AddTorrentPort = 18081
MaxLogLines = 20
DownloadSeconds = 20
Stagger = 5  #seconds between each fake torrent being dropped
PollInterval = 1  #how often we (as a stand-in browser) poll /Updates
DebounceMargin = 3 + 2  #ProcessFile's own 3s "let the file finish writing" sleep, plus safety margin
FinishMargin = 2 * PollInterval + 5  #extra time to wait after a fake download completes, so our poller has a chance to notice it

#(Category, Show, Filename) for each fake torrent, in drop order
Torrents = [
	("Cat1", "ShowA", "T1.torrent"),
	("Cat1", "ShowA", "T2.torrent"),
	("Cat1", "ShowA", "T3.torrent"),
	("Cat2", "ShowB", "T4.torrent"),
	("Cat2", "ShowB", "T5.torrent"),
]


def RandomToken():
	return base64.urlsafe_b64encode(os.urandom(6)).decode()


def RunCmd(Args, Cwd=None):
	Result = subprocess.run(Args, cwd=Cwd, capture_output=True, text=True)
	if Result.returncode != 0:
		print(f"Command failed: {' '.join(Args)}\n--- stdout ---\n{Result.stdout}\n--- stderr ---\n{Result.stderr}")
		sys.exit(1)


def PostUpdates(LastLogLine, ServerStartTime):
	Req = urllib.request.Request(
		f"http://127.0.0.1:{AddTorrentPort}/Updates",
		data=json.dumps({"LastLogLine": LastLogLine, "ServerStartTime": ServerStartTime}).encode(),
		headers={"Content-Type": "application/json"},
		method="POST",
	)
	with urllib.request.urlopen(Req, timeout=5) as Resp:
		return json.loads(Resp.read())


def PollerThread(StopEvent):
	LastLogLine = -1
	ServerStartTime = 0
	while not StopEvent.is_set():
		try:
			Data = PostUpdates(LastLogLine, ServerStartTime)
			if "Error" not in Data:
				ServerStartTime = Data.get("ServerStartTime", ServerStartTime)
				LastLogLine = Data.get("NewLastLogLine", LastLogLine)
		except Exception:
			pass  #server may not be up yet, or a transient hiccup; just retry next tick
		StopEvent.wait(PollInterval)


def WaitForUrl(Url, Timeout=15):
	Deadline = time.time() + Timeout
	while time.time() < Deadline:
		try:
			with urllib.request.urlopen(Url, timeout=2):
				return True
		except Exception:
			time.sleep(0.5)
	return False


def CheckLogOrder(LogPath, Expected):
	"""Expected is a list of (Description, CompiledRegex). Scans LogPath forward, requiring each entry's
	regex to be found at or after the line where the previous entry matched. Returns (True, "") or
	(False, FailureMessage)."""
	with open(LogPath, "r", errors="replace") as F:
		LogLines = F.readlines()

	Cursor = 0
	for Description, Pattern in Expected:
		Found = False
		for I in range(Cursor, len(LogLines)):
			if Pattern.search(LogLines[I]):
				Cursor = I + 1
				Found = True
				break
		if not Found:
			Context = "".join(LogLines[max(0, Cursor - 3):Cursor + 5])
			return False, f"Expected log line not found (in order) for step {Description!r} (pattern: {Pattern.pattern!r}) after line {Cursor}.\nNearby log context:\n{Context}"
	return True, ""


def Main():
	print("Building AddTorrent...")
	RunCmd(["go", "build", "-o", os.path.join(RepoRoot, "AddTorrent"), os.path.join(RepoRoot, "AddTorrent.go")])
	print("Building fakeqbit...")
	RunCmd(["go", "build", "-o", os.path.join(ScriptDir, "fakeqbit"), os.path.join(ScriptDir, "fakeqbit.go")])

	Workspace = os.path.join(ScriptDir, "tmp_workspace")
	if os.path.exists(Workspace):
		shutil.rmtree(Workspace)
	Root = os.path.join(Workspace, "root") + "/"
	NewRoot = os.path.join(Workspace, "newRoot") + "/"
	DoneDir = os.path.join(Workspace, "done") + "/"
	for Category, Show, _ in Torrents:
		os.makedirs(os.path.join(Root, Category, Show), exist_ok=True)

	AddTorrentLogPath = os.path.join(Workspace, "addtorrent.log")
	FakeQbitLogPath = os.path.join(Workspace, "fakeqbit.log")
	ExpectedLogPath = os.path.join(Workspace, "expected_log.json")

	AddTorrentLog = open(AddTorrentLogPath, "w")
	FakeQbitLog = open(FakeQbitLogPath, "w")

	FakeQbitProc = subprocess.Popen(
		[os.path.join(ScriptDir, "fakeqbit"), "-listenOn", f":{QbitPort}", "-downloadSeconds", str(DownloadSeconds)],
		stdout=FakeQbitLog, stderr=subprocess.STDOUT,
	)

	StopEvent = threading.Event()
	Poller = threading.Thread(target=PollerThread, args=(StopEvent,), daemon=True)
	AddTorrentProc = None

	Passed = False
	FailureReason = ""
	try:
		#fakeqbit must already be accepting connections before AddTorrent's first (synchronous,
		#startup-time) login attempt, or that attempt silently fails with no retry for 15 seconds
		#(qbitMonitor only retries while disconnected; the very first failure logs nothing, since
		#"Lost connection" is only logged on a connected->disconnected transition).
		if not WaitForUrl(f"http://127.0.0.1:{QbitPort}/api/v2/auth/login"):
			raise RuntimeError("fakeqbit never came up")

		AddTorrentProc = subprocess.Popen(
			[
				os.path.join(RepoRoot, "AddTorrent"),
				"-config", os.path.join(RepoRoot, "config.default.json"),
				"-root", Root,
				"-newRoot", NewRoot,
				"-doneDir", DoneDir,
				"-listenOn", f":{AddTorrentPort}",
				"-qbitURL", f"http://127.0.0.1:{QbitPort}",
				"-qbitUser", "testuser",
				"-qbitPass", "testpass",
				"-fileRegex", r"^[^/]+/.+/[^/]+$",
				"-maxLogLines", str(MaxLogLines),
			],
			stdout=AddTorrentLog, stderr=subprocess.STDOUT, cwd=RepoRoot,
		)

		if not WaitForUrl(f"http://127.0.0.1:{AddTorrentPort}/"):
			raise RuntimeError("AddTorrent's web server never came up")
		print(f"AddTorrent is up. Open http://127.0.0.1:{AddTorrentPort}/ in a browser now to watch it live.")
		Poller.start()

		#Build the ordered list of expected log lines, and drop the fake torrents on schedule
		ExpectedLog = [
			("qBittorrent connected", re.compile(r"Connected to qbittorrent WebUI")),
			("Started watching root", re.compile(r"Started watching:")),
		]
		Tokens = {}
		StartTime = time.time()
		for I, (Category, Show, Filename) in enumerate(Torrents):
			TargetTime = StartTime + I * Stagger
			SleepFor = TargetTime - time.time()
			if SleepFor > 0:
				time.sleep(SleepFor)

			Token = RandomToken()
			Tokens[(Category, Show, Filename)] = Token
			Path = os.path.join(Root, Category, Show, Filename)
			with open(Path, "w") as F:
				F.write(Token)
			print(f"Dropped {Category}/{Show}/{Filename} (token={Token})")

			ExpectedLog.append((f"{Filename} copied", re.compile(rf"File copied to:.*{re.escape(Filename)}")))
			ExpectedLog.append((f"{Filename} added to qbittorrent", re.compile(rf"Added to qbittorrent:.*{re.escape(Filename)}")))

		#Finish order matches add order here since Stagger is comfortably larger than PollInterval
		for Category, Show, Filename in Torrents:
			ExpectedLog.append((f"{Filename} finished downloading", re.compile(rf"Finished downloading:.*{re.escape(Filename)}")))

		LastDropTime = StartTime + (len(Torrents) - 1) * Stagger
		LastFinishEstimate = LastDropTime + DebounceMargin + DownloadSeconds + FinishMargin
		LastOutputFile = os.path.join(NewRoot, Torrents[-1][0], Torrents[-1][1], Torrents[-1][2].replace(".torrent", ".mp4"))

		print(f"Waiting for the last fake download to finish (up to {LastFinishEstimate - time.time():.0f}s)...")
		Deadline = LastFinishEstimate + 30  #extra safety margin on top of our own estimate
		while time.time() < Deadline and not os.path.exists(LastOutputFile):
			time.sleep(1)
		#Give the poller a final chance to notice the last completion before we stop it
		time.sleep(FinishMargin)

		#One plain page load, to exercise/verify that path too
		with urllib.request.urlopen(f"http://127.0.0.1:{AddTorrentPort}/", timeout=5):
			pass
		ExpectedLog.append(("final page load logged", re.compile(r"Log request:")))
		time.sleep(0.5)

		StopEvent.set()
		Poller.join(timeout=5)

		with open(ExpectedLogPath, "w") as F:
			json.dump([{"desc": D, "regex": P.pattern} for D, P in ExpectedLog], F, indent=2)

		#Filesystem checks
		Failures = []
		for Category, Show, Filename in Torrents:
			SourcePath = os.path.join(Root, Category, Show, Filename)
			DonePath = os.path.join(DoneDir, Filename)
			OutPath = os.path.join(NewRoot, Category, Show, Filename.replace(".torrent", ".mp4"))
			if os.path.exists(SourcePath):
				Failures.append(f"{Filename}: original file still present at {SourcePath}")
			if not os.path.exists(DonePath):
				Failures.append(f"{Filename}: not found in doneDir at {DonePath}")
			if not os.path.exists(OutPath):
				Failures.append(f"{Filename}: fake downloaded output not found at {OutPath}")
			else:
				with open(OutPath) as F:
					Content = F.read()
				ExpectedToken = Tokens[(Category, Show, Filename)]
				if Content != ExpectedToken:
					Failures.append(f"{Filename}: output content {Content!r} != expected token {ExpectedToken!r}")

		#Negative checks: nothing suggesting a real problem should appear anywhere in the log.
		#Note: "Error statting: ... no such file or directory" is NOT included here even though it
		#contains "Error" - it's an expected, benign race: AddTorrent's own rename-into-doneDir fires a
		#RENAME fsnotify event on the (now already-moved) source path, which legitimately fails to stat.
		with open(AddTorrentLogPath, errors="replace") as F:
			FullLog = F.read()
		for BadPattern in ("Gave up tracking", "Lost connection to qbittorrent", "Warning: Could not", "Error adding torrent", "Error evaluating fileRegex", "Could not finish copying", "Refusing to write outside"):
			if BadPattern in FullLog:
				Failures.append(f"Unexpected text found in log: {BadPattern!r}")

		OrderOk, OrderMessage = CheckLogOrder(AddTorrentLogPath, ExpectedLog)
		if not OrderOk:
			Failures.append(OrderMessage)

		if Failures:
			FailureReason = "\n".join(Failures)
		else:
			Passed = True
	except Exception as E:
		FailureReason = str(E)
	finally:
		StopEvent.set()
		Procs = [P for P in (AddTorrentProc, FakeQbitProc) if P is not None]
		for Proc in Procs:
			Proc.terminate()
		for Proc in Procs:
			try:
				Proc.wait(timeout=5)
			except subprocess.TimeoutExpired:
				Proc.kill()
		AddTorrentLog.close()
		FakeQbitLog.close()

	print()
	if Passed:
		print("PASS")
	else:
		print("FAIL")
		print(FailureReason)
	print(f"\nWorkspace (logs, expected_log.json, output files): {Workspace}")
	sys.exit(0 if Passed else 1)


if __name__ == "__main__":
	Main()

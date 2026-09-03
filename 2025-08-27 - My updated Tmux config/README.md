<div>
Original Post: <a href="https://www.castledragmire.com/Posts/My_updated_Tmux_config">My updated Tmux config</a><br>
Date: 2025-08-27<br>
Tags: Bash, Configuration, Linux, screen, Tmux<br><br>
</div>

<p>Tmux is a great alternative to gnu screen. I’ve <a href="../2020-01-10 - My Tmux config">previously posted my cygwin tmux config</a> in 2020, and I figured I’d update with my additions since then, including running under real linux. So here it is. Its features include:</p>
<ul>
	<li>Uses <b>ctrl+a</b>, like gnu screen, instead of <b>ctrl+b</b></li>
	<li>Mouse interaction is enabled</li>
	<li>Tab bar/windows:<ul>
		<li>Current tab is highlighted in cyan</li>
		<li>Cycle through tabbed windows with a <b>click</b> on its tab or <b>ctrl+shift+alt+arrow_keys</b></li>
		<li>Reorder tabbed windows with a <b>drag</b> of its tab <strike>or <b>alt+arrow_keys</b></strike></li>
		<li><b>ctrl+a,/</b> to rename a tab on the tab bar</li>
		<li>Create new window with <b>ctrl+a,c</b></li>
	</ul></li>
	<li>Panes<ul>
		<li>Create split panes with vertical=<b>ctrl+a,|</b> and horizontal=<b>ctrl+a,-</b> (and start in the cwd)</li>
		<li>Move around panes with <b>click</b> or <b>ctrl+shift+arrow_keys</b></li>
		<li>Resize panes by <b>dragging</b> on the separator bar or use <b>alt+shift+arrow_keys</b></li>
		<li>Swap panes directionally with <b>alt+arrow_keys</b></li>
		<li>Panes automatically resize to fit OS window</li>
	</ul></li>
	<li>Clipboard/highlighting<ul>
		<li>Copy text to clipboard by <b>highlighting</b> it</li>
		<li>Paste from clipboard with <b>right click</b></li>
		<li><b>Middle mouse button+drag</b> highlights and copies too. I did this since nano takes over normal mouse click+drag.</li>
	</ul></li>
	<li>Copy mode<br>
		<table>
			<tr><th>Key*</th><th>Action</th></tr>
			<tr><td>Ctrl+PageUp</td><td>Start copy mode, scroll page up</td></tr>
			<tr><td>Ctrl+PageDown</td><td>Start copy mode</td></tr>
			<tr><td><code>Enter</code></td><td>Copies selected text to clipboard and ends copy mode</td></tr>
			<tr><td><code>Shift+(Left/Right)</code></td><td>Starts selection mode on previous/next character (or continues selection navigation like normal left/right)</td></tr>
			<tr><td><code>Escape</code></td><td>If something is selected, end current selection. Otherwise, exit copy mode</td></tr>
			<tr><td><code>Ctrl+(Home/End)</code></td><td>Go to top/bottom of the buffer</td></tr>
			<tr><td><code>Ctrl+(Left/Right)</code></td><td>Moves to previous/next word like when not in copy mode</td></tr>
			<tr><td><code>Ctrl+F3</code>, Ctrl+F3</td><td>Start incremental search in copy mode</td></tr>
			<tr><td><code>F3/Shift+F3</code></td><td>Search next/previous (Note, if you use gnome-terminal it may override this)</td></tr>
		</table>
		* When <code>in highlighted block</code> the key is for when in copy mode
	</li>
	<li>Start the session on the current bash directory</li>
	<li>Escape time is lowered for quicker response to copy mode access</li>
</ul><br>
<p>To use this, save the file to <b>~/.tmux.conf</b></p>

```bash
#Mouse interaction
set -g mouse on

#Lower escape timing from 500ms to 50ms for quicker response to copy-mode access
set -s escape-time 50

#Window always takes up largest possible max size
set-window-option -g aggressive-resize

#Highlight active window in tab-bar at bottom in cyan
set-window-option -g window-status-current-style bg=cyan

#Reorder windows in status bar by drag & drop
bind-key -n MouseDrag1Status swap-window -t=

#Copy to clipboard on text selection
bind -Tcopy-mode MouseDragEnd1Pane { send -X copy-selection-and-cancel; run-shell -b "tmux show-buffer | xclip -selection clipboard" }

#Paste from clipboard with right click
bind-key -n MouseDown3Pane { run-shell 'tmux set-buffer -b winclip "$(xclip -o -selection clipboard)"'; paste-buffer -db winclip }

#Middle drag runs copy (Since some programs like nano take control of the mouse for normal selection)
bind -n MouseDrag2Pane copy-mode -M
bind -Tcopy-mode MouseDragEnd2Pane { send -X copy-selection-and-cancel; run-shell -b "tmux show-buffer | xclip -selection clipboard" }
bind-key -n MouseDown2Pane run-shell -b "echo '' > /dev/null"

#Remap prefix to Control+a
set -g prefix C-a
unbind C-b
#bind 'C-a C-a' to type 'C-a'
bind C-a send-prefix

#Start in CWD when creating or splitting tabs; move the splitting planes keys to | and -
bind '|' split-window -h -c '#{pane_current_path}'  # Split panes horizontal
bind '\' split-window -h -c '#{pane_current_path}'  # Split panes horizontal
bind '-' split-window -v -c '#{pane_current_path}'  # Split panes vertically
bind c new-window -c '#{pane_current_path}' # Create new window
unbind '"'
unbind %

#Select next/prev window with Ctrl+Shift+Alt+(Left|Right)
bind-key -n C-S-M-Right next-window
bind-key -n C-S-M-Left previous-window

#Reorder panes with Alt+arrow
bind-key -n M-Left swap-pane -t '{left-of}' -d
bind-key -n M-Right swap-pane -t '{right-of}' -d
bind-key -n M-Up swap-pane -t '{up-of}' -d
bind-key -n M-Down swap-pane -t '{down-of}' -d

#Switch panes using Ctrl+shift+arrow - I rather would have used Ctrl+arrow, but that can interfeer with other programs like nano
bind -n C-S-Left select-pane -L
bind -n C-S-Right select-pane -R
bind -n C-S-Up select-pane -U
bind -n C-S-Down select-pane -D

#Resize panes using Alt+Shift+arrow
bind-key -n M-S-Up resize-pane -U 1
bind-key -n M-S-Down resize-pane -D 1
bind-key -n M-S-Left resize-pane -L 1
bind-key -n M-S-Right resize-pane -R 1

#Ctrl+F3 starts searching in copy mode
bind -n C-F3 { copy-mode; command-prompt -T search -ip'(search down)' -I'#{pane_search_string}' { send -X search-forward-incremental '%%' } }

#prefix, / -- Renames window, but starts blank
bind-key / command-prompt "rename-window '%%'"

#---Copy mode---

#Enter key in copy mode copies to clipboard and ends copy mode
bind-key -Tcopy-mode Enter { send -X copy-selection-and-cancel; run-shell -b "tmux show-buffer | xclip -selection clipboard" }
#Ctrl+PageUp = Copy mode, scroll page up
bind-key -n C-PageUp copy-mode -u
#Ctrl+PageDown = Copy mode
bind-key -n C-PageDown copy-mode 
#Shift+Left = Start selection mode on previous character (if selection mode is active, do not restart selection)
bind-key -Tcopy-mode S-Left { run-shell "tmux display-message -p \"#{selection_active}\" | awk '{print ($0 == 1) ? \"\" : \"tmux send -X begin-selection\"}' | sh"; send -X cursor-left }
#Shift+Right = Start selection mode on current character (if selection mode is active, do not restart selection)
bind-key -Tcopy-mode S-Right { run-shell "tmux display-message -p \"#{selection_active}\" | awk '{print ($0 == 1) ? \"\" : \"tmux send -X begin-selection\"}' | sh"; send -X cursor-right }
#Escape = If something is selected, end current selection. Otherwise, exit copy mode.
bind-key -Tcopy-mode Escape run-shell "tmux display-message -p \"#{selection_active}\" | awk '{print ($0 == 1) ? \"tmux send -X clear-selection\" : \"tmux send -X cancel\"}' | sh"
#Ctrl+Home = Go to top of buffer
bind-key -Tcopy-mode C-Home send -X history-top
#Ctrl+End = Go to bottom of buffer
bind-key -Tcopy-mode C-End send -X history-bottom
#Ctrl+Left = Previous work
bind-key -Tcopy-mode C-Left send -X previous-word
#Ctrl+Right = Next work
bind-key -Tcopy-mode C-Right send -X next-word
#Ctrl+F3 = Start incremental search
bind -Tcopy-mode C-F3 { command-prompt -T search -ip'(search down)' -I'#{pane_search_string}' { send -X search-forward-incremental '%%' } }
#F3 = Search next (Note, if you use gnome-terminal it may override this)
bind -Tcopy-mode F3 send -X search-again
#Shift+F3 = Search previous (Note, if you use gnome-terminal it may override this)
bind -Tcopy-mode S-F3 send -X search-reverse
```
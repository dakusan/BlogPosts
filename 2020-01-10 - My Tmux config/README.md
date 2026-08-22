<div>
Original Post: <a href="https://www.castledragmire.com/Posts/My_Tmux_config">My Tmux config</a><br>
Date: 2020-01-10<br>
Tags: Bash, Configuration, Cygwin, screen, Tmux<br><br>

<b>Updated config at <a href="2025-08-27 - My updated Tmux config">My updated Tmux config</a></b>

<p>Tmux is a great alternative to gnu screen. I can’t believe I’ve never posted my custom Tmux config for Cygwin after all the work I put into it years ago. So here it is. Its features include:</p>
<ul>
	<li>Uses <b>ctrl+a</b>, like gnu screen, instead of <b>ctrl+b</b></li>
	<li>Mouse interaction is enabled</li>
	<li>Tab bar/windows:<ul>
		<li>Current tab is highlighted in cyan</li>
		<li>Cycle through tabbed windows with a <b>click</b> on its tab or <b>ctrl+arrowkeys</b></li>
		<li>Reorder tabbed windows with a <b>drag</b> of its tab or <b>alt+arrowkeys</b></li>
		<li><b>ctrl+a,/</b> to rename a tab on the tab bar</li>
		<li>Create new window with <b>ctrl+a,c</b></li>
	</ul></li>
	<li>Panes<ul>
		<li>Create split panes with vertical=<b>ctrl+a,|</b> and horizontal=<b>ctrl+a,-</b></li>
		<li>Move around panes with <b>click</b> or <b>ctrl+shift+arrowkeys</b></li>
		<li>Resize panes by <b>dragging</b> on the separator bar or use <b>ctrl+shift+alt+arrowkeys</b></li>
		<li>Panes automatically resize to fit OS window</li>
	</ul></li>
	<li>Clipboard/highlighting<ul>
		<li>Copy text to clipboard by <b>highlighting</b> it. Had to use a minor hack to fix a cygwin selection problem</li>
		<li>Paste from clipboard with <b>right click</b></li>
		<li><b>Middle mouse button+drag</b> starts copy mode<ul>
			<li>When in copy mode, <b>u</b> runs the selection as a command in a separate window (Instead of “cygstart” for cygwin, use “xdg-open” for linux, or “open” for MacOS X)</li>
		</ul></li>
		<li><b>Double click</b> selects word</li>
		<li><b>Double middle click</b> runs the word under the mouse as a command</li>
	</ul></li>
	<li>Start the session on the current bash directory</li>
	<li>Escape time is lowered for quicker response to scroll buffer access (<b>ctrl+a,pageup</b>)</li>
</ul>
<p>To use this, save the file to <b>~/.tmux.conf</b></p>

```bash
#Set current directory setting for cygwin
set-environment -g CHERE_INVOKING 1

#Mouse interaction
set -g mouse on

#Lower escape timing from 500ms to 50ms for quicker response to scroll-buffer access
set -s escape-time 50

#Window always takes up largest possible max size
set-window-option -g aggressive-resize

#Highlight active window in tab-bar at bottom in cyan
set-window-option -g window-status-current-bg cyan

#Reorder windows in status bar by drag & drop
bind-key -n MouseDrag1Status swap-window -t=

#Copy to clipboard on text selection in cygwin. Move cursor position 1 to the right before copy to bypass a bug
bind -Tcopy-mode MouseDragEnd1Pane send-keys -X cursor-right\; send -X copy-selection-and-cancel\; run-shell -b "tmux show-buffer > /dev/clipboard"

#Paste from clipboard with right click in cygwin
bind-key -n MouseDown3Pane run-shell 'tmux set-buffer -b winclip "$(cat /dev/clipboard)"'\; paste-buffer -db winclip

#Middle drag starts copy mode
bind -n MouseDrag2Pane copy-mode -M

#When in copy mode, "u" runs the selection as a command in a separate window (Instead of "cygstart" for cygwin, use "xdg-open" for linux, or "open" for MacOS X)
bind -Tcopy-mode u send -X copy-selection-and-cancel\; run-shell -b "tmux show-buffer | xargs cygstart"

#Double click selects word
bind-key -n DoubleClick1Pane copy-mode -M\; send-keys -X select-word

#Double middle click runs the word under the mouse as a command. See description for MouseDown3Pane above
bind-key -n DoubleClick2Pane copy-mode -M\; send-keys -X select-word\; send -X copy-selection-and-cancel\; run-shell -b "tmux show-buffer | xargs cygstart"

#Remap prefix to Control+a
set -g prefix C-a
unbind C-b
#bind 'C-a C-a' to type 'C-a'
bind C-a send-prefix

#Start in CWD when creating or splitting tabs; move the splitting planes keys to | and -
bind '|' split-window -h -c '#{pane_current_path}'  # Split panes horizontal
bind '-' split-window -v -c '#{pane_current_path}'  # Split panes vertically
bind c new-window -c '#{pane_current_path}' # Create new window
unbind '"'
unbind %

#prefix, / -- Renames window, but starts blank
bind-key / command-prompt "rename-window '%%'"

#Select next/prev window with Ctrl+(Left|Right)
bind-key -n C-Right next-window
bind-key -n C-Left previous-window

#Reorder window with Alt+(Left|Right)
bind-key -n M-Left swap-window -t -1
bind-key -n M-Right swap-window -t +1

#Switch panes using Ctrl+Shift+arrow
bind -n C-S-Left select-pane -L
bind -n C-S-Right select-pane -R
bind -n C-S-Up select-pane -U
bind -n C-S-Down select-pane -D

#Resize panes using Ctrl+Shift+Alt+arrow
bind-key -n C-S-M-Up resize-pane -U 1
bind-key -n C-S-M-Down resize-pane -D 1
bind-key -n C-S-M-Left resize-pane -L 1
bind-key -n C-S-M-Right resize-pane -R 1
```
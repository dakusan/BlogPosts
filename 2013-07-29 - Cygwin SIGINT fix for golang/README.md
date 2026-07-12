<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Cygwin_SIGINT_fix_for_golang">Cygwin SIGINT fix for golang</a><br>
Date: 2013-07-29<br>
Tags: Command Line Utility, Cygwin, Examples, GoLang, Quick Fixes, Servers, SIGINT<br><br>

<p>Cygwin has had a long time problem that, depending on your configuration, may cause you to be unable to send a SIGINT (interrupt signal via Ctrl+C) to a native Windows command line executables. As a matter of fact, trying to do so may completely freeze up the console, requiring a process kill of the actual console, bash, and the executable you ran. This problem can crop up for many reasons including the version of Cygwin you are running and your terminal emulator. I specifically installed mintty as my default Cygwin console to get rid of this problem a long time ago (among many other features it had), and now it even has this problem.</p>
<p>While my normal solution is to try and steer clear of native Windows command line executables in Cygwin, this is not always an option. Golang was also causing me this problem every time I ran a network server, which was especially problematic as I would have to ALSO manually kill the server process or it would continue to hold the network port so another test of the code could not use it. An example piece of code is as follows:</p>
</div>

```go
package main
import ( "net/http"; "fmt" )
func main() {
	var HR HandleRequest
	if err := http.ListenAndServe("127.0.0.1:81", HR); err!=nil {
		fmt.Println("Error starting server") }
}

//Handle a server request
type HandleRequest struct{}
func (HR HandleRequest) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received connection from: %s\n", req.RemoteAddr)
}
```

```bash
go run example.go
```

The first solution I found to this problem, which is by far the best solution, was to build the executable and then run it, instead of just running it straight from go.

```bash
go build example.go && example.exe
```

However, as of this post, it seems to no longer work! The last time I tested it and confirmed it was working was about 3 months ago, so who knows what has changed since then.

The second solution is to just build in some method of killing the process that uses “os.Exit”. For example, the following will exit if the user types “exit”

```go
func ListenForExitCommand() {
	for s:=""; s!="exit"; { //Listen for the user to type exit
		if _, err:=fmt.Scanln(&s); err!=nil {
			if err.Error()!="unexpected newline" {
				fmt.Println(err) }
		} else if s=="flush" || s=="exit" {
			//Clean up everything here
		}
	}
	fmt.Println("Exit received, closing process")
	os.Exit(1)
}
```

and then add the following at the top of the main function:

```bash
go ListenForExitCommand() //Listen for "exit" command
```
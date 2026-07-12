package main
import ( "net/http"; "fmt" )

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

func main() {
	var HR HandleRequest
	go ListenForExitCommand() //Listen for "exit" command
	if err := http.ListenAndServe("127.0.0.1:81", HR); err!=nil {
		fmt.Println("Error starting server") }
}

//Handle a server request
type HandleRequest struct{}
func (HR HandleRequest) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received connection from: %s\n", req.RemoteAddr)
}
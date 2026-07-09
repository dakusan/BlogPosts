package main
import ( "./encryptedtcp"; "fmt"; "log"; "os" )

func main() {
	//Confirm parameters, and if invalid, print the help
	if len(os.Args) < 4 || (os.Args[2] != "y" && os.Args[2] != "n") {
		log.Print("First Parameter: ip address to connect to\nSecond Parameter: y = encrypted, n = unencrypted\nAdditional Parameters (at least 1 required): messages to send\n"); return }

	//Initialize the connection
	conn, err := encryptedtcp.StartClient("client.pem", "client.key", os.Args[1]+":16001", os.Args[2]=="y" )
	if err != nil {
		log.Printf("%q\n", err); return }
	defer conn.Close()

	//Process all parameters past the first
	buf := make([]byte, 512)
	for _, msg := range os.Args[3:] {
		//Send the parameter
		if(len(msg)==0) {
			continue }
		n, err := conn.Write([]byte(msg))
		if err != nil {
			log.Printf("Error Writing: %q\n", err); break }
		fmt.Printf("Sent: %q\n", msg[:n])

		//Receive the reply
		n, err = conn.Read(buf)
		if err != nil {
			log.Printf("Error Reading: %q\n", err); break }
		fmt.Printf("Received: %q\n", string(buf[:n]))
	}
}
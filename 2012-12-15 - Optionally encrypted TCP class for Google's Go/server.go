package main
import ( "./encryptedtcp"; "fmt"; "log" )

func main() {
	if err := encryptedtcp.StartServer("server.pem", "server.key", "0.0.0.0:16001", handleClient); err != nil {
		log.Printf("%q\n", err) }
}

func handleClient(conn encryptedtcp.ReadWriteClose) {
	buf := make([]byte, 512)
	for {
		//Read data
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error Reading: %q\n", err); break }
		fmt.Printf("Received: %q\n", string(buf[:n]))

		//Reverse data
		for i, m := 0, n/2; i<m; i++ { //Iterate over half the list
			buf[i], buf[n-i-1] = buf[n-i-1], buf[i] } //Swap first and half of list 1 char at a time

		//Echo back reversed data
		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("Error Writing: %q\n", err); break }
		fmt.Printf("Sent: %q\n", string(buf[:n]))
	}
}
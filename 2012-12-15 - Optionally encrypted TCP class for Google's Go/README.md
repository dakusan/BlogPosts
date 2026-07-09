<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Optionally_encrypted_TCP_class_for_Google%27s_Go">Optionally encrypted TCP class for Google&#039;s Go</a><br>
Date: 2012-12-15<br>
SubTitle: Yet another new language to play with<br>
Tags: Encryption, Examples, Go, GoLang, Networking<br><br>

<p>I wanted to play around with Google's go language a little so I ended up decided on making a simple class that helps create a TCP connection between a server and client that is encrypted via TLS, or not, depending upon a flag. Having the ability to not encrypt a connection is useful for debugging and testing purposes, especially if other people are needing to create clients to connect to your server.</p><br>

<p>The example server.go file listens on port 16001 and for every set of data it receives, it sends the reversed string back to the client. (Note there are limitations to the string lengths in the examples due to buffer and packet payload length restrictions).</p><br>
<p>The example client.go file connects to the server (given via the 1st command line parameter), optionally encrypts the connection (depending upon the 2nd command line parameter), and sends the rest of the parameters to the server as strings.</p><br>
The encryptedtcp.go class has the following exported functions:
<ul>
	<li>StartServer: Goes into a connection accepting loop. Whenever a connection is accepted, it checks the data stream for either the "ENCR" or "PTXT" flags, which control whether a TLS connection is created or not. The passed "clientHandler" function is called once the connection is completed.</li>
	<li>StartClient: Connects to a server, passes either the "ENCR" or "PTXT" flag as noted above, and returns the finished connection.</li>
</ul><br>
<p>Connections are returned as "ReadWriteClose" interfaces. Creating the pem and key certificate files is done via openssl. You can just google for examples.</p><br>

</div>

<b><a href="server.go">server.go</a>:</b>
```go
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
```

<b><a href="client.go">client.go</a>:</b>
```go
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
```

<b><a href="encryptedtcp/encryptedtcp.go">encryptedtcp/encryptedtcp.go</a>:</b>
```go
//A simple TCP client/server that can be encrypted (via tls) or not, depending on a flag passed from the client

package encryptedtcp

import ( "crypto/rand"; "crypto/tls"; "net"; "log" )

//Goes into a loop to accept clients. Returns a string on error
func StartServer(certFile, keyFile, listenOn string, clientHandler func(ReadWriteClose)) (error) {
	//Configure the certificate information
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return MyError{"Cannot Load Keys", err} }
	conf := tls.Config{Certificates:[]tls.Certificate{cert}, ClientAuth:tls.RequireAnyClientCert, Rand:rand.Reader}

	//Create the listener
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return MyError{"Cannot Listen", err} }
	defer listener.Close()

	//Listen and dispatch clients
	for {
		conn, err := listener.Accept()
		if err != nil {
			return MyError{"Cannot Accept Client", err} }
		go startHandleClient(conn, &conf, clientHandler)
	}

	//No error to return - This state is unreachable in the current library
	return nil
}

//Return the io stream for the connected client
func startHandleClient(conn net.Conn, conf* tls.Config, clientHandler func(ReadWriteClose)) {
	defer conn.Close()

	//Confirm encrypted connection flag (ENCR = yes, PTXT = no)
	isEncrypted := make([]byte, 4)
	amountRead, err := conn.Read(isEncrypted)
	if err != nil {
		log.Printf("Cannot get Encrypted Flag: %q\n", err); return }
	if amountRead != 4 {
		log.Printf("Cannot get Encrypted Flag: %q\n", "Invalid flag length"); return }
	if string(isEncrypted) == "PTXT" { //If plain text, just pass the net.Conn object to the client handler
		clientHandler(conn); return
	} else if string(isEncrypted) != "ENCR" { //If not a valid flag value
		log.Printf("Invalid flag value: %q\n", isEncrypted); return }

	//Initialize the tls session
	tlsconn := tls.Server(conn, conf)
	defer tlsconn.Close()
	if err := tlsconn.Handshake(); err != nil {
		log.Printf("TLS handshake failed: %q\n", err); return }

	//Pass the tls.Conn object to the client handler
	clientHandler(tlsconn)
}

//Start a client connection
func StartClient(certFile, keyFile, connectTo string, isEncrypted bool) (ReadWriteClose, error) {
	//Configure the certificate information
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, MyError{"Cannot Load Keys", err} }
	conf := tls.Config{Certificates:[]tls.Certificate{cert}, InsecureSkipVerify:true}

	//Connect to the server
	tcpconn, err := net.Dial("tcp", connectTo)
	if err != nil {
		return nil, MyError{"Cannot Connect", err} }

	//Handle unencrypted connections
	if !isEncrypted {
		tcpconn.Write([]byte("PTXT"))
		return tcpconn, nil //Return the base tcp connection
	}

	//Initialize encrypted connections
	tcpconn.Write([]byte("ENCR"))
	conn := tls.Client(tcpconn, &conf)
	conn.Handshake()

	//Confirm handshake was successful
	state := conn.ConnectionState()
	if !state.HandshakeComplete || !state.NegotiatedProtocolIsMutual {
		conn.Close()
		if !state.HandshakeComplete {
			return nil, MyError{"Handshake did not complete successfully", nil}
		} else {
			return nil, MyError{"Negotiated Protocol Is Not Mutual", nil} }
	}

	//Return the tls connection
	return conn, nil
}

//Error handling
type MyError struct {
	Context string
	TheError error
}
func (e MyError) Error() string {
	return e.Context+": "+e.TheError.Error(); }

//Interface for socket objects (read, write, close)
type ReadWriteClose interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}
```
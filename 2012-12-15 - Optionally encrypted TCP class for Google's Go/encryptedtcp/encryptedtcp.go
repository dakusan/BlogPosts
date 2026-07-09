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
var NewSock=new XMLSocket();
NewSock.onData=function(msg) { GetPacket(msg); }
NewSock.onConnect=function(Success) { SendInfoToJS(Success ? "Connected!" : "Could not connect"); }
SendInfoToJS(NewSock.connect(IP, Port) ? "Trying to Connect" : "Could not start connecting");
import flash.external.ExternalInterface;
import flash.events.Event;
ExternalInterface.addCallback("OpenSocket", OpenSocket);

function OpenSocket(IP:String, Port:Number):void
{
	SendInfoToJS("Trying to connect");
	var TheSocket:Socket = new Socket();
	TheSocket.addEventListener(Event.CONNECT, function(Success) { SendInfoToJS(Success ? "Connected!" : "Could not connect"); });
	TheSocket.addEventListener(Event.CLOSE, function() { SendInfoToJS("Connection Closed"); });
	TheSocket.addEventListener(IOErrorEvent.IO_ERROR, function() {SendInfoToJS("Could not connect");});
	TheSocket.addEventListener(ProgressEvent.SOCKET_DATA, function(event:ProgressEvent):void { ExternalInterface.call("GetPacket", TheSocket.readUTFBytes(TheSocket.bytesAvailable)); });
	TheSocket.connect(IP, Port);
}
function SendInfoToJS(str:String) { ExternalInterface.call("GetInfoFromFlash", str); }
stop();
/*
*         ReDEX.class
*         Coded: Timothy Strazzere
*
*         10/29/2008
*
*         Pass a .dex file as an argument and it will output the current
*         checksum and signature that is in the file, then it will output
*         the checksum and signature as it calculates them.
*
*/

import java.security.*;
import java.util.zip.Adler32;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;

public class ReDex {
	public static void main(String[] args) {
	args = new String[1];
	args[0]="C:\\riamer\\Finished Projects\\MiniAcuity\\com.castledragmire.testproject\\classes.dex";
	if (args.length == 1) {
	try {
	File file = new File(args[0]);
	
	byte[] barr = null;
	barr = getBytesFromFile(file);
	
	System.out.print("Original Checksum: ");
	for(int i = 8; i<12; i+=4)
	System.out.printf("0x%02X%02X%02X%02X ", barr[i+3], barr[i+2], barr[i+1], barr[i]);
	System.out.print("\nOriginal Signature: 0x");
	for(int i = 12; i<32; i+=4)
	System.out.printf("%02X%02X%02X%02X ", barr[i], barr[i+1], barr[i+2], barr[i+3]);
	
	calcSignature(barr);
	calcChecksum(barr);
	System.out.print("\n\nNew Checksum: ");
	for(int i = 8; i<12; i+=4)
	System.out.printf("0x%02X%02X%02X%02X ", barr[i+3], barr[i+2], barr[i+1], barr[i]);
	System.out.print("\nNew Signature: 0x");
	for(int i = 12; i<32; i+=4)
	System.out.printf("%02X%02X%02X%02X ", barr[i], barr[i+1], barr[i+2], barr[i+3]);
	}
	catch (Exception e) {
	System.err.println("File input error");
	}
	}
	else
	System.out.println("Invalid parameters");
	}
	
	public static byte[] getBytesFromFile(File file) throws IOException {
	InputStream is = new FileInputStream(file);
	
	// Get the size of the file
	long length = file.length();
	
	if (length > Integer.MAX_VALUE) {
	// File is too large
	}
	
	// Create the byte array to hold the data
	byte[] bytes = new byte[(int)length];
	
	// Read in the bytes
	int offset = 0;
	int numRead = 0;
	while (offset < bytes.length && (numRead=is.read(bytes, offset, bytes.length-offset)) >= 0) {
	offset += numRead;
	}
	
	// Ensure all the bytes have been read in
	if (offset < bytes.length) {
	throw new IOException("Could not completely read file "+file.getName());
	}
	
	// Close the input stream and return bytes
	is.close();
	return bytes;
	}
	private static void calcSignature(byte bytes[])
	{
	MessageDigest md;
	try
	{
	md = MessageDigest.getInstance("SHA-1");
	}
	catch(NoSuchAlgorithmException ex)
	{
	throw new RuntimeException(ex);
	}
	md.update(bytes, 32, bytes.length - 32);
	try
	{
	int amt = md.digest(bytes, 12, 20);
	if(amt != 20)
	throw new RuntimeException((new StringBuilder()).append("unexpected digest write:").append(amt).append("bytes").toString());
	}
	catch(DigestException ex)
	{
	throw new RuntimeException(ex);
	}
	}
	
	private static void calcChecksum(byte bytes[])
	{
	Adler32 a32 = new Adler32();
	a32.update(bytes, 12, bytes.length - 12);
	int sum = (int)a32.getValue();
	bytes[8] = (byte)sum;
	bytes[9] = (byte)(sum >> 8);
	bytes[10] = (byte)(sum >> 16);
	bytes[11] = (byte)(sum >> 24);
	}
}
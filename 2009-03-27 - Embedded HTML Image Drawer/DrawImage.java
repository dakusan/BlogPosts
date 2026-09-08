/***Dakusan's Embedded Image Drawer - v1.0 http://www.castledragmire.com/Projects/Other_Web_Scripts#DrawImage ***/

import java.applet.Applet;

//Image manupulation
import java.awt.Graphics;
import java.awt.Image;
import java.awt.image.MemoryImageSource;

//String manipulation
import java.util.StringTokenizer;

//Mouse events
import java.awt.event.MouseListener;
import java.awt.event.MouseEvent;
import java.awt.event.InputEvent;

//Communicating with JavaScript
import netscape.javascript.*;

public class DrawImage extends Applet implements MouseListener
{
	//Screen update functions
	Image CurrentImage; //The rendered image
	public void update(Graphics screen) { paint(screen); }
	public void paint(Graphics g) { g.drawImage(CurrentImage, 0, 0, this); } //Copy rendered image to the display
	
	//This takes an array of 24 bit pixels (encoded as a string separated by commas) and outputs it to the CurrentImage Image, then redraws the screen
	public void DrawFrame(Integer Width, Integer Height, String ImageData)
	{
		int TheWidth=Width.intValue(), TheHeight=Height.intValue(); //Change Width/Height values to native ints
		int NumPixels=TheWidth*TheHeight; //Calculate the number of pixels
		int PixelArray[]=new int[NumPixels]; //Allocate the pixel array
		StringTokenizer st = new StringTokenizer(ImageData, ","); //Split the encoded string list around commas
		for(int i=0;i<NumPixels;i++) //Loop over the encoded string list and extract each pixel's color
				PixelArray[i]=Integer.parseInt(st.nextToken()) | (255 << 24); //Set alpha channel to opaque
		CurrentImage=createImage(new MemoryImageSource(TheWidth, TheHeight, PixelArray, 0, TheWidth)); //Transfer the pixel data to an Image
		repaint(); //Update the screen
	}
	
	//Mouse events
	public void init() {addMouseListener(this);} //Init function initialized mouse listener
	public void mouseClicked(MouseEvent e) //When the mouse is clicked...
	{
		try { //Send mouse data (X, Y, If_Right_Button) to JavaScript "OnImageClickApplet" function
			JSObject Window=JSObject.getWindow(this);
			String Point[]={new Integer(e.getX()).toString(), new Integer(e.getY()).toString(), (((InputEvent)e).getModifiers()&InputEvent.BUTTON1_MASK)!=0 ? "0" : "1"};
			Window.call("OnImageClickApplet", Point);
			
		} catch(JSException err) { //If there was an error, output it to the screen
			Graphics g=this.getGraphics();
			g.clearRect(0, 0, this.WIDTH, this.HEIGHT);
			g.drawString(err.toString(), 0, 0);
		}
	}
	//Unused MouseListener events
	public void mouseReleased(MouseEvent e) {}
	public void mouseEntered(MouseEvent e) {}
	public void mouseExited(MouseEvent e) {}
	public void mousePressed(MouseEvent e) {}
}
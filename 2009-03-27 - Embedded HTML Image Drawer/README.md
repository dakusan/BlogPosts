<div>
Original Web Script: <a href="https://www.castledragmire.com/Projects/Other_Web_Scripts#DrawImage">Draw Image v1.0</a><br>
Date: 2009-03-27<br>
</div><br>

Draw JavaScript dynamically generated images to a web page.

See <a href="https://www.castledragmire.com/ProjectContent/WebScripts/DrawImage1.0/DrawFractal_Example.html">here</a> for live example.

## Information
<ul>
	<li>
		There are 2 user-called functions and 1 user-callback function for this script:
		<ul>
			<li>
				<b>CreateRenderObject</b>(<i>CodeBase</i>, <i>Width</i>, <i>Height</i>)
				<ul>
					<li>This creates an object where the pictures are drawn to; which needs to be added to the HTML document. In Internet Explorer, this is a Java Applet, and in all other browsers, its an IMG DOM object.</li>
					<li>The parameters are only used when in Internet Explorer</li>
					<li><i>CodeBase</i> is the base path where the Java applet resides</li>
				</ul>
			</li>
			<li>
				<b>DrawImage</b>(<i>ImageObject</i>, <i>Width</i>, <i>Height</i>, <i>PixelData</i>)
				<ul>
					<li>This draws the image to a render object</li>
					<li><i>ImageObject</i> is the object created through <b>CreateRenderObject</b></li>
					<li><i>PixelData</i> is a 1D array of color values (0xRRGGBB format) for each pixel (rows then columns)</li>
				</ul>
			</li>
			<li>
				<b>ImageClicked</b>(<i>X</i>, <i>Y,</i> <i>Flag</i>)
				<ul>
					<li>This is a callback function that occurs when an object created through <b>CreateRenderObject</b> is clicked.</li>
					<li><i>Flag</i> specifies if the Ctrl key is held down during the click in Internet Explorer, or if the right mouse button was used in all other browser.</li>
				</ul>
			</li>
		</ul>
	</li>
	<li>
		How this works:
		<ul>
			<li>In Internet Explorer, this sends all the data to a Java Applet that created a picture through the <b>createImage</b> and <b>MemoryImageSource</b> Java interfaces.</li>
			<li>In all other browsers, this sends the data to an IMG object by using the <a href="https://en.wikipedia.org/wiki/Data_URI_scheme">Data_URI_scheme</a>. Internet Explorer 8 supports this scheme, but only if the data set is smaller than 32KB, so I still consider it unsupported for Internet Explorer.</li>
		</ul>
	</li>
</ul>

Example:
<ul>
	<li><a href="https://www.castledragmire.com/Projects/Other_Web_Scripts#DrawImage">This</a> is an example of the technology which produces a Mandelbrot fractal using the same code base as my <a href="https://www.castledragmire.com/Projects/Julia_Fractals">Julia Fractals project</a>.</li>
	<li>It might be preferable to <a href="https://www.castledragmire.com/ProjectContent/WebScripts/DrawImage1.0/DrawFractal_Example.html">go directly to the example page</a> to remove size restrictions from the below IFRAME.</li>
	<li>The fractal dimensions, detail level, and area can be controlled by manually changing the values in the top table, or by clicking on the fractal per the below instructions.</li>
</ul>
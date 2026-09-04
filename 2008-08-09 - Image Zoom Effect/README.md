<div>
Original Web Script: <a href="https://www.castledragmire.com/Projects/Other_Web_Scripts#ImageZoomEffect">Image Zoom Effect v0.95</a><br>
Date: 2008-08-09<br>
</div><br>

An animated image zoom (growing/shrinking) effect for web pages that transitions from a displayed thumbnail into a larger version of the image. Animation occurs when mouse is moved over the thumbnail and shrinks when mouse leaves zoomed image. Zooming images attempt to stay within the window’s boundaries.

See <a href="https://www.castledragmire.com/Projects/Other_Web_Scripts#ImageZoomEffect">project page</a> for example.

<hr>

To use, the thumbnail must be an “img” object containing the class-name “GrowImage”, and its parent must be an anchor link (“a”) to the full size image. For example:
<pre><code>&lt;a href="<kbd>LINK_TO_FULL_SIZE_IMAGE.jpg</kbd>"&gt;&lt;img src="<kbd>LINK_TO_THUMBNAIL.jpg</kbd>" width=<kbd>THE_THUMBNAIL_WIDTH</kbd> height=<kbd>THE_THUMBNAIL_HEIGHT</kbd> class=GrowImage&gt;&lt;/a&gt;<br>
&lt;script type="text/javascript" src="ImageZoomEffect.js"&gt;&lt;/script&gt;</code></pre>
The ImageZoomEffect.js script <b>must be included after</b> all GrowImages have been added.<br><br>

Variables:
<ul>
	<li><kbd>RefreshTime</kbd>: Time (in milliseconds) between animation frame updates.</li>
	<li><kbd>TimeToAnimate</kbd>: The time (in milliseconds) to grow the animation to full size. For example, if this value is 100 (1/10th of a second) and <kbd>RefreshTime</kbd> is 10 (1/100th of a second), then there will be 10 frames of the image growing or shrinking at each 1/100th of a second interval.</li>
	<li><kbd>LoadingImageInfo</kbd>: The image to temporarily show while the full size image is being downloaded. An array of 3 values is used for this variable. “['<kbd>IMAGE_LINK</kbd>', <kbd>IMAGE_WIDTH</kbd>, <kbd>IMAGE_HEIGHT</kbd>]”</li>
</ul><br>

Technical details:
<ul>
	<li>After JS script is included, it searches through all images for ones that contain “GrowImage” in the className and adds appropriate events to the object.</li>
	<li>Zoomed images are placed inside an <i>absolutely positioned</i> container so it doesn’t affect the rest of the layout. The zoomed images themselves are <i>relatively positioned</i> inside the parent.</li>
	<li>Full size images are not loaded until the first time they are needed/requested by a mouse over action, and are unloaded when finished shrinking.</li>
	<li>A zoomed image currently stays centered over its thumbnail until it hits a window border, in which case it will not overlap the border if possible. If the image grows too large to fit within the window’s width or height, the image will overlap on the bottom and right sides.</li>
	<li>A zoomed image will start growing again if moused over while shrinking.</li>
	<li>There are still a few bugs I need to work out in this one... see the updates below.</li>
</ul>
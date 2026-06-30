<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Android_Stuff">Android Stuff</a><br>
Date: 2010-05-13<br>
Tags: Android, Examples, Java, Libraries, Phones, Security<br><br>

<p>Having recently finished my <a href="https://www.castledragmire.com/Projects/Mini_Acuity">first Android project</a> (and hopefully not last), I decided to supply some notes I took about the process.</p><br>
<p>While I am going to try and keep impressions to a minimum on the rest of this post, and keep it to tangible notes, I must first comment that trying to find out things for the <a href="http://en.wikipedia.org/wiki/Android_%28operating_system%29">Android platform</a> was often like pulling teeth. While its typical Java reference online documentation is all there with all the classes and cross-linking, that is about all it is, very dry and virtually useless beyond a base reference. The comments on variable parameters (and many other sections) in the reference are often coarse and not descriptive at all, for example, one parameter named <i>mask</i> has the basic description as “this is a mask”. Some functions don’t even have descriptions at all.</p>
<p>Perhaps I am getting too complacent as a programmer and getting used to excellent documentation like for Python or <a href="http://www.gtk.org/">GTK</a> (I’ve even grown to <b>love</b> Microsoft documentation after having used it for long enough!). After all, most required information is just a <a href="http://www.google.com">Google</a> away, and being a programmer is often just about <a href="https://www.castledragmire.com/Misc/Entertainment_Titles#SnowCrash">finding the proper magical incantations</a> to hook into a certain library. Unfortunately, however, even web searches were often yielding less than fruitful results when dealing with Android, as the platform is relatively new.</p><br><br>

<ul>
	<li>Some useful tasks and some problems:
	<ul>
		<li><b>Using the virtual (soft) keyboard without a TextView</b>:
		<ul>
			<li>
				<b>Showing the virtual keyboard</b>:<br>
				<code>((InputMethodManager)getSystemService(INPUT_METHOD_SERVICE)).toggleSoftInput(InputMethodManager.SHOW_FORCED, InputMethodManager.HIDE_IMPLICIT_ONLY);</code>
			</li>
			<li>
				<b>Hiding the virtual keyboard</b>:<br>
				<code>((InputMethodManager)getSystemService(INPUT_METHOD_SERVICE)).hideSoftInputFromWindow(getWindow().getDecorView().getApplicationWindowToken(), 0);</code><br>
				Note: “getWindow().getDecorView()” can also be replaced by a <i>View</i> on your screen
			</li>
			<li>
				<b>Getting the keyboard input</b>: Add the following function to the <i>Activity</i> that opened the keyboard:<br>
				<code>@Override public boolean onKeyDown(int keyCode, KeyEvent msg)</code><br>
				Note: This will not work if you’re not using default keyboard input (like if it’s set to enter Japanese or Chinese characters).
			</li>
		</ul></li>
		<li>
			<b>Determining the physical dimensions of the screen</b>:
			<p>This should be a trivial task using the <a href="http://developer.android.com/reference/android/util/DisplayMetrics.html">DisplayMetrics</a> (<code>getWindowManager().getDefaultDisplay()</code>) interface to get <a href="http://en.wikipedia.org/wiki/Dots_per_inch">dpi</a>s and multiply by the screen dimensions <code>getWindowManager().getDefaultDisplay().getWidth()</code> (and <i>.getHeight</i>). However, it doesn’t always work as it should.</p>
			<p>The best method to get the DPI would be to use “DisplayMetrics.xdpi” and “DisplayMetrics.ydpi”, but unfortunately, these are misreported by at least the Motorola Droid. I’ve found “DisplayMetrics.density”*160 to be pretty accurate, but if true accuracy is needed, a <a href="https://www.castledragmire.com/Projects/Mini_Acuity#CalibrateScreen">calibration screen</a> might be required.</p>
		</li>
		<li>
			<b>Inform user of touch events</b>: Many Android widgets (<i>Views</i>) change their visual state (highlight) when the user presses down on them to let the user know something is going to happen if the user lifts their finger while still on the widget. Unfortunately, there seems to be no text widget or <i>layout view</i> that does this automatic highlighting by itself (<i>ListView</i>s do in groups). The following is some example code to produce this effect.
			</li></ul></li></ul></div>

```java
import android.view.View.OnTouchListener;

public class CLASSNAME extends Activity //Replace CLASSNAME
{
	@Override public void onCreate(Bundle savedInstanceState)
	{
		View HighlightView=findViewById(R.id.THE_VIEWS_ID); //Replace THE_VIEWS_ID
		HighlightView.setOnTouchListener(HighlightState);
	}	
	
	private OnTouchListener HighlightState = new OnTouchListener() { public boolean onTouch(View v, MotionEvent event)
	{
		if(event.getAction()==MotionEvent.ACTION_DOWN)
			v.setBackgroundColor(0xFF0000FF); //Set background color to blue
		else if(event.getAction()==MotionEvent.ACTION_CANCEL || event.getAction()==MotionEvent.ACTION_UP)
			v.setBackgroundResource(0); //No background color
			
		return false;
	} };
}
```

<div><ul><ul>
		<li>
			<b>Retrieving the names and IDs of all resources in a resource group: </b>
			</li></ul></li></ul></div>

```java
import java.lang.reflect.Field;

Field[] FieldList=R.drawable.class.getDeclaredFields(); //replace "drawable"
String[] Names=new String[FieldList.length];
int[] IDs=new int[FieldList.length];
for(int i=0;i<FieldList.length;i++)
	IDs[i]=getResources().getIdentifier(Names[i]=FieldList[i].getName(), "drawable", getClass().getPackage().getName()); //replace "drawable"
```

<div><ul><ul>
		<li>
			<b>Setting a color matrix on an image</b>: If you have 2 ImageViews that display the same resource image, and either has a color matrix set on it, the will both share one of the color matrices. If this occurs, copy the image the resource, or use a separate image resource. For kicks, here is an example of setting an inverse color matrix on an image.<br>
			<code>((ImageView)findViewById(R.id.<b>IMAGE_ID</b>)).setColorFilter(new ColorMatrixColorFilter(new float[] {-1,0,0,0,255, 0,-1,0,0,255, 0,0,-1,0,255, 0,0,0,1,0}));</code>
		</li>
		<li>
			<b>Setting to full screen</b>:<br>
			<code>requestWindowFeature(Window.FEATURE_NO_TITLE); //This must be called before "setContentView", and hides the title bar
getWindow().setFlags(<b>FULLSCREEN</b> ? WindowManager.LayoutParams.FLAG_FULLSCREEN : 0, WindowManager.LayoutParams.FLAG_FULLSCREEN); //Turns on/off the status bar</code>
		</li>
		<li>
			<b>Starting another local activity</b>: Instead of using <i>Intent(String action)</i> for <i>Context.StartActivity</i>, as the <a href="http://developer.android.com/reference/android/content/Intent.html">Reference</a> explains, it is much easier to use <i>Intent(Context packageContext, Class&lt;?&gt; cls)</i> like the following: (called from inside an Activity)<br>
			<code>startActivity(new Intent(this, <b>OTHER_ACTIVITY_NAME</b>.class);</code>
		</li>
		<li>
			<b>Creating a timed event that updates the UI</b>: A function running through java.util.Timer cannot interact with the GUI. One solution to make a timer is with the android.os.Handler interface.
			</li></ul></li></ul></div>

```java
import android.os.Handler;

public class ExampleActivity extends Activity
{
	final int InitialDelay, RepeatDelay; //Set these variables
	Handler TimedHandler=new Handler();
	
	public void ExampleFunction()
	{
		TimedHandler.postDelayed(new Runnable() { public void run() {
			//Do GUI stuff...
			TimedHandler.postDelayed(this, RepeatDelay);
		} }, InitialDelay);
	}
}
```

<div><ul><ul>
			Another solution is to post to a Handler from the Timer function.
	</ul></li>
	<li>When dealing with putting on the market place:
	<ul>
		<li>Getting an account to put applications on the Android Market cost $25.</li>
		<li>Screenshots shown on the Android Market description page are somewhat buggy, and seemingly randomly either stretch properly or crop. Viewing the full sized screenshots does seem to work properly.</li>
		<li>
			Seeing as there are a number of applications on the market that have both a “Free” and “Full” version,  you’d think this would be an easy thing to accomplish. Unfortunately, the marketplace uses an application’s package name as its unique identifier, so both versions have to have a different package name, which is again, a bit of a nuisance.
			<p>One method of remedying this is just having a recursive string replace through all the files to change the package names. However, if using eclipse,  so you don’t have to reopen it, it’s quicker to update the string first in the manifest, and then renaming the package under the “src” folder by pressing F2 (rename) on it when it is selected.</p>
			<p>Also, unfortunately, if you do this, when a person upgrades from the lite to the full version, preferences are not automatically transferred 😕.</p>
		</li>
		<li>The <a href="http://market.android.com/publish/Home">publisher’s market place</a> page is very sparse and leaves a lot to be desired. It also seems to update only once every 24 hours or so (not sure of exact times).</li>
		<li>If an application is put up, it WILL get downloads immediately. For example, I put up an application with a description of “This is a test, do not download this” for doing security tests that I took down within like 10 minutes. It already had 2 comments/ratings on it within that time ~.~; .</li>
		<li><b>Google Checkout</b>: <a href="http://market.android.com/support/bin/answer.py?answer=112622">Fees</a>. When a copy of your application is purchased, the user has 24 hours to return it. The money is not deposited into your bank account until after this time (when it’s not a weekend). If you want to give your application to someone for free, they need to purchase it through the market, and then you can cancel the purchase transaction before the 24 hours are up. Unfortunately, this has to be done every time they want to update the application. It also seems you cannot buy your own applications, as the purchase server throws an error.</li>
	</ul></li>
	<li id=ApplicationProtection>Application Protection:
		<p>You can download any Android application by default from your phone to your computer, modify them, and reinstall them back to any phone. An example use for this would be to crack a shareware application where just a single byte probably needs to be changed to make it a full version.</p>
		<p>The applications themselves are in an .apk file (which is just a .zip file), and the source code (classes) are encoded as a “Dalvik Executable” file within it (classes.dex), which as I understand it, is optimized Java bytecode. So, AFAIK right now, there is no way to decompile the .dex file back to the original source, like you can with normal Java. However, the <a href="http://developer.android.com/guide/developing/tools/emulator.html">Android emulator</a>, part of the <a href="http://developer.android.com/sdk/index.html">Android SDK</a>, includes a tool called <i>dexdump</i>, which allows you to decompile it to bytecode.</p>
		<p>Once you have the bytecode, you can use that as reference to modify the compiled .dex file however you want, which is pretty much synonymous with assembly editing. Once that is done, the <a href="http://strazzere.com/blog/?p=3">signature and checksum of the .dex file must be recalculated</a> (<a href="ReDex.java">Java source by Timothy Strazzere</a>), and then the apk file must be <a href="http://developer.android.com/guide/publishing/app-signing.html">resigned</a>, and then it’s good to go!</p>
		<p>The marketplace also has an option to turn on <i>Copy Protection</i>. When this is turned on for an application, the user cannot backup or access the applications package file. I would assume however with a rooted phone you could still grab it from “/data/app-private”, and the rest of the process should be the same. I have not tested this as rooting Android 2.1 is much more of a pain in the butt, ATM, than I want to deal with.</p>
	</li>

</ul></div>
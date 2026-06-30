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
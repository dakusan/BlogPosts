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
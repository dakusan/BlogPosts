//Showing the virtual keyboard
((InputMethodManager)getSystemService(INPUT_METHOD_SERVICE)).toggleSoftInput(InputMethodManager.SHOW_FORCED, InputMethodManager.HIDE_IMPLICIT_ONLY);

//Hiding the virtual keyboard
((InputMethodManager)getSystemService(INPUT_METHOD_SERVICE)).hideSoftInputFromWindow(getWindow().getDecorView().getApplicationWindowToken(), 0);

//Getting the keyboard input: Add the following function to the Activity that opened the keyboard
@Override public boolean onKeyDown(int keyCode, KeyEvent msg) { }

//Setting a color matrix on an image
((ImageView)findViewById(R.id.IMAGE_ID)).setColorFilter(new ColorMatrixColorFilter(new float[] {-1,0,0,0,255, 0,-1,0,0,255, 0,0,-1,0,255, 0,0,0,1,0}));

//Setting to full screen
requestWindowFeature(Window.FEATURE_NO_TITLE); //This must be called before "setContentView", and hides the title bar
getWindow().setFlags(FULLSCREEN ? WindowManager.LayoutParams.FLAG_FULLSCREEN : 0, WindowManager.LayoutParams.FLAG_FULLSCREEN); //Turns on/off the status bar

//Starting another local activity
startActivity(new Intent(this, OTHER_ACTIVITY_NAME.class);
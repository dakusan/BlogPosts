//Run a function on primary thread
public static void RunOnPrimary(Runnable R)
{
    Looper L=MainActivity.getMainLooper();
    //Start commenting here so that items are always added to the queue, forcing in-order processesing
    if(Looper.myLooper()==Looper.getMainLooper())
        R.run();
    else
    //End commenting here
        new Handler(Looper.getMainLooper()).post(R);
}
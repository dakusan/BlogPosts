<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Optimization_gone_bad">Optimization gone bad</a><br>
Date: 2015-11-22<br>
SubTitle: Or, the case of the Android app out-of-order calls<br>
Tags: Android, Concurrency, Examples, Java<br><br>

<link href="layoutx/github-markdown-new.css" rel=stylesheet type="text/css">
<p>On Android, there is a primary thread which runs all UI stuff. If a GUI operation is ran in a different thread, it just won't work, and may throw an error. If you block this thread with too much processing... well... bad things happen. Due to this design, you have to push all UI operations to this main thread using <a href="https://developer.android.com/reference/android/os/Looper.html">Looper</a>.<a href="https://developer.android.com/reference/java/lang/Thread.html#run()">run()</a>.</p>
<p><a href="https://developer.android.com/reference/java/lang/Runnable.html">Runnables</a> pushed to this thread are always ran in FIFO execution order, which is a useful guarantee for programming.</p>
<p>So I decided to get smart and create the following function to add asynchronous calls that needed to be run on the primary thread. It takes a Runnable and either runs it immediately, if already on the Primary thread, or otherwise adds it to the Primary Thread’s queue.</p>
</div>

```java
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
```

<div>
<p>I was getting weird behaviors though in a part of the project where some actions pulled in from JavaScript were happening before subsequent actions.
After the normal debugging one-by-one steps to figure it out, I realized that <i>MAYBE</i> some of the JavaScript calls were, for some bizarre reason, already running on the primary thread. In this case they would run immediately, before the queued items started coming in. This turned out to be the case, so I ended up having to comment out the first 3 lines after the function’s first comment (if/R.run/else), and it worked great.</p>
<p>I found it kind of irritating that I had to add actions to the queue when it could have been run immediately on the current thread, but oh well, I didn’t really have a choice if I want to make sure everything is always run in order across the system.</p></div>
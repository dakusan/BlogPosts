<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Thread_synchronization_in_C%23">Thread synchronization in C#</a><br>
Date: 2013-01-09<br>
SubTitle: Building the wheel that should have already existed<br>
Tags: CSharp, Examples, Synchronization, Threading<br><br>

<p>I have been working heavily in C# CE (Compact Edition) v2.0 for the last 2 years for clients, and one of the very many things that I was never really happy with (in at least that version of the language, though it looks like it might plague all versions of C#) is the available thread synchronization tools. I’ve come to love the lock/wait/notify model (in Java it’s synchronized/wait/notify and in Perl it’s <a href="http://perldoc.perl.org/threads/shared.html">lock/cond_wait/cond_signal</a>), but I have found nothing as intuitive and safe to use in C#. To alleviate this, I went ahead and wrote my own ThreadLockAndWait class that achieves this functionality.</p><br>

This works the same as the POSIX lock, unlock, cond_wait, cond_signal, and cond_timedwait functions, except:
<ul>
	<li>Lock is not required before CondSignal (it does its own inner lock and unlock)</li>
	<li>If ReacquireLockAfterWait is false, in which case CondWait will not lock again after signaled and just continue immediately</li>
	<li>Only 1 thread can be CondWaiting at a time (If one is CondWaiting and is signaled but not reacquired the lock, its ok for another to start CondWaiting)</li>
</ul><br>

</div>

```cs
public class ThreadLockAndWait
{
	private Mutex TheLock=new Mutex(), CondWaitLock=new Mutex(); //CondWaitLock makes sure 1 thread stops waiting before the next one starts waiting
	private ManualResetEvent WaitTimer=new ManualResetEvent(false);
	private string OwnersThreadName=null;
	private int OwnerLockCount=0;

	public void Lock()
	{
		TheLock.WaitOne();
		if(OwnerLockCount++==0)
			OwnersThreadName=Thread.CurrentThread.Name;
	}
	public void UnLock()
	{
		TheLock.WaitOne();
		if(OwnerLockCount==0)
		{
			TheLock.ReleaseMutex();
			throw new Exception("Cannot unlock if not locked");
		}
		TheLock.ReleaseMutex();
		if(--OwnerLockCount==0)
			OwnersThreadName=null;
		TheLock.ReleaseMutex();
	}
	public void CondWait() { RealCondWait(-1, true); }
	public void CondWait(bool ReacquireLockAfterWait) { RealCondWait(-1, ReacquireLockAfterWait); }
	public void CondTimedWait(int TimeToWait) { RealCondWait(Math.Max(0, TimeToWait), true); }
	public void CondTimedWait(int TimeToWait, bool ReacquireLockAfterWait) { RealCondWait(Math.Max(0, TimeToWait), ReacquireLockAfterWait); }
	private void RealCondWait(int TimeToWait, bool ReacquireLockAfterWait)
	{
		//Prepare to wait
		TheLock.WaitOne();
		if(OwnerLockCount==0)
		{
			TheLock.ReleaseMutex();
			throw new Exception("Cannot wait if not locked");
		}
		CondWaitLock.WaitOne(); //Release this wait before the next one starts
		WaitTimer.Reset();
		TheLock.ReleaseMutex();

		//Release all locks
		int PreviousLockCount=OwnerLockCount;
		OwnersThreadName=null;
		OwnerLockCount=0;
		if(PreviousLockCount!=1)
			System.Diagnostics.Debug.Print("Warning, mutex has multiple locks from thread!");
		for(int i=0;i<PreviousLockCount;i++)
			TheLock.ReleaseMutex();

		//Wait
		if(TimeToWait>0)
			WaitTimer.WaitOne(TimeToWait, false);
		else if(TimeToWait!=0)
			WaitTimer.WaitOne();
		CondWaitLock.ReleaseMutex();

		//Reacquire lock
		if(!ReacquireLockAfterWait)
			return;
		for(int i=0;i<PreviousLockCount;i++)
			TheLock.WaitOne();
		OwnerLockCount=PreviousLockCount;
		OwnersThreadName=Thread.CurrentThread.Name;
	}

	public void CondSignal()
	{
		TheLock.WaitOne();
		WaitTimer.Set();
		TheLock.ReleaseMutex();
	}
}
```

<div>
</div>
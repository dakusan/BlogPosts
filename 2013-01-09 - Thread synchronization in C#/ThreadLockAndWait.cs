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
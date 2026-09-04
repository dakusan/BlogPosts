/***Copyright and coded by Dakusan, please do not steal code, but you may study and learn from it... Free to use and redistribute - See https://www.castledragmire.com/Copyright for more information. ***/
/***Dakusan's Image Zoom Effect - v0.95 https://www.castledragmire.com/Projects/Other_Web_Scripts ***/

//Standard functions
var
IsIE=(navigator.appName.indexOf("Microsoft")!=-1),
$CP =function(O,P)	{for(var i in P)eval("O."+i+"=P[i];");return O;}, //Copy Properties
$GE =function(O)	{return typeof(O)=="string"?$(O):O;}, //GetElement by String or Object
$CE =function(T,P,A)	{T=document.createElement(IsIE&&typeof(P)=="object"&&'name' in P?'<'+T+' name="'+P.name+'">':T);if(P)$CP(T,P);if(A)$GE(A).appendChild(T);return T;}, //CreateElement of Type "T", Copy Properties "P", and Append to "A"
$SE =function(O,E,F)	{ if(IsIE) eval("O[E]=function() {"+F+"};"); else O.setAttribute(E,F); }; //SetEvent
function I(N) { return parseInt(N); }

var GI={
	//User variables
	RefreshTime:10, //Milliseconds between updates/checks
	TimeToAnimate:100, //Time to grow an image
	LoadingImageInfo:['layoutx/Loading.png',104,30], //Loading Picture: Path, width, height
	//***************DO NOT EDIT PAST HERE***************
	//Internal Variables
	Cur:[], //List of current images that are grown
	//Functions
	Init:function() //Search for and create all grow images
	{
		//Search all images for ones with class containing "GrowImage"
		var MyImgs=document.getElementsByTagName('img');
		for(var i=0;i<MyImgs.length;i++)
			if(MyImgs[i].className.indexOf('GrowImage')!=-1) //When found, Turn into a grow image box
			{
				$SE(MyImgs[i], 'onmouseover', 'GI.Grow(this);'); //When mouse over, start growing
				$SE(MyImgs[i], 'onmouseout', 'GI.EndGrow(this);'); //If mouse out, cancel growing
			}
	},

	Find:function(O, In) //Find the current info in the GrowImage list array - Parameters: ObjectToFind, Index to search
	{
		for(var i=0;i<GI.Cur.length;i++)
			if(GI.Cur[i][In]==O)
				return i;
		return -1;
	},

	Grow:function(O) //Create the growing image effect
	{
		//Confirm that our object hasn't already been created
		if(GI.Find(O, 0)!=-1)
			return;
	
		//Add object to the GrowImage list array
		var GIInfo=[O, null, O.offsetWidth, O.offsetHeight, 0, 0, 0, -1, null]; //Items are in the format [OriginalImageObject, NewImageObject, OldWidth, OldHeight, NewWidth, NewHeight, CurrentStepInGrowing, State (-1=Loading, 0=Idle, 1=Growing, 2=Shrinking), IntervalAction, ParentX, ParentY]
		GI.Cur.push(GIInfo);
	
		//Create and find elements
		var P=O.parentNode, //Find the parent
		ElCont=P.insertBefore($CE('span', {'style.position':'absolute'}), O),  //Create the growing images's container which is absolute so it doesn't affect the rest of the document layout
		LoadEl=$CE('img', {src:GI.LoadingImageInfo[0], 'style.position':'relative', 'style.left':I((GIInfo[2]-GI.LoadingImageInfo[1])/2)+'px', 'style.top':I((GIInfo[3]-GI.LoadingImageInfo[2])/2)+'px', 'border':0}, ElCont), //Create a temporary loading image
		NewEl=$CE('img', {src:P.href, 'style.visibility':'hidden', 'style.position':'absolute', 'onmouseout':function() {GI.NewAnimate(this,2);}, 'onmouseover':function() {GI.NewAnimate(this,1);}}, ElCont);
	
		//Copy over new properties into our info object
		GIInfo[1]=NewEl;
		GIInfo[8]=setInterval(function() { GI.StartGrow(NewEl); }, GI.RefreshTime);
		
		//Determine the objects positioning
		var MyX=0, MyY=0, CurEl=O;
		while(CurEl!=null && CurEl.tagName!='BODY')
		{
			MyX+=CurEl.offsetLeft;
			MyY+=CurEl.offsetTop;
			CurEl=CurEl.offsetParent;
		}
		GIInfo[9]=MyX;
		GIInfo[10]=MyY;
	},

	EndGrow:function(O) //If moused out of image before image is loaded, cancel grow effect
	{
		//If growing and still loading, stop the process
		var i=GI.Find(O, 0);
		if(i!=-1 && GI.Cur[i][7]==-1)
			GI.Remove(O,true);
	},
	
	StartGrow:function(O) //If grow image has completed loading, start growing it
	{
		if(O.offsetWidth==0) //If image has not yet loaded, exit
			return;

		//Find the GI info in the current list
		var i=GI.Find(O, 1);
		if(i==-1) //If not found, there was a problem and interval will continue indefinitely since we don't have its handle... shouldn't happen though
			return;
		var CurGI=GI.Cur[i];

		//Prepare to start growing
		clearInterval(CurGI[8]);	//Stop the loading interval
		CurGI[4]=O.offsetWidth;		//Get the full size width
		CurGI[5]=O.offsetHeight;	//Get the full size height
		CurGI[7]=1;			//Set the State to Growing
		O.parentNode.removeChild(O.previousSibling); //Remove the loading image
		O.width=CurGI[2];		//Set the starter width
		O.height=CurGI[3];		//Set the starter height
		O.style.visibility='visible';	//show the object
	},

	NewAnimate:function(O,State) //Start shrinking or growing again
	{
		//Find the GI info in the current list
		var i=GI.Find(O, 1); //Find the GI info in the current list
		if(i==-1) //If not found, there was a problem... shouldn't happen though
			return;
		var CurGI=GI.Cur[i];

		//If not idle, stop current animation
		if(CurGI[7]!=0)
			clearInterval(CurGI[8]);

		//Start the shrinking animation
		CurGI[7]=State; //Set the new State
		CurGI[8]=setInterval(function() { GI.DoAnimate(O); }, GI.RefreshTime); //Start the animation
	},

	DoAnimate:function(O) //Shrink or Grow animation
	{
		//Find the GI info in the current list
		var i=GI.Find(O, 1); //Find the GI info in the current list
		if(i==-1) //If not found, there was a problem and interval will continue indefinitely since we don't have its handle... shouldn't happen though
			return;
		var CurGI=GI.Cur[i];

		//Grow the image
		var DP=[document.documentElement.scrollLeft, document.documentElement.scrollTop, document.documentElement.clientWidth, document.documentElement.clientHeight, CurGI[9], CurGI[10]]; //Document positioning numbers
		CurGI[6]=Math.max(Math.min(CurGI[6]+GI.RefreshTime*(CurGI[7]==1 ? 1 : -1), GI.TimeToAnimate),0); //Update grow counter (+for grow, -for shrink) - clamp from 0 to TimeToAnimate
		var Progress=CurGI[6]/GI.TimeToAnimate;			//Calculate the progress
		O.style.display='none';					//Hide so multiple refreshes are not needed
		O.width=I(CurGI[2]+(CurGI[4]-CurGI[2])*Progress);	//Calculate new width
		O.height=I(CurGI[3]+(CurGI[5]-CurGI[3])*Progress);	//Calculate new height
		O.style.left=Math.max(DP[0]-DP[4],Math.min(DP[2]-O.width -DP[4]+DP[0],I((CurGI[2]-O.width )/2)))+'px'; //Recenter horizontally - but keep within windows bounds
		O.style.top =Math.max(DP[1]-DP[5],Math.min(DP[3]-O.height-DP[5]+DP[1],I((CurGI[3]-O.height)/2)))+'px'; //Recenter vertically - but keep within windows bounds
		O.style.display='block';				//Show the image again

		if((CurGI[7]==1 && CurGI[6]!=GI.TimeToAnimate) || (CurGI[7]==2 && CurGI[6]!=0)) //If still growing or shrinking, exit
			return;

		//If done growing, stop the animating
		var IsGrowing=(CurGI[7]==1);
		CurGI[7]=0; //Set the State to Idle
		clearInterval(CurGI[8]); //Stop the animation interval
		if(IsGrowing) //If growing, exit here
			return;
		//If shrinking, delete GrowImage objects and data
		GI.Remove(O,false);
	},

	Remove:function(O,Loading) //Delete current GrowImage data and objects
	{
		//Find the GI info in the current list
		var i=GI.Find(O, 1); //Find the GI info in the current list
		if(i==-1) //If not found, there was a problem... shouldn't happen though
			return;
		var CurGI=GI.Cur[i], P=CurGI[1].parentNode;

		clearInterval(CurGI[8]); //Stop the current interval
		if(Loading) P.removeChild(CurGI[1].previousSibling); //Remove the loading image - technically not needed since we are unloading the parent, but helps prevent memory leaks 
		P.removeChild(CurGI[1]); //Remove the grow image - technically not needed since we are unloading the parent, but helps prevent memory leaks
		P.parentNode.removeChild(P); //Remove the image container
		GI.Cur.splice(i, 1); //Remove object from the GrowImage list array
	}
};

GI.Init();
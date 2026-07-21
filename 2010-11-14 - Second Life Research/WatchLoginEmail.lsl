//Users to watch
key DetectPersons=[ //List of UIDs of users to watch. (Real UIDs redacted)
	"fdf1fbff-f19f-ffff-ffff-ffffffffffff", //Person 1
	"f0fffaff-f61f-ffff-ffff-ffffffffffff" //Person 2
];

//Other Global Variables
integer NumUsers;
integer UsersParsed=0;
list UserNames;
list Status;

default
{
	state_entry()
	{
		NumUsers=llGetListLength(DetectPersons); //Number of users to watch

		//Get User Names
		integer i;
		for(i=0;i<NumUsers;i++)
		{
			llListInsertList(UserNames, [''], i);
			llListInsertList(Status, [0], i);
			llRequestAgentData(llList2Key(DetectPersons, i), DATA_NAME);
		}
	}

	dataserver(key requested, string data)
	{
		//Find User Position
		integer i;
		for(i=0;i<NumUsers;i++)
			if(llList2Key(DetectPersons, i)==requested)
				llListReplaceList(UserNames, [data], i, 1);

		if(++UsersParsed==NumUsers)
			state Running;
	}
}

state Running
{
	state_entry()
	{
		llOwnerSay((string)UserNames);
		llOwnerSay((string)Status);
		llSetTimerEvent(30);
	}

	timer()
	{
		llRequestAgentData(DetectPerson, DATA_ONLINE);
	}

	dataserver(key requested, string data)
	{
		if(data==IsOnline)
			return;
		IsOnline=data;
		if(data=="0")
			return;
		string Message="The user you are watching '"+UserName+"' signed on at "+llGetTimestamp();
		llEmail(EMAIL_ADDRESS, "User Signed on", Message);
		llOwnerSay(Message);
	}
}
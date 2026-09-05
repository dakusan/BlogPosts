<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Second_Life_Research">Second Life Research</a><br>
Date: 2010-11-14<br>
Tags: Examples, Games, Reverse Engineering, Second Life, Security<br><br>
</div>

<p>Back in May of 2007 one of my friends got me onto <a href="http://en.wikipedia.org/wiki/Second_Life">Second Life</a>, the first and only MMORPG I’ve touch since <a href="https://www.castledragmire.com/ragnarok">my Ragnarok days</a>. While Second Life had a strong pull for me due to its similarities to The MetaVerse in <a href="https://www.castledragmire.com/Misc/Entertainment_Titles#SnowCrash">Snow Crash</a>, my favorite book, I was of course more drawn to playing with the Engine and seeing what I could do with it.</p>
<p>I felt no real need to delve into the code or packet level of the client as it was open source, so I stayed mostly on the scripting level side of things in the world. IIRC I did find at least a dozen major security holes, but I unfortunately cannot seem to find logs of my research 🙁.</p>
I do however remember at least 2 of the security holes I found:
<ul>
	<li>While an avatar could not pass through solid walls normally, if an object was visible that allowed “sitting” beyond the walls, the user could issue the sit command on that object which transported the avatar past the barriers.</li>
	<li>While there were optional restrictions on areas pertaining to if/where an object could be placed, once an object was placed somewhere, it could be “pushed” to almost any other location no matter the restrictions. When an object was pushed into another area beyond where it was placed, it was still inventoried as being in the originally placed location, but could interact with the world at the location it was actually at. Objects could even pass through solid barriers if the proper push velocities were given. The only way at the time to combat this was to have whole private islands as blocking anonymous objects. This security hole opened up multiple other security holes including:
		<ul>
			<li>If a user “sat” on the object, they could get to anywhere the object could.</li>
			<li>These objects could be used to interact with the immediate world around them, including repeating private conversations in a private area.</li>
		</ul>
	</li>
</ul><br>

<p>I had also at the time planned on writing an application that allowed hijacking and reuploading any encountered texture or construct, which was trivial due to the open nature of the system. I never did get around to it for two reasons. First, I got distracted by other projects, and second, because it could have seriously destabilized the Second Life economy, which was built around selling said textures and constructs. I actually liked what Second Life was trying to accomplish and had no wish of making <a href="http://en.wikipedia.org/wiki/Linden_Lab">Linden Lab</a>’s life harder or ruining the experiment in open economy.</p><br>

<p>I was however able to find a few pieces of my research and scripts that I figured I could post here. First, I do not recall what I did to find this, but the entire list of pre-defined “Last Names” was accessible, and IIRC the proprietary last names could be used for character creation if you knew how to access them (not 100% sure if this latter hack was available). Here was the list as of when I acquired it in 2007. I had the list separated into two columns, and <b><i>I think</i></b> they were “open” names and “proprietary” names. Each name is followed by its identifier.</p>
<p align=center><b><ins>Open Names</ins></b></p>
Congrejo(339), Spitteler(957), Boucher(1716), Kohime(2315), Korobase(2363), Bingyi(3983), Hyun(3994), Qunhua(4003), Yiyuan(4010), Nikolaidis(4032), Bikcin(4040), Laryukov(4112), Bamaisin(4127), Choche(4136), Ultsch(4140), Coage(4164), Cioc(4173), Barthelmess(4212), Koenkamp(4322), Daviau(4340), Menges(4345), Beaumont(4390), Lubitsch(4392), Taurog(4408), Negulesco(4418), Beresford(4466), Babenco(4468), Catteneo(4483), Dagostino(4509), Ihnen(4511), Basevi(4517), Gausman(4530), Heron(4533), Fegte(4535), Huldschinsky(4539), Juran(4543), Furse(4548), Heckroth(4550), Perfferle(4552), Reifsnider(4553), Hotaling(4559), DeCuir(4560), Carfagno(4561), Mielziner(4573), Bechir(4592), Zehetbauer(4615), Roelofs(4624), Hienrichs(4647), Rau(4654), Oppewall(4657), Bonetto(4659), Forwzy(4677), Repine(4680), Fimicoloud(4685), Bleac(4687), Anatine(4688), Gynoid(4745), Recreant(4748), Hapmouche(4749), Ceawlin(4758), Balut(4760), Peccable(4768), Barzane(4778), Eilde(4783), Whitfield(4806), Carter(4807), Vuckovic(4808), Rehula(4809), Docherty(4810), Riederer(4811), McMahon(4812), Messmer(4813), Allen(4814), Harrop(4815), Lilliehook(4816), Asbrink(4817), Laval(4818), Dyrssen(4819), Runo(4820), Uggla(4822), Mayo(4823), Handrick(4824), Grut(4825), Szondi(4826), Mannonen(4827), Korhonen(4828), Beck(4829), Nagy(4830), Nemeth(4831), Torok(4832), Mokeev(4833), Lednev(4834), Balczo(4835), Starostin(4836), Masala(4837), Rasmuson(4838), Martinek(4839), Mizser(4840), Zenovka(4841), Dovgal(4842), Capalini(4843), Kuhn(4845), Platthy(4846), Uriza(4847), Cortes(4848), Nishi(4849), Rang(4850), Schridde(4851), Dinzeo(4852), Winkler(4853), Broome(4854), Coakes(4855), Fargis(4856), Beerbaum(4857), Pessoa(4858), Mathy(4859), Robbiani(4860), Raymaker(4861), Voom(4862), Kappler(4863), Katscher(4864), Villota(4865), Etchegaray(4866), Waydelich(4867), Johin(4868), Blachere(4869), Despres(4871), Sautereau(4872), Miles(4873), Lytton(4874), Biedermann(4875), Noel(4876), Pennell(4877), Cazalet(4878), Sands(4879), Tatham(4880), Aabye(4881), Soderstrom(4882), Straaf(4883), Collas(4884), Roffo(4885), Sicling(4886), Flanagan(4887), Seiling(4888), Upshaw(4889), Rodenberger(4890), Habercom(4891), Kungler(4892), Theas(4893), Fride(4894), Hirons(4895), Shepherd(4896), Humphreys(4897), Mills(4898), Ireton(4899), Meriman(4900), Philbin(4901), Kidd(4902), Swindlehurst(4903), Lowey(4904), Foden(4905), Greggan(4906), Tammas(4907), Slade(4908), Munro(4909), Ebbage(4910), Homewood(4911), Chaffe(4912), Woodget(4913), Edman(4914), Fredriksson(4915), Larsson(4916), Gustafson(4917), Hynes(4918), Canning(4919), Loon(4920), Bekkers(4921), Ducatillon(4923), Maertens(4924), Piek(4925), Pintens(4926), Jansma(4927), Sewell(4928), Wuyts(4929), Hoorenbeek(4930), Broek(4931), Jacobus(4932), Streeter(4933), Babii(4934), Yifu(4935), Carlberg(4936), Palen(4937), Lane(4938), Bracken(4939), Bailey(4940), Morigi(4941), Hax(4942), Oyen(4943), Takacs(4944), Saenz(4945), Lundquist(4946), Tripsa(4947), Zabelin(4948), McMillan(4950), Rosca(4951), Zapedzki(4952), Falta(4953), Wiefel(4954), Ferraris(4955), Klaar(4956), Kamachi(4957), Schumann(4958), Milev(4959), Paine(4960), Staheli(4961), Decosta(4962), Schnyder(4963), Umarov(4964), Pinion(4965), Yoshikawa(4966), Mertel(4967), Iuga(4968), Vollmar(4969), Dollinger(4970), Hifeng(4971), Oh(4972), Tenk(4973), Snook(4974), Hultcrantz(4975), Barbosa(4976), Heberle(4977), Dagger(4978), Amat(4979), Jie(4980), Qinan(4981), Yalin(4982), Humby(4983), Carnell(4984), Burt(4985), Hird(4986), Lisle(4987), Huet(4988), Ronmark(4989), Sirbu(4990), Tomsen(4991), Karas(4992), Enoch(4993), Boa(4994), Stenvaag(4995), Bury(4996), Auer(4997), Etzel(4998), Klees(4999), Emmons(5000), Lusch(5001), Martynov(5002), Rotaru(5003), Ballinger(5004), Forcella(5005), Kohnke(5006), Kurka(5007), Writer(5008), Debevec(5009), Hirvi(5010), Planer(5011), Koba(5012), Helgerud(5013), Papp(5014), Melnik(5015), Hammerer(5016), Guyot(5017), Clary(5018), Ewing(5019), Beattie(5020), Merlin(5021), Halasy(5022), Rossini(5024), Halderman(5025), Watanabe(5026), Bade(5027), Vella(5028), Garrigus(5029), Faulds(5030), Pera(5031), Bing(5032), Singh(5033), Maktoum(5034), Petrov(5035), Panacek(5036), Dryke(5037), Shan(5038), Giha(5039), Graves(5040), Benelli(5041), Jun(5042), Ling(5043), Janus(5044), Gazov(5045), Pfeffer(5046), Lykin(5047), Forder(5048), Dench(5049), Hykova(5050), Gufler(5051), Binder(5052), Shilova(5053), Jewell(5054), Sperber(5055), Meili(5056), Matova(5057), Holmer(5058), Balogh(5059), Rhode(5060), Igaly(5061), Demina(5062)<br><br>

<p align=center><b><ins>Proprietary Names</ins></b></p>
ACS(1353), FairChang(1512), Teacher(2186), Learner(2213), Maestro(2214), Aprendiz(2215), Millionsofus(2746), Playahead(2833), RiversRunRed(2834), SunMicrosystems(2836), Carr(2917), Dell(3167), Reuters(3168), Hollywood(3173), Sheep(3471), YouTopia(3816), Hillburn(3817), Bradford(3820), CiscoSystems(3958), PhilipsDesign(3959), MadeVirtual(4205), DuranDuran(4210), eBay(4665), Vodafone(4666), Xerox(4667), TGDev(4668), Modesto(4669), Sensei(4670), Ideator(4671), Autodesk(4789), MovieTickets(4790), AvaStar(4791), DiorJoaillerie(4793), AOL(4795), Gabriel(4805), Tequila(5064), Loken(5065), Matlin(5066), GeekSquad(5067), Bradesco(5068), CredicardCiti(5069), PontiacGXP(5070), KAIZEN(5071), McCain(5072), Schomer(5074), Showtime(5075), OzIslander(5076), Meltingdots(5077), Allanson(5083), Sunbelter(5084), SaxoBank(5085), Esslinger(5086), Stengel(5087), Lemeur(5088), Tsujimoto(5089), KaizenGames(5090), Uphantis(5091), OurVirtualHolland(5092), McKinseyandCompany(5093), Lempert(5094), Affuso(5095), Gkguest(5096), Eye4You(5097), OShea(5098), Citibank(5099), Citicard(5100), Citigroup(5101), Citi(5102), Credicard(5103), Diners(5104), Citifinancial(5105), CitiBusiness(5106), BnT(5107), Yensid(5108), Helnwein(5111), Grindstaff(5112), Shirk(5113), SolidWorks(5114), Storm(5115), CarterFinancial(5116), Parkinson(5117), Lear(5118), FiatBrasil(5119), RossiResidencial(5120), Brooklintolive(5121), Calmund(5123), Briegel(5124), Herde(5125), Pfetzing(5126), Triebel(5127), Roemer(5128), Reacher(5129), Thomas(5130), Fraser(5131), Gabaldon(5132), NBA(5133), Accubee(5134), Brindle(5135), Searer(5136), Ukrop(5137), Ponticelli(5138), Belcastro(5139), Glin(5140), Rice(5141), DavidStern(5142), Totti(5144), onrez(5145), DeAnda(5146), Grandi(5147), Pianist(5148), osMoz(5149), PaulGee(5150)<br><br>

<p>The second piece I was able to find was a script I used to alert me via email whenever one of my friends signed on. I have unfortunately not tested this script before posting it as I no longer have Second Life installed or wish to waste the time testing it, but here it is none the less. 😅</p>
</div>

```lsl
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
```

<p>Of course all this research was from 2007 and I have no idea what is capable now. I do really hope though that they at least updated the client’s interface because it was <b>incredibly</b> clunky. Also, Second Life has always been a neat experiment, and I hope it still is and continues to keep doing well 🙂.</p>
import java.lang.reflect.Field;

Field[] FieldList=R.drawable.class.getDeclaredFields(); //replace "drawable"
String[] Names=new String[FieldList.length];
int[] IDs=new int[FieldList.length];
for(int i=0;i<FieldList.length;i++)
	IDs[i]=getResources().getIdentifier(Names[i]=FieldList[i].getName(), "drawable", getClass().getPackage().getName()); //replace "drawable"
<div>
Original Post: <a href="https://www.castledragmire.com/Posts/Deep_object_compare_for_javascript">Deep object compare for javascript</a><br>
Date: 2016-09-24<br>
Tags: Comparison, Examples, JavaScript<br><br>
</div>

```js
function DeepObjectCompare(O1, O2)
{
	try {
		DOC_Val(O1, O2, ['O1->O2', O1, O2]);
		return DOC_Val(O2, O1, ['O2->O1', O1, O2]);
	} catch(e) {
		console.log(e.Chain);
		throw(e);
	}
}
function DOC_Error(Reason, Chain, Val1, Val2)
{
	this.Reason=Reason;
	this.Chain=Chain;
	this.Val1=Val1;
	this.Val2=Val2;
}

function DOC_Val(Val1, Val2, Chain)
{
	function DoThrow(Reason, NewChain) { throw(new DOC_Error(Reason, NewChain!==undefined ? NewChain : Chain, Val1, Val2)); }

	if(typeof(Val1)!==typeof(Val2))
		return DoThrow('Type Mismatch');
	if(Val1===null || Val1===undefined)
		return Val1!==Val2 ? DoThrow('Null/undefined mismatch') : true;
	if(Val1.constructor!==Val2.constructor)
		return DoThrow('Constructor mismatch');
	switch(typeof(Val1))
	{
		case 'object':
			for(var m in Val1)
			{
				if(!Val1.hasOwnProperty(m))
					continue;
				var CurChain=Chain.concat([m]);
				if(!Val2.hasOwnProperty(m))
					return DoThrow('Val2 missing property', CurChain);
				DOC_Val(Val1[m], Val2[m], CurChain);
			}
			return true;
		case 'number':
			if(Number.isNaN(Val1))
				return !Number.isNaN(Val2) ? DoThrow('NaN mismatch') : true;
		case 'string':
		case 'boolean':
			return Val1!==Val2 ? DoThrow('Value mismatch') : true;
		case 'function':
			if(Val1.prototype!==Val2.prototype)
				return DoThrow('Prototype mismatch');
			if(Val1!==Val2)
				return DoThrow('Function mismatch');
			return true;
		default:
			return DoThrow('Val1 is unknown type');
	}
}
```
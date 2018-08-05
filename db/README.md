# DB

Below describes all the tables and their cols

## pumps
Stores all the info about the connected pumps, this is created for you by running the provided schema file
 * pump_id (int) - integer id of the pump, I prefer to number them in the order that users will see them from left to right
 * flow_rate (real) - the amount of fluid the pump moves in ml / min
 * ingredient (text) - what is hooked up to this pump
 * gpio_pin (text) - the pin this pump (or relay) is connected to in the form 'GPIO_X'

## ingredients
All the various alcohols, or whatever else you want to hook up to this go here
 * ingredient (text) - the name of whatever it is that is hooked up
 * viscosity (real)- unused at this moment, this was meant to allow users to adjust the flow rate when pumping this liquid.  In theory how this will be used is if you pump something that only gets 80ml/min but the pump would do 100, you set this to .8 and then we mutliply the flow rate by it before pumping.  In some very minor tests I did it didn't seem to really change so I haven't used it at all, your mileage may vary.

## drinks
 * drink_id (int, pk) - numeric id of the drink, not exposed directly to users
 * name (text) - whatever you call this concoction
 * notes (text) - this is for when the machine is done making the drink, we will show the user notes of anything that needs done to finish the drink, garnish, whatever you want.  Don't add extra ingredients here (well you can if you want) but you can do that via normal drink ingredients by unchecking the dispense checkbox

## drink_ingredients
 * ingredient (text) - which ingredient we are adding
 * amount (real) - how much of the ingredient
 * units (text) - what the amount represents, ml, cl, oz, dash
 * dispense (bool) - are we going to pump this or not, an example of something that we would not dispense are bitters
 * drink_id (int, fk -> drinks) - links this to a drink

# Parts and Diagrams

In this directory are two Fritzing diagrams and their exported PNG formats.  DrinkMachine.fzz and DrinkMachine.png are the current / supported diagrams.

Note on the _AllGPIO versions.  This was my initial go at buildin this using a GPIO LCD I had lying around.  This isn't supported in the code at all, but I've included it for historical sake, and maybe one day I'll add support for it.

The LCD is very lightly used today, to where you might ask "why is this here?".  I had plans to do more but as I wrote it all I just kind of forgot about the LCD.  I stil intend to go back and make more use of it, until then the answer to the question is "eh, why not?".

## Parts List

Everything I used, and I tried to include links to either where I got them or at least to show generally what it is.  You might be able to find deals elsewhere, when I bought a lot of this it was whatever was cheapest.

```
Raspberry Pi Zero W (with headers) - https://www.raspberrypi.org/products/raspberry-pi-zero-w/
SainSmart 8ch relay - https://www.sainsmart.com/products/8-channel-5v-relay-module
SainSmart I2C LCD - https://www.sainsmart.com/products/16x2-iic-i2c-twi-serial-lcd-display-blue-on-white
8 x 12v Peristaltic pump - https://www.amazon.com/gp/product/B01IUVHB8E
12v 2A power supply - https://www.amazon.com/gp/product/B00Q2E5IXW
12v -> 5v step down - https://www.amazon.com/gp/product/B00CXKBJI2

Misc Parts:
2.1 x 5.5mm female dc plug - https://www.amazon.com/gp/product/B00CUKHN0S
Micro USB connector - https://www.amazon.com/gp/product/B01KV4RANE
Butt connectors
Female crimp disconnect (pwr / gnd to pumps)
Heat shrink tubing
Various wiring, 22awg and I used some thermostat wiring for the pumps since it was cheap at your local home improvement store
ideal in-sure push-in wire connectors (seemed easiest to tie the pumps to the power) - http://www.idealindustries.ca/products/prodSelect.php?prodId=30-1034
```

## Diagrams

Just a quick note on the diagrams.  I'm new to fritzing and a lot of this hardware work, so I tried to find exact parts as I could.  You'll notice some notes where I could not find the parts like the 12v power and the step-down.  If anyone knows where I can find a parts file for these that would be great!

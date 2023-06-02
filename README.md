### deskpi-ctl - A simple DeskPi Pro tool
#### Yawning Angel <yawning@schwanenlied.me>

As part of the on-going futile search for an ARM development target that
doesn't suck, I ended up using a Raspi 4 (8 GiB) in a [DeskPi Pro][1]
enclosure.  While the way the SATA support is connected is kind of
jank (a loopback USB connector on the back), at least it is a sort-of
ok AArch64 development target that has an actual m.2 SATA SSD, that
even supports TRIM.

Naturally I would love [something][2] [better][3], but that would cost
money.

The DeskPi creators did release the source code to the scripts that
allow you to control the fan and to "safely" power the unit off, but
I feel compelled to rewrite the tooling.

The actual control scheme is dead trivial.  There is a "QinHeng
Electronics CH340 serial converter" (`1a86:7523`), that comes up as
`/dev/ttyUSB0`, that appears to be connected to the fan and power
control.  Writing various plain text commands (eg: even just with
`echo`), does useful things.

- `pwm_000` ...  `pwm_100`: Set fan speed by percent, needs to be 3 digits.
- `power_off`: Turn off the unit "safely".

##### Annoyances

- The fan probably just uses low-frequency PWM, and the only duty
cycles that don't make a lot of really annoying noises are 0% and
100%.  I also did not bother checking the extension board to see
if there is a flyback diode present.
- The vendor provided C code's serial port initialization is "odd".
Thankfully, 9600 8-N-1 "just works".
- Recent Ubuntu (that I don't use, because Debian/Ubuntu is banned on
my hardware), has a Braille tty package that claims ownership of
the USB serial device that controls the fan and power.  Remove
`brltty` to get things to work.

[1]: https://deskpi.com/products/deskpi-pro-for-raspberry-pi-4
[2]: http://www.orangepi.org/html/hardWare/computerAndMicrocontrollers/details/Orange-Pi-5-plus.html
[3]: https://www.solid-run.com/arm-servers-networking-platforms/honeycomb-servers-workstation/#honeycomb-lx2-workstation
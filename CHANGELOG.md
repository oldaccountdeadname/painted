# v0.1.0
or: everything implemented is known-working

Painted is a lightweight notification daemon which will eventually
contain a complete feature set. As it currently stands, painted is
incomplete. However, it's incomplete yet stable in that incompleteness.

If you want to write scripts for Painted, you can do so now, with
*enough* assurance that little will break in the future. Features will
mostly only be added before a 1.0.0. Note, however, that we're still in
0.x.x, and you should heed SemVar's advice: anything may change at any
time.

The biggest notable missing feature is the lack of format strings, so
notifications are all pretty ugly (see the following sample):

```
&{OriginApp:notify-send Summary:painted test Id:1}
```

But while output is subject to change, existing input is pretty stable
at the moment. Thus, 0.1.0.

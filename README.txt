== AKAI LPD8 command line tool =================================================

This is a tool for backing up and configuring the AKAI LPD8 MIDI controller.
It is a reimplementation of a very similar tool written in Python.

    https://github.com/boomlinde/lpd8

I decided to reimplement because the Python tool only works on Linux.


== Building & Running ==========================================================

Install a recent Go version, then run:

   go build .

To back up your LPD8, run:

   ./lpd8 backup backup.json

To restore a complete backup:

   ./lpd8 restore backup.json

You can, of course, also make adjustments to backup.json before restoring.

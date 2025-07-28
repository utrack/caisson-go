# levels

This directory splits the platform into 'levels' - arbitrarily separating the lower-level/more generic packages from the high-level ones.

Lower levels are intended to be used by the specific platform's developers only; they should ban the end-users/application developers from importing them via the `forbidigo` or `depguard` linters.

You can use the higher-level packages without using the lower levels, however, it may not be possible to swap the lower-level package while still using the higher level.

A level is forbidden from using the higher-level packages; it only can use same level or lower-level packages.
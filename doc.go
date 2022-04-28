// Copyright 2020 Penguin Statistics. All rights reserved.
// Use of this source code is governed by MIT license that
// can be found in the LICENSE file.

/*
widget-backend renders minimum matrix data and their corresponding data dependencies
with the widget frontend all at once in order to minimize data consumption and to
reduce load time for users of the widget.
Usage:
	go run .

Notice that default values have been set for all data controllers and their caches,
and the one SHOULD be checking them before deploying such service. Please also notice
that, such service which probably has already been deployed, is NOT suitable for public
usage except for the render of widgets. Developers SHOULD use the Penguin Statistics
public API which is documented at https://developer.penguin-stats.io/docs/ instead of hosting
this service themselves in order to reduce server load.
*/
package main

SPLASH
================

Downlaods images from [Splashbase](http://www.splashbase.co/)!

Images are stored in `$HOME/.splash_api`. Run this periodically to download the latest images. Set your OS to randomly retrieve images from `$HOME/.splash_api` and you're good to go.

### Usage Instructions

    go get github.com/drewwells/splash_api/cmd/splash
    splash

To check for new images, but not download them.

    splash -check

Don't have Go? Choose one of the precompiled binaries on the [Release](https://github.com/drewwells/splash_api/releases) page.

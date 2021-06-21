# Grasp
Grasp is a tiny command line utility to generate MD5 digest for website
addresses.

## Overview
I tried to be simple, but correct, using only the standard library, as requested
by the challenge description. The entire code is in a single file, called `grasp.go`.
There are three functions, besides main:
* `normalizeAddresses` to add *http* to the addresses whenever it is necessary
* `fetcher` the worker itself, it fetches the body and publishes the MD5 digest
* `run` to isolate the code execution from main

### Workers, channels, and synchronization
There are two types of workers here
* `fetcher` to retrieve the page and calculate the MD5
* `printer` to print out the addresses and their MD5 values

The communication happens using two channels
* `inputChannel` receives the normalized URLs passed as arguments from command line
* `printChannel` receives the digest version of the response body

The synchronization happens using two wait groups
* `wgFetcher` ensures that only after finishing fetching we may close the print channel
* `wgPrinter` ensures that we wait for all the messages to be print before moving on with the execution

### Error messages
For simplicity, `grasp` does not use another channel to send out error messages,
every error is sent out using the same `print` channel.

## Do not add any extra features
If I were to add a couple of flags, I would definitely add one to control the
length of the input channel buffer, and another one to control the length of the
print channel buffer. Also, having a way to control the amount of printer workers
would be interesting. This way you would have all the buttons and knobs to change
and see what variables play the most important roles here.

## Running
In order to run this challenge, please, move to a temporary directory,
clone this repository, build and run it. The process may change from
environment to environment, but it should be something like this:

```bash
$ go build grasp.go
$ ./grasp yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny
$ ./grasp -parallel 12 adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com
```

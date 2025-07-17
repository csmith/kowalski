# Changelog

## 6.0.2 - 2025-07-17

_No code changes, just build process fixes._

## 6.0.1 - 2025-07-17

_No code changes, just build process fixes._

## 6.0.0 - 2025-07-17

### Breaking changes

* Discord bot moved from `cmd/kowalski` to `cmd/discord`
* Discord bot docker images are now named `kowalski/discord`

### Features

* Added a web UI

### Other changes

* Fixed `wordlink` not lowercasing its inputs properly
* Updated docker images to use `greboid/dockerbase`

## 5.4.0 - 2025-06-05

### Changes

* Dependency updates

## 5.3.0 - 2022-03-14

### Changes

* New algorithm for scoring "English-like" words
* Analysis now shows the score for "might be English" entries
* Fix Discord image commands not working with plain URLs

## 5.2.0 - 2022-03-14

### Changes

* FST functions are now in their own package instead of tied to Discord

## 5.1.0 - 2022-03-13

### Changes

* Fix caesar shift analysis being off-by-one still.
* Add Discord commands for using FSTs to find words by regexp or morse

## 5.0.0 - 2021-04-27

### Breaking changes

* Some generic functions have moved to the [csmith/cryptography](https://github.com/csmith/cryptography) library

### Changes

* Fix the RGB extraction swapping green and blue channels
* File commands on Discord will accept a URL if attachments aren't available
* Add function for finding hidden pixels in images
* Analysis now shows how long the input is

## 4.0.2 - 2021-04-02

### Changes

* Fix panic in analysis in some cases

## 4.0.1 - 2021-04-01

### Changes

* Discord bot now designates its messages as replies
* Analysis notes if input only contains a small subset of letters,
  with a specific hint for ADFG(V)X ciphers.
* Analysis now checks for run-length-encoding

## 4.0.0 - 2021-03-29

### Breaking changes

* Long-running functions now take a Context param so they can be cancelled 

### Changes

* Add chunk command to divide text into blocks
* Add transpose command to swap the axes of a grid
* Add help command to Discord bot
* Fix caesar shift output in analysis being off-by-one
* Add function for counting colours in an image
* Add function to split an image into RGB channels

## 3.1.0 - 2021-03-24

### Changes

* Analysis now highlights if text contains most English letters but misses a few
* Analysis now shows the split into symbols if its input is made up of chemical elements
  (previously it was just noted that it did, and dividing it up was left as an exercise for the reader.)
* Add letter distribution function

## 3.0.0 - 2020-11-15

### Breaking changes

* The `Analysis()` function now requires a SpellChecker parameter.

### Changes

* Text analysis now notes if caesar shift might work

## 2.3.0 - 2020-11-15

### Changes

* Text analysis now notes if characters all come from one or two rows of a QWERTY keyboard
* Support for Caesar shifting
* Support for finding words as substrings, or within a word search grid

## 2.2.0 - 2020-10-21

### Changes

* Add generic text analysis function

## 2.1.0 - 2020-10-21

### Changes

* Add support for decoding T9

## 2.0.0 - 2020-10-08

### Breaking changes

* All functions that need to check if a word exists now take
  a SpellCheck parameter.

### Changes

* Add `Multiplex*` functions to allow using multiple dictionaries
  in parallel.

## 1.4.0 - 2020-10-07

### Changes

* Add `FromMorse()` for finding words from undelimited morse code

## 1.3.0 - 2020-10-07

### Changes

* Switch to using bloom filters for checking dictionary words.
  This reduces memory usage for a ~99MB word list from ~2GB to ~100MB.

## 1.2.0 - 2020-04-13

### Changes

* `Match()` and `Anagrams()` now lowercase their inputs automatically

## 1.1.0 - 2020-04-13

### Changes

* Faster anagram implementation

## 1.0.0 - 2020-04-11

_Initial release._

# Kowalski

![Kowalski, analysis](kowalski.jpg)

Kowalski is a Discord bot that can perform rudimentary analysis on text to help solve puzzles.

It currently supports these commands:

- `match <term>` returns all known words that match the given term,
  where '?' is a single-character wildcard.
  e.g. `match melism?` will return `melisma`.

- `anagram <term>` returns all known anagrams that match the given term,
  where '?' is a single-character wildcard.
  e.g. `anagram lismem?` will return `melisma`.

- `morse <term>` returns all possible words that match the given morse code
  (specified using `-` and `.`), ignoring all spaces/pauses.

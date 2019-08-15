# ipic

A quick ’n’ dirty rewrite of [Dr. Drang’s ipic][orig] in Go because I was too lazy to figure out how to install [Requests][] for a single-file script.

My version improves on the original slightly: it asks for 4096×4096 images, not measly 600×600 images. Might as well ask for an image big enough to not need upscaling, right?

## Usage

```
Usage: ipic (-i | -m | -a | -f | -t | -b | -h) SEARCH_TERM

Generates a page of thumbnails and links to larger images for items in the iTunes/App/macOS App Stores.

Options:
  -a	album
  -b	book
  -f	film
  -h	show this help message
  -i	iOS app
  -m	macOS app
  -t	TV show

Only one option is allowed. The HTML file for the generated webpage is saved to ~/Desktop.
```

## Bugs

- “It works on my machine.”
- Probably missing a `defer foo.Close()` somewhere.

## License

[Apache License, version 2](https://www.apache.org/licenses/LICENSE-2.0).



[orig]: https://github.com/drdrang/ipic
[requests]: https://github.com/kennethreitz/requests

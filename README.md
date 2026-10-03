# ImgOps

This is a CLI tool for reverse searching images through ImgOps website.
It supports files and URLs.

## How It Works

1. You upload a file or pass a URL to the tool
2. Depending on your choice, it will open search results from sites you want.  
  These are Google, TinEye, Yandex, Bing, Reddit and few others  
  If you don't provide a “target”, it will open ImgOps page with the provided image.

## Commands

| Command | Description |
| --- | --- |
| `imgops search <file-or-url>` | Search using `--targets`, alias `a` |
| `imgops iqdb <file-or-url>` | Reverse search on [IQDB](https://iqdb.org) (anime) |
| `imgops saucenao <file-or-url>` | Reverse search on [SauceNAO](https://saucenao.com) (anime) |
| `imgops tracemoe <file-or-url>` | Reverse search on [trace.moe](https://trace.moe) (anime, alias `trace`) |
| `imgops ascii2d <file-or-url>` | Reverse search on [ascii2d](https://ascii2d.net) (anime) |
| `imgops copy <file-or-url>` | Upload an image and copy its link to the clipboard (aliases `clip`, `url`) |

Every command accepts `--return`/`-r` to print the resulting URL instead of opening it in the browser.

The new anime engines are also available as regular targets, so you can mix them:

```sh
imgops search -t "google, iqdb, saucenao" -r image.png
imgops search -i -r image.png   # type initials, e.g. "qm" for iqdb + trace.moe
```

## Download

[Click here](https://github.com/dogancelik/imgops/releases/latest) to download latest version.

## How To Use

See `imgops -h` or check out [the Wiki](https://github.com/dogancelik/imgops/wiki/Examples).

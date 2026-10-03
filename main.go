package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/pkg/browser"
	"github.com/urfave/cli"
)

import . "github.com/visionmedia/go-debug"

var debug = Debug("imgops")

var authors = []cli.Author{
	{
		Name:  "Doğan Çelik",
		Email: "dogancelik.com",
	},
}

var Version string

var returnFlag = cli.BoolFlag{
	Name:  "return, r",
	Usage: "Output the result URL instead of opening it in the browser",
}

var copyFlag = cli.BoolFlag{
	Name:  "copy, c",
	Usage: "Copy the uploaded image link (ImgOps image URL) to the clipboard",
}

func cliSelect() string {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	fmt.Println(genSelectText(false))
	ret := ""

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		} else if key == keyboard.KeyEsc {
			break
		}

		m := getKeyToNameTargets(availableTargets)
		target, mapOk := m[char]
		if mapOk {
			return target
		} else if char == 'i' {
			return defaultTarget
		}
	}

	return ret
}

// checkSource makes sure the given path/URL can be used.
func checkSource(srcPath string) error {
	if isUrl(srcPath) {
		return nil
	}

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		if isCommandName(srcPath) {
			return cli.NewExitError("'"+srcPath+"' is a command, not a file.\nUse 'imgops "+srcPath+" <file-or-url>' or pass the flag '-c' to copy the image link while searching, e.g. 'imgops iqdb -c image.png'", 2)
		}
		return cli.NewExitError("File doesn't exist: "+srcPath, 2)
	}

	return nil
}

// isCommandName reports whether the given string is one of the CLI command
// names or aliases.
func isCommandName(s string) bool {
	for _, name := range []string{"search", "a", "iqdb", "saucenao", "tracemoe", "trace", "ascii2d", "copy", "clip", "url", "help", "h"} {
		if s == name {
			return true
		}
	}
	return false
}

// resolveImageLink returns the link of the image as stored by ImgOps.
// It prefers the link found in the last response body and falls back to the
// final redirect URL.
func resolveImageLink(srcPath string) string {
	if isUrl(srcPath) {
		return srcPath
	}

	if url, err := getUploadedImageURL(lastBody); err == nil && strings.TrimSpace(url) != "" {
		return url
	}

	return imageUrlFromFinal(finalUrl)
}

// copyImageLink copies the uploaded image link to the clipboard and prints it.
func copyImageLink(srcPath string) error {
	imageLink := resolveImageLink(srcPath)
	if strings.TrimSpace(imageLink) == "" {
		return cli.NewExitError("Could not get the image link", 8)
	}

	if err := ClipboardCopy(imageLink); err != nil {
		return cli.NewExitError("Could not copy to the clipboard: "+err.Error(), 7)
	}

	fmt.Println("Copied to clipboard: " + imageLink)

	return nil
}

// runSearch uploads the given file/URL and returns the search result links
// for the requested targets.
func runSearch(srcPath, targets string) ([]string, error) {
	debug("Path: %s", srcPath)
	debug("Targets: %s", targets)

	if err := checkSource(srcPath); err != nil {
		return nil, err
	}

	if isUrl(srcPath) {
		debug("Start URL upload")
		return UploadURL(srcPath, targets)
	}

	debug("Start file upload")
	return UploadFile(srcPath, targets)
}

// openOrPrint opens every URL in the browser, or prints them when the
// --return flag is used.
func openOrPrint(urls []string, returnOnly bool) {
	for _, url := range urls {
		if returnOnly {
			fmt.Println(url)
		} else {
			browser.OpenURL(url)
		}
	}
}

func cliSearch(c *cli.Context) error {

	if c.NArg() == 0 {
		return cli.NewExitError("No file or URL is given", 1)
	}

	srcPath := c.Args().First()
	targets := c.String("targets")

	// Select flag
	if c.Bool("select") {
		targets = cliSelect()
		if targets == "" {
			return cli.NewExitError("Upload cancelled", 4)
		}
	}

	// Input flag
	if c.Bool("input") {
		fmt.Print(genSelectText(true))
		fmt.Scanln(&targets)
		targets = initialsToTargets(targets)
		if targets == "" {
			return cli.NewExitError("No target is specified", 5)
		}
	}

	urls, errUpload := runSearch(srcPath, targets)

	// Upload result
	if errUpload != nil && len(urls) == 0 {
		return cli.NewExitError("Error during upload: "+errUpload.Error(), 3)
	}

	if errUpload != nil {
		fmt.Fprintf(os.Stderr, "Unknown targets '%s', will open default page instead", targets)
	}

	openOrPrint(urls, c.Bool("return"))

	if c.Bool("copy") {
		if err := copyImageLink(srcPath); err != nil {
			return err
		}
	}

	return nil
}

// cliSite returns an action that searches a fixed target website.
func cliSite(target string) cli.ActionFunc {
	return func(c *cli.Context) error {
		if c.NArg() == 0 {
			return cli.NewExitError("No file or URL is given", 1)
		}

		srcPath := c.Args().First()

		urls, errUpload := runSearch(srcPath, target)
		if errUpload != nil && len(urls) == 0 {
			return cli.NewExitError("Error during upload: "+errUpload.Error(), 3)
		}

		openOrPrint(urls, c.Bool("return"))

		if c.Bool("copy") {
			if err := copyImageLink(srcPath); err != nil {
				return err
			}
		}

		return nil
	}
}

// cliClipboard uploads the given file/URL and copies the resulting image
// link to the system clipboard.
func cliClipboard(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.NewExitError("No file or URL is given", 1)
	}

	srcPath := c.Args().First()

	if err := checkSource(srcPath); err != nil {
		return err
	}

	imageURL, err := ImageURL(srcPath)
	if err != nil {
		return cli.NewExitError("Error while getting the image URL: "+err.Error(), 6)
	}

	if err := ClipboardCopy(imageURL); err != nil {
		return cli.NewExitError("Could not copy to the clipboard: "+err.Error(), 7)
	}

	fmt.Println(imageURL)

	return nil
}

func cliMain(c *cli.Context) error {
	cli.ShowAppHelp(c)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "ImgOps CLI"
	app.Usage = "Reverse search images"
	app.Version = Version
	app.Authors = authors
	app.Action = cliMain
	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"a"},
			Usage:   "Search a file or a URL",
			Action:  cliSearch,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "targets, t",
					Value: defaultTarget,
					Usage: "Target website to search at (e.g. Google)",
				},
				cli.BoolFlag{
					Name:  "select, s",
					Usage: "Show a list of targets to select from",
				},
				cli.BoolFlag{
					Name:  "input, i",
					Usage: "Type the targets you want to open",
				},
				returnFlag,
				copyFlag,
			},
		},
		{
			Name:   "iqdb",
			Usage:  "Reverse search an image on IQDB (anime)",
			Action: cliSite("iqdb"),
			Flags:  []cli.Flag{returnFlag, copyFlag},
		},
		{
			Name:   "saucenao",
			Usage:  "Reverse search an image on SauceNAO (anime)",
			Action: cliSite("saucenao"),
			Flags:  []cli.Flag{returnFlag, copyFlag},
		},
		{
			Name:    "tracemoe",
			Aliases: []string{"trace"},
			Usage:   "Reverse search an image on trace.moe (anime)",
			Action:  cliSite("tracemoe"),
			Flags:   []cli.Flag{returnFlag, copyFlag},
		},
		{
			Name:   "ascii2d",
			Usage:  "Reverse search an image on ascii2d (anime)",
			Action: cliSite("ascii2d"),
			Flags:  []cli.Flag{returnFlag, copyFlag},
		},
		{
			Name:    "copy",
			Aliases: []string{"clip", "url"},
			Usage:   "Upload an image and copy its link to the clipboard",
			Action:  cliClipboard,
		},
	}
	app.Run(os.Args)
}

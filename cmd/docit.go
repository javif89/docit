package main

import (
	"fmt"
	"io"
	"javif89/docit/internal/pagebuilder"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name: "docit",
		Usage: "A simple static site generator for documentation",
		Authors: []*cli.Author{
			{
				Name: "Javier Feliz",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "build",
				Aliases:   []string{"b"},
				Usage:     "Build your site to the build folder",
				UsageText: "docit build -i [input directory] -o [output directory] -t [site title]",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "input", Value: ".", Usage: "Path to the `input folder`", Aliases: []string{"i"}},
					&cli.StringFlag{Name: "output", Value: "./build", Usage: "Path to the `output folder`", Aliases: []string{"o"}},
					&cli.StringFlag{Name: "title", Value: "docs", Usage: "Project title", Aliases: []string{"t"}},
				},
				Action: func(c *cli.Context) error {
					build(c.String("input"), c.String("output"), c.String("title"))
					return nil
				},
			},
			{
				Name:      "serve",
				Aliases:   []string{"s"},
				Usage:     "start the development server",
				UsageText: "docit serve -h [host (default localhost)] -p [port (default 8000)] -s [build folder (default ./build)]",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Value: "localhost", Usage: "Host to serve the site"},
					&cli.StringFlag{Name: "port", Value: "8000", Usage: "Port to serve the site", Aliases: []string{"p"}},
					&cli.StringFlag{Name: "site", Value: "./build", Usage: "Path to the `build folder`", Aliases: []string{"s"}},
					&cli.BoolFlag{Name: "watch", Value: false, Usage: "Watch for changes and rebuild the site"},
					&cli.BoolFlag{Name: "build", Value: false, Usage: "Build the site before serving"},
					&cli.StringFlag{Name: "input", Value: ".", Usage: "Path to the `input folder` (default ./docs) if building", Aliases: []string{"i"}},
					&cli.StringFlag{Name: "title", Value: "docs", Usage: "Project title if building (default docs)", Aliases: []string{"t"}},
				},
				Action: func(c *cli.Context) error {
					if c.Bool("build") {
						build(c.String("input"), c.String("site"), c.String("title"))
					}

					http.Handle("/", http.FileServer(http.Dir(c.String("site"))))
					fmt.Println("Serving site...")
					address := c.String("host")+":"+c.String("port")
					fmt.Println(address)
					
					if c.Bool("watch") {
						go watch("../docs")
					}

					http.ListenAndServe(address, nil)

					return nil
				},
			},
		},
	}

	// Pre flight check
	if _, err := os.Stat("./page.html"); os.IsNotExist(err) {
		fmt.Println("Downloading default template...")
		downloadFile("https://raw.githubusercontent.com/javif89/docit-assets/main/dist/page-v1.html", "./page.html")
	}

	// Run the app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func downloadFile(url string, filename string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error downloading file", err)
	}

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}

func build(input string, output string, sitename string) {
	start := time.Now().UnixMilli()
	b := pagebuilder.NewBuilder("./page.html", input, output)
	b.ProjectTitle = sitename
	b.Build()
	end := time.Now().UnixMilli()
	fmt.Println("Built in", (end-start),"ms")
}

// Watch a directory for file changes. We'll use this to rebuild the site when
// needed
func watch(dir string) {
	fmt.Println("Watching", dir, "for changes")
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		// since fsnotify can watch all the files in a directory, watchers only need
		// to be added to each nested directory
		if fi.Mode().IsDir() {
			return watcher.Add(path)
		}

		return nil
	}); err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Println("Changed", event.Name)
				fmt.Println("Rebuilding site...")
				build("../docs", "../build", "docs")

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}
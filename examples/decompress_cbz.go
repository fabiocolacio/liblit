package main

import(
    "github.com/fabiocolacio/liblit/cbz"
    "io/ioutil"
    "flag"
    "fmt"
)

/*
 * This program demonstrates how to parse a cbz file.
 * The program writes each page (an image) of the comic book archive
 * to a file in the current directory.
 */
func main() {
    var infile string
    flag.StringVar(&infile, "f", "", "The cbz file to decompress")
    flag.Parse()

    if infile == "" {
        fmt.Println("Please specify a file to open with the '-f' flag.")
        return
    }

    pages, err := cbz.NewFromFile(infile)
    if err != nil {
        fmt.Println(err)
        return
    }

    for i, page := range pages {
        filename := fmt.Sprintf("Page%d.jpg", i)
        ioutil.WriteFile(filename, page, 0666)
    }
}


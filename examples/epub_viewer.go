package main

import(
    "github.com/fabiocolacio/liblit/epub"
    "golang.org/x/net/html"
    "flag"
    "fmt"
)

/*
 * This simple example demonstrates how to open an epub file,
 * extract some metadata from it, and parsing the html contents
 * so that the book can be read as plain-text in the terminal.
 *
 * For the best viewing experience, you can pipe the output of this
 * program to pager like less, or pipe it to a file to be opened with
 * another program.
 */
func main() {
    var infile string
    flag.StringVar(&infile, "f", "", "The epub file to open")
    flag.Parse()

    if infile == "" {
        fmt.Println("Please specify an epub file to open with the '-f' flag.")
        return
    }

    book, err := epub.NewFromFile(infile)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("Title:", book.Metadata.Title)
    fmt.Println("Author:", book.Metadata.Author)
    fmt.Println("Subjects:")
    for _, subj := range book.Metadata.Subjects {
        fmt.Println("*", subj)
    }

    fmt.Println("Contents:")
    for _, token := range book.Tokens {
        if token.Type == html.StartTagToken {
            if token.Data == "p" {
                fmt.Print("\t")
            }
        }

        if token.Type == html.SelfClosingTagToken {
            if token.Data == "b" {
                fmt.Print("\n")
            }
        }

        if token.Type == html.TextToken {
            fmt.Print(token.Data)
        }
    }
}


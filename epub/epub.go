package epub

import(
    "golang.org/x/net/html"
    "archive/zip"
    "errors"
    "encoding/xml"
    "io/ioutil"
    "strings"
    "fmt"
)

var(
    ErrInvalidEpub error = errors.New("Invalid Epub file")
)

type Epub struct {
    Metadata   Metadata
    Tokens   []html.Token
}

type rootFile struct {
    FullPath  string `xml:"full-path,attr"`
    MediaType string `xml:"media-type,attr"`
}

type rootFiles struct {
    RootFiles []rootFile `xml:"rootfile"`
}

type container struct {
    RootFiles rootFiles `xml:"rootfiles"`
}

type Metadata struct {
    Title          string `xml:"http://purl.org/dc/elements/1.1/ title"`
    Author         string `xml:"http://purl.org/dc/elements/1.1/ creator"`
    Contributors []string `xml:"http://purl.org/dc/elements/1.1/ contributor"`
    Subjects     []string `xml:"http://purl.org/dc/elements/1.1/ subject"`
    Language       string `xml:"http://purl.org/dc/elements/1.1/ language"`
    Date           string `xml:"http://purl.org/dc/elements/1.1/ date"`
    Source         string `xml:"http://purl.org/dc/elements/1.1/ source"`
}

type item struct {
    Id        string `xml:"id,attr"`
    MediaType string `xml:"media-type,attr"`
    Href      string `xml:"href,attr"`
}

type manifest struct {
    Items []item `xml:"item"`
}

type opf struct {
    XMLName  xml.Name `xml:"package"`
    Metadata Metadata `xml:"metadata"`
    Manifest manifest `xml:"manifest"`
}

func NewFromFile(filename string) (*Epub, error) {
    epub := new(Epub)

    reader, err := zip.OpenReader(filename)
    if err != nil {
        return nil, err
    }
    defer reader.Close()

    // Create mapping of all files in the archive
    zipFiles := make(map[string]*zip.File)
    for _, file := range reader.File {
        zipFiles[file.Name] = file 
    }

    // Check the mimetype
    if mimeFile := zipFiles["mimetype"]; mimeFile == nil {
        return nil, ErrInvalidEpub
    } else {
        mimeReader, err := mimeFile.Open()
        if err != nil {
            return nil, err
        }

        contents, err := ioutil.ReadAll(mimeReader)
        if err != nil {
            return nil, err
        }

        targetMime := "application/epub+zip"
        if string(contents) != targetMime {
            return nil, ErrInvalidEpub }

        mimeReader.Close()
    }

    // Parse container.xml
    var opfPath string
    if containerFile := zipFiles["META-INF/container.xml"]; containerFile == nil {
        return nil, ErrInvalidEpub
    } else {
        containerReader, err := containerFile.Open()
        if err != nil {
            return nil, err
        }

        contents, err := ioutil.ReadAll(containerReader)
        if err != nil {
            return nil, err
        }

        containerReader.Close()

        var cont container
        if err := xml.Unmarshal(contents, &cont); err != nil {
            return nil, err
        }

        if cont.RootFiles.RootFiles == nil {
            return nil, ErrInvalidEpub
        }

        opfPath = cont.RootFiles.RootFiles[0].FullPath
    }

    // Parse oebps file
    if opfFile := zipFiles[opfPath]; opfFile == nil {
        return nil, ErrInvalidEpub
    } else {
        opfReader, err := opfFile.Open()
        if err != nil {
            return nil, err
        }

        contents, err := ioutil.ReadAll(opfReader)
        if err != nil {
            return nil, err
        }

        opfReader.Close()
        
        var opfContents opf
        if err := xml.Unmarshal(contents, &opfContents); err != nil {
            return nil, err
        }

        epub.Metadata = opfContents.Metadata

        opfBase := opfPath[:strings.LastIndex(opfPath, "/") + 1]
        for _, item := range opfContents.Manifest.Items {
            if item.MediaType == "application/xhtml+xml" {
                if htmlFile := zipFiles[opfBase + item.Href]; htmlFile == nil {
                    fmt.Println(opfBase + item.Href)
                    return nil, ErrInvalidEpub
                } else {
                    htmlReader, err := htmlFile.Open()
                    if err != nil {
                        return nil, err
                    }

                    tokenizer := html.NewTokenizer(htmlReader)

                    for {
                        tokenType := tokenizer.Next()
                        
                        if tokenType == html.ErrorToken {
                            break
                        }

                        epub.Tokens = append(epub.Tokens, tokenizer.Token())
                    }

                    htmlReader.Close()
                }
            }
        }
    }

    return epub, nil
}


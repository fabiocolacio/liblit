package cbz

import(
    "archive/zip"
    "io/ioutil"
)

func NewFromFile(filename string) ([][]byte, error) {
    reader, err := zip.OpenReader(filename)
    if err != nil {
        return nil, err
    }
    defer reader.Close()

    var pages [][]byte

    for _, file := range reader.File {
        fileReader, err := file.Open()
        if err != nil {
            return nil, err
        }

        content, err := ioutil.ReadAll(fileReader)
        if err != nil {
            return nil, err
        }

        fileReader.Close()

        pages = append(pages, content)
    }

    return pages, nil
}


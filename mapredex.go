package main

import (
    "fmt"
    "os"
    "io"
    "bytes"
    "unicode"
//    "sync"
)


const gLetterType int = -1
const gWordType int = -2

type ResultMap map[string]int
type ResultMapChan chan ResultMap
type StringChan chan string


func Reducer(och ResultMapChan, ich StringChan, rType int) {
    rmap := make(ResultMap)
    
    defer func () {
        fmt.Println("\nReducer end!")
    }()
    
    fmt.Println("Reducer start!")
    
    rmap["__type"] = rType
    
    for r := range ich {
        _, exists := rmap[r]
        if (exists) {
            rmap[r] = rmap[r] + 1
            fmt.Printf("r")
        } else {
            rmap[r] = 1
            fmt.Printf("R")
        }
    }
    
    och <- rmap
}


func LetterReader(ch StringChan, fn string) {
    buff := make([]byte, 1024)
    
    fmt.Println("LetterReader start!")
    
    f, e := os.Open(fn)
    if (e != nil){
        panic(e)
    }
    
    // Close the file on exit
    defer func () {
        f.Close()
        fmt.Println("\nLetterReader end!")
    }()
    
    // file read loop
    n := 0
    e = nil
    for {
        n, e = f.Read(buff)
        if (n == 0) {
            break
        }
        
        rn := bytes.Runes(buff)
        for _, r := range rn {
            if unicode.IsLetter(r) {
                ch <- string(r)
                fmt.Printf(".")
            }
        }
    }
    
    if ((e != nil) && (e != io.EOF)) {
        panic(e);
    }
}


func HandleResults(resCh ResultMapChan) {
    msg := ""
    total := 0
    
    for resMap := range resCh {
        resType := resMap["__type"]
        delete (resMap, "__type")
        switch resType {
            case gLetterType:
                msg = "Letters"
            case gWordType:
                msg = "Words"
            default:
                msg = "Unknown"
        }
        fmt.Printf("Result %s:\n", msg);
        
        for k, v := range resMap {
            total += v
            fmt.Printf("    %s: %d\n", k, v)
        }
        
        fmt.Printf("    Total: %d\n", total)
    }
}


func main() {
    //var wg sync.WaitGroup
    
    fmt.Println("Main start!")
    
    args := os.Args
    resCh := make(ResultMapChan)
    letCh := make(StringChan, 1024)
    
    if len(args) > 1 {
        //wg.Add(2)
        go func () {
            defer func () {
                close(letCh)
                //wg.Done()
            }()
            
            LetterReader(letCh, args[1])
        }()
        
        go func () {
            defer func () {
                close(resCh)
                //wg.Done()
            }()
            
            Reducer(resCh, letCh, gLetterType)
        }()
        
        //wg.Wait()
        
        HandleResults(resCh)
    } else {
        fmt.Printf("Usage: %s <file-name>\n", args[0])
    }
    
    fmt.Println("Main end!")
}


package main

import (
    "fmt"
    BG "battleGrounds_server/source"
)

func main() {
    go BG.RunServer()

    var input string
    fmt.Scanln(&input)
}

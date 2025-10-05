package main

import (
    "typing_cli/colors"
    "strings"
    "fmt"
    "log"
    "math/rand"
    "github.com/mattn/go-tty"
    "time"
)

type printString struct{
    content string
    color string
}
type boardRow struct{
    segments []printString
    changed bool
}

var cursorLine int
var maxLine int
var board []boardRow

var lines = make(map[string]int)

func streakString(streak int) string {

    i:=0
    var str string
    str=""
    for i < streak{
        str+="\u2588"
        i+=1
    }
    str+=colors.Primary
    return str
}

func updateLetter(letter string, color string){
    var currLine = lines["question"]

    letterLines := strings.Split(letter, "\n")

    for i, q := range letterLines{
        if len(board) == currLine + i {
            newRow := boardRow{segments: []printString{}, changed: true}
            board = append(board, newRow)
        }
        newSeg := printString{content: q, color: color}
        board[currLine + i].segments = append(board[currLine + i].segments, newSeg)
    }
}

func updateStreak(streak int){
    var streakLine = lines["streak"]

    for len(board) <= streakLine{
        newRow := boardRow{segments: []printString{}, changed: true}
        board = append(board,newRow)
    }
    newSeg := printString{content: fmt.Sprintf("Streak: %d",streak), color: colors.Green}
    board[streakLine].segments = append(board[streakLine].segments, newSeg)
}

func refresh(){
    for i := 0; i < len(board); i++ {
        // if !board[i].changed { continue }

        // clear the current row
        fmt.Print("\r\033[2K")

        // print the next row from board
        for _, seg := range board[i].segments{
            fmt.Print(seg.color)
            fmt.Println(seg.content)
        }
    }

    // always move the cursor back to the top
    fmt.Print(fmt.Sprintf("\033[%dA",len(board) - 1))
}

func main() {
    cursorLine=0
    maxLine=0
    lines["question"] = 0
    lines["streak"] = 8

    rand.Seed(time.Now().UnixNano())

    var question int
    var streak int
    streak=0

    tty, err := tty.Open()
    if err != nil {
        log.Fatal(err)
    }
    defer tty.Close()

    // infinite loop
    for {
        question=rand.Intn(26)
        updateLetter(alph[question], colors.Green)
        updateStreak(streak)
        refresh()

        for{
            char, err := tty.ReadRune()
            if err != nil {
                log.Fatal(err)
            }

            // Correct answer
            if int(char)-97 == question {
                // time.Sleep(100*time.Millisecond)
                streak+=1
                break

            }else{
                // time.Sleep(50*time.Millisecond)
                streak=0
                updateStreak(streak)
                refresh()
            }
        }
    }
}

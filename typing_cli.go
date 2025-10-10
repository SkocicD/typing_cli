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

type screenTile struct{
    color string
    content rune
}
type box struct {
    r int
    c int
    height int
    width int
}

var termwidth = 181
var linecount = 10
// const boardsize = termwidth * linecount
const boardsize = 181*10
var board [boardsize]screenTile

var slot1 = box{0,0,10,14}
var slot2 = box{0,14,10,14}
var slot3 = box{0,28,10,14}
var pastBox = box{1,1,8,12}
var currBox = box{1,15,8,12}
var nextBox = box{1,29,8,12}

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


func update(){
    var prevcolor string
    var toPrint = ""
    for _, tile := range board{
        if tile.color != prevcolor{
            toPrint+=tile.color
            prevcolor = tile.color
        }
        toPrint += string(tile.content)
    }
    fmt.Println(toPrint)
    fmt.Print(fmt.Sprintf("\033[%dA", linecount))
}

// func updateLetter(letter string, color string){
//     letterLines := strings.Split(letter, "\n")
//
//     for i, q := range letterLines{
//         if len(board) == currLine + i {
//             newRow := boardRow{segments: []printString{}, changed: true}
//             board = append(board, newRow)
//         }
//         newSeg := printString{content: q, color: color}
//         board[currLine + i].segments = append(board[currLine + i].segments, newSeg)
//     }
// }
func set(bounds box, r int, c int, ch rune){
    board[termwidth * (bounds.r + r) + bounds.c + c].content = ch
}
func outlineBox(bounds box) {
    for c:=0; c < bounds.width; c++{
        set(bounds, 0, c, '-')
        set(bounds, bounds.height-1, c, '-')
    }
    for r:=0; r < bounds.height; r++{
        set(bounds, r, 0, '|')
        set(bounds, r, bounds.width-1, '|')
    }
}

func fillBox(bounds box, contents string, color string) int {
    letterLines := strings.Split(contents, "\n")
    var maxLen = 0
    for _, str := range letterLines{
        if len(str) > maxLen{
            maxLen = len(str)
        }
    }
    if len(letterLines) > bounds.height || maxLen > bounds.width{
        fmt.Println("the string passed in isnt the same size as the box")
        fmt.Println(fmt.Sprintf("The string has dimensions %dx%d", len(letterLines), maxLen))
        return -1
    }

    var ch rune
    for r, row := range letterLines{
        for c := 0; c < bounds.width; c++{
            if c >= len(row){
                ch = ' '
            } else {
                ch = rune(row[c])
            }
            board[termwidth * (bounds.r + r) + bounds.c + c].content = ch
            board[termwidth * (bounds.r + r) + bounds.c + c].color = color
        }
    }
    return maxLen
}


// func updateStreak(streak int){
//     var streakLine = lines["streak"]
//
//     for len(board) <= streakLine{
//         newRow := boardRow{segments: []printString{}, changed: true}
//         board = append(board,newRow)
//     }
//     newSeg := printString{content: fmt.Sprintf("Streak: %d",streak), color: colors.Green}
//     board[streakLine].segments = append(board[streakLine].segments, newSeg)
// }

// func refresh(){
//     for i := 0; i < len(board); i++ {
//         // if !board[i].changed { continue }
//
//         // clear the current row
//         fmt.Print("\r\033[2K")
//
//         // print the next row from board
//         for _, seg := range board[i].segments{
//             fmt.Print(seg.color)
//             fmt.Println(seg.content)
//         }
//     }
//
//     // always move the cursor back to the top
//     fmt.Print(fmt.Sprintf("\033[%dA",len(board) - 1))
// }

func main() {
    // init board
    for i := 0; i < len(board); i++{
        board[i].color = colors.Primary
        board[i].content = ' '
    }

    // outlineBox(slot1)
    // outlineBox(slot2)
    // outlineBox(slot3)
    update()

    rand.Seed(time.Now().UnixNano())
    var streak = 0

    tty, err := tty.Open()
    if err != nil {
        log.Fatal(err)
    }
    defer tty.Close()

    // infinite loop
    var thisQ int
    var nextQ = rand.Intn(26)
    var correct bool
    for {
        correct = true
        thisQ = nextQ
        nextQ = rand.Intn(26)
        // updateStreak(streak)

        fillBox(currBox, alph[thisQ], colors.Primary)
        fillBox(nextBox, alph[nextQ], colors.Primary)
        update()

        for{
            char, err := tty.ReadRune()
            if err != nil {
                log.Fatal(err)
            }

            // Correct answer
            if int(char)-97 == thisQ{
                // time.Sleep(100*time.Millisecond)
                streak+=1
                if correct{
                    fillBox(pastBox, alph[thisQ], colors.Green)
                }else{
                    fillBox(pastBox, alph[thisQ], colors.Red)
                }
                break

            }else{
                correct = false
                streak=0
                fillBox(currBox, alph[thisQ], colors.Red)
                update()
                // updateStreak(streak)
                // refresh()
            }
        }
    }
}

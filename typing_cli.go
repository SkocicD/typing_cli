package main

import (
    "strings"
    "fmt"
    "log"
    "math/rand"
    "github.com/mattn/go-tty"
    "time"
)

type colorstruct struct{
    green string
    red string
    end string
}

const alphlen int=26
var alph [alphlen]string
var colors colorstruct
var cursorLine int
var maxLine int

func declare_alph(){
    alph[0]=
`

  __ _
 / _  |
| (_| |
 \__,_|

`
    alph[1]=
`  _     
 | |    
 | |__  
 | '_ \ 
 | |_) |
 |_.__/

`
    alph[2]=
`

   ___ 
  / __|
 | (__ 
  \___|

`
    alph[3]=
`      _ 
     | |
   __| |
  / _  |
 | (_| |
  \__,_|

`
    alph[4]=
`

   ___ 
  / _ \
 |  __/
  \___|

`
    alph[5]=
`   __ 
  / _|
 | |_ 
 |  _|
 | |  
 |_|

`
    alph[6]=
`

   __ _ 
  / _  |
 | (_| |
  \__, |
   __/ |
  |___/`
    alph[7]=
`  _     
 | |    
 | |__  
 | '_ \ 
 | | | |
 |_| |_|

`
    alph[8]=
`  _ 
 (_)
  _ 
 | |
 | |
 |_|

`
    alph[9]=
`    _ 
   (_)
    _ 
   | |
   | |
   | |
  _/ |
 |__/ `
    alph[10]=
`  _    
 | |   
 | | __
 | |/ /
 |   < 
 |_|\_\

`
    alph[11]=
`  _ 
 | |
 | |
 | |
 | |
 |_|

`
    alph[12]=
`

  _ __ ___  
 | '_ ' _ \ 
 | | | | | |
 |_| |_| |_|

`
    alph[13]=
`

  _ __  
 | '_ \ 
 | | | |
 |_| |_|

`
    alph[14]=
`

   ___  
  / _ \ 
 | (_) |
  \___/ 

`
    alph[15]=
`

  _ __  
 | '_ \ 
 | |_) |
 | .__/ 
 | |    
 |_|`
    alph[16]=
`

   __ _ 
  / _  |
 | (_| |
  \__, |
     | |
     |_|`
    alph[17]=
`

  _ __ 
 | '__|
 | |   
 |_|

`
    alph[18]=
`

  ___ 
 / __|
 \__ \
 |___/

`
    alph[19]=
`  _   
 | |  
 | |_ 
 | __|
 | |_ 
  \__| 

`
    alph[20]=
`

  _   _ 
 | | | |
 | |_| |
  \__,_|

`
    alph[21]=
`

 __   __
 \ \ / /
  \ V / 
   \_/

`
    alph[22]=
`

 __      __
 \ \ /\ / /
  \ V  V / 
   \_/\_/

`
    alph[23]=
`

 __  __
 \ \/ /
  >  < 
 /_/\_\

`
    alph[24]=
`

  _   _ 
 | | | |
 | |_| |
  \__, |
   __/ |
  |___/`
    alph[25]=
`

  ____
 |_  /
  / / 
 /___|

`
}

func streakString(streak int) string {

    i:=0
    var str string
    str=""
    for i < streak{
        str+="\u2588"
        i+=1
    }
    str+=colors.end
    return str
}

func moveCursorUp(numLines int){
    if numLines>cursorLine{
        numLines=cursorLine
    }
    if numLines==0{
        return
    }
    var s string
    s=fmt.Sprintf("\033[%dA",numLines)
    fmt.Print(s)
    cursorLine-=numLines
}

func clearLines(numLines int){
    i:=0
    startingCursorLine:=cursorLine
    for i < numLines && i < maxLine{
        fmt.Print("\033[2K\r")
        fmt.Print("\033[1B")
        cursorLine+=1
        i+=1
    }
    moveCursorUp(cursorLine-startingCursorLine)
}

func printAndCountLines(str string, color string, newLine bool){
    if newLine{ str += "\n" }
    cursorLine+=strings.Count(str,"\n")
    fmt.Print(color)
    fmt.Print(str)
    fmt.Print(colors.end)
    if cursorLine > maxLine{ maxLine=cursorLine }
}

func main() {
    colors=colorstruct{green: "\033[32m", end: "\033[0m", red: "\033[31m"}
    cursorLine=0
    maxLine=0

    rand.Seed(time.Now().UnixNano())
    declare_alph()
    
    var question int
    var streak int
    streak=0
    

    tty, err := tty.Open()
    if err != nil {
        log.Fatal(err)
    }
    defer tty.Close()

    for {
        question=rand.Intn(alphlen)
        moveCursorUp(cursorLine)
        clearLines(strings.Count(alph[question],"\n")+1)
        printAndCountLines(alph[question],colors.end,true) // print the letter in default terminal color
        fmt.Print(streakString(streak,true))
        
        for{
            char, err := tty.ReadRune()
            if err != nil {
                log.Fatal(err)
            }
            // Correct answer
            if int(char)-97 == question {
                moveCursorUp(cursorLine)
                clearLines(strings.Count(alph[question],"\n"))
                printAndCountLines(alph[question],colors.green) 
                time.Sleep(100*time.Millisecond)
                streak+=1
                break
            }else{
                moveCursorUp(cursorLine)
                clearLines(strings.Count(alph[question],"\n")+1)
                printAndCountLines(alph[question],colors.red) 
                fmt.Print(streakString(streak,false))
                time.Sleep(50*time.Millisecond)
                moveCursorUp(cursorLine)
                clearLines(strings.Count(alph[question],"\n")+1)
                printAndCountLines(alph[question],colors.end) 

                streak=0
            }
        }
    }
}

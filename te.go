package main

import (
	"os"
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/mattn/go-runewidth"
	"bufio"
	"strconv"
	"strings"
)

var mode int
var ROWS, COLS int
var offX, offY int
var currX, currY int
var srcFile string

var textBuff = [][]rune{}

var undoBuff = [][]rune{}
var copyBuff = []rune{}

var modified bool

func readFile(filename string){
	file, err := os.Open(filename)

	if err != nil {
		srcFile = filename
		textBuff = append(textBuff, []rune{})
 	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		textBuff = append(textBuff, []rune(line))

		// for i :=0 ; i < len(line); i++ {
		// 	textBuff[lineNum] = append(textBuff[lineNum], rune(line[i]))
		// }
		lineNum++
	}

	if lineNum == 0 {
		textBuff = append(textBuff, []rune{})
	}
}



func writeFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for row, line := range textBuff {
		newLine := "\n"
		if row == len(textBuff)-1{
			newLine = ""
		}
		writeLine := string(line) + newLine
		_, err := writer.WriteString(writeLine)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		writer.Flush()
		modified = false
	}
}


func insertRune(event termbox.Event){
	insRune := make([]rune, len(textBuff[currY])+1)
	copy(insRune[:currX], textBuff[currY][:currX])

	if event.Key == termbox.KeySpace {
		insRune[currX] = rune(' ')
	}else if event.Key == termbox.KeyTab{
		insRune[currX] = rune(' ')
	}else{
		insRune[currX] = rune(event.Ch)
	}

	copy(insRune[currX+1:], textBuff[currY][currX:])

	textBuff[currY] = insRune
	currX++
}

func deleteRune() {
	if currX > 0 {
		currX--
		deleteLine := make([]rune, len(textBuff[currY])-1)
		copy(deleteLine[:currX], textBuff[currY][:currX])
		copy(deleteLine[currX:], textBuff[currY][currX+1:])
		textBuff[currY] = deleteLine
	}else if currY>0{
		appendLine := make([]rune, len(textBuff[currY]))
		copy(appendLine, textBuff[currY][currX:])
		newTextBuff := make([][]rune, len(textBuff)-1)
		copy(newTextBuff[:currY], textBuff[:currY])
		copy(newTextBuff[currY:], textBuff[currY+1:])
		textBuff = newTextBuff
		currY--
		currX = len(textBuff[currY])
		insertLine := make([]rune, len(textBuff[currY])+len(appendLine))
		copy(insertLine[:len(textBuff[currY])], textBuff[currY])
		copy(insertLine[len(textBuff[currY]):], appendLine)
		textBuff[currY] = insertLine
	}
}


func insertLine() {
	rightLine := make([]rune, len(textBuff[currY][currX:])) 
	copy(rightLine, textBuff[currY][currX:])
	leftLine := make([]rune, len(textBuff[currY][:currX]))
	copy(leftLine, textBuff[currY][:currX])
	textBuff[currY] = leftLine
	currY++
	currX = 0
	newTextBuff := make([][]rune, len(textBuff)+1)
	copy(newTextBuff, textBuff[:currY])
	newTextBuff[currY] = rightLine
	copy(newTextBuff[currY+1:], textBuff[currY:])
	textBuff = newTextBuff
}



func scrollTextBuff() {
	if currY < offY{
		offY = currY
	}
	if currX < offX {
		offX = currX
	}
	if currY >= offY + ROWS {
		offY = currY - ROWS + 1
	}
	if currX >= offX + COLS {
		offX = currX - COLS + 1
	}
}



func displayTextBuff() {
	for row := 0; row < ROWS; row++ {
		textBuffRow := row + offY
		for col := 0; col < COLS; col++ {
			textBuffCol := col + offX
			if textBuffRow >= 0 && textBuffRow < len(textBuff) && textBuffCol < len(textBuff[textBuffRow]) {
				ch := textBuff[textBuffRow][textBuffCol]
				if ch != '\t' {
					termbox.SetChar(col, row, ch)
				} else {
					termbox.SetCell(col+1, row, rune(' '), termbox.ColorDefault, termbox.ColorDefault)
				}
			}else if row+ offY >= len(textBuff) {
				termbox.SetCell(0, row, rune('*'), termbox.ColorBlue, termbox.ColorDefault)
				termbox.SetChar(col,row,rune('\n'))
			}
		}
	}
}


func displayStatusBar() {
	var modeStatus string
	var copyStatus string
	var undoStatus string
	var fileStatus string
	var cursorStatus string

	if mode >0{
		modeStatus = "EDIT MODE "
	}else{
		modeStatus = "CMD MODE "
	}

	filenameLength := len(srcFile)
	if filenameLength > 8{
		filenameLength = 8
	}

	fileStatus = srcFile[:filenameLength] + "  " + strconv.Itoa(len(textBuff)) + " " + "line(s)"

	if modified{
		fileStatus += " modified"
	}else{
		fileStatus += " saved"
	}

	cursorStatus = " Row " + strconv.Itoa(currY+1) + ", Col " + strconv.Itoa(currX+1) + " "
	
	if len(copyBuff) > 0 {
		copyStatus = " [Copy]"
	}

	if len(undoBuff) > 0 {
		undoStatus = " [Undo]"
	}

	usedSpace := len(modeStatus) + len(copyStatus) + len(undoStatus) + len(fileStatus) + len(cursorStatus)
	spaces := strings.Repeat(" ", COLS - usedSpace)

	message := modeStatus + fileStatus + copyStatus + undoStatus + spaces + cursorStatus

	printMsg(0, ROWS, termbox.ColorBlack, termbox.ColorWhite, message)

}



func printMsg(col, row int, fg, bg termbox.Attribute, msg string) {
	for _, ch := range msg {
		termbox.SetCell(col, row, ch, fg, bg)
		col += runewidth.RuneWidth(ch)
	}
}

func getKey() termbox.Event {
	var keyEvent termbox.Event

	switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey: keyEvent = event
		case termbox.EventError: panic(event.Err)
	} 
	return keyEvent
		
}

func keyPress(){
	keyEvent := getKey()
	if keyEvent.Key == termbox.KeyEsc {
		mode = 0
	}else if keyEvent.Key == termbox.KeyCtrlC {
		termbox.Close()
		os.Exit(0)
	}else if keyEvent.Ch != 0 {
		if mode == 1 {
			insertRune(keyEvent)
			modified = true
		}else{
			switch keyEvent.Ch {
				case 'i': mode = 1 
				case 'w': writeFile(srcFile)
			}

		}

	}else{
		switch keyEvent.Key {
			case termbox.KeyBackspace2: {
				if mode == 1 {
					deleteRune()
					modified = true
				}
			}
			case termbox.KeyBackspace: {
				if mode == 1 {
					deleteRune()
					modified = true
				}
			}
			case termbox.KeyTab: 
				if mode ==1{
					for i := 0; i < 4; i++ {
						insertRune(keyEvent)
						modified = true
					}
				}
			case termbox.KeySpace:
				if mode == 1 {
					insertRune(keyEvent)
					modified = true
				}
			case termbox.KeyEnter:
				if mode == 1 {
					insertLine()
					modified = true
				}
			case termbox.KeyHome: currY = 0
			case termbox.KeyEnd: currY = ROWS-1
			case termbox.KeyArrowUp: if currY != 0{currY--}
			case termbox.KeyArrowDown: if currY < len(textBuff)-1{currY++}
			case termbox.KeyArrowLeft: 
				if currX != 0{
					currX--
				}else if currY > 0{
					currY--
					currX = len(textBuff[currY])
				}
			case termbox.KeyArrowRight:
				if currX < len(textBuff[currY]){
					currX++
				}else if currY < len(textBuff)-1 {
					currX = 0
					currY++
				}
		}
		if currX > len(textBuff[currY]) {
			currX = len(textBuff[currY])
		}
	}
}

func main(){
	err :=  termbox.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing termbox: %v\n", err)
  		os.Exit(1)
	}

	if len(os.Args) > 1 {
		srcFile = os.Args[1]
		readFile(srcFile)
	}else{
		srcFile = "text.txt"
		textBuff = append(textBuff, []rune{})
	}


	for{
		COLS, ROWS = termbox.Size()
		ROWS--

		if COLS < 80{
			COLS = 80
		}
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		scrollTextBuff()

		displayTextBuff()

		displayStatusBar()

		termbox.SetCursor(currX-offX, currY-offY)
		termbox.Flush()
		
		keyPress()
	}
}
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: " + os.Args[0] + " <src-file> <dst-file>")
		fmt.Println("rest of the process is interactive.")
		os.Exit(-1)
	}
	srcfilename := os.Args[1]
	dstfilename := os.Args[2]
	srcfile, srcerr := os.Open(srcfilename)
	if srcerr != nil {
		fmt.Println("Cant open source file:" + srcerr.Error())
		os.Exit(-1)
	}
	defer srcfile.Close()
	dstfile, dsterr := os.Open(dstfilename)
	if dsterr != nil {
		fmt.Println("Cant open destination file:" + dsterr.Error())
		os.Exit(-1)
	}
	log.Println("Reading logisim destination file")
	destlines := readAllLines(bufio.NewReader(dstfile))
	log.Println("Starting partial parser for source file")
	destTokens := getRange(destlines, "circuit")
	dstfile.Close()
	// try parse
	log.Println("Reading logisim source file")
	allinput := readAllLines(bufio.NewReader(srcfile))
	log.Println("Starting partial parser for source file")
	expTokens := getRange(allinput, "circuit")
	log.Println("Parsing Complete. Initiating interactive process..")
	fmt.Println("Source circuits are:")
	for id, circ := range expTokens {
		fmt.Printf("[%2d] %s\n", id, circ.Params["name"])
	}
	fmt.Println("Destination circuits are:")
	for id, circ := range destTokens {
		fmt.Printf("[%2d] %s\n", id, circ.Params["name"])
	}
	for true {
		fmt.Println("Select operation:")
		fmt.Println("c - Copy a series of circuits to destination file.")
		fmt.Println("r - Copy a circuit to destination file, allow you to rename it.")
		fmt.Println("p - Reprint the source circuits list.")
		fmt.Println("d - Print destination circuits list.")
		fmt.Println("q - Quit the program.")
		var op byte = 'q'
		_, operr := fmt.Scanf("%c\n", &op)
		if operr != nil {
			fmt.Println("Operation reading failed.")
			fmt.Println("Prevent data loss, writing changes into file..")
			dstfile, dsterr = os.Create(dstfilename)
			for _, str := range destlines {
				_, destwrerr := dstfile.WriteString(str + "\n")
				if destwrerr != nil {
					log.Fatalf("Writing to destination file failed. Your file may be corrupted!!!")
				}
			}
			dstfile.Close()
			return
		}
		switch op {
		case 'c':
			fmt.Println("Enter a series of index, ending with a negative number.")
			tarIdx := make([]int, 0)
			for true {
				tmpIdx := -1
				_, idxerr := fmt.Scanf("%d", &tmpIdx)
				if idxerr != nil {
					fmt.Println("Index read failed, assume termination.")
					break
				}
				if tmpIdx < 0 {
					fmt.Println("Termination detected.")
					break
				}
				if tmpIdx >= len(expTokens) {
					fmt.Printf("Overflow detected. Ignoring %d\n", tmpIdx)
					continue
				}
				tarIdx = append(tarIdx, tmpIdx)
			}
			fmt.Println("Copy list are:")
			for _, idx := range tarIdx {
				fmt.Printf("[%2d] %s\n", idx, expTokens[idx].Params["name"])
			}
			fmt.Println("Completing copy process..")
			for _, idx := range tarIdx {
				tmpTok := expTokens[idx]
				tmpfooter := destlines[len(destlines)-1]
				destlines = append(destlines[0:len(destlines)-1], allinput[tmpTok.Top:tmpTok.Bottom+1]...)
				destlines = append(destlines, tmpfooter)
				log.Printf("Copy for id %d completed.\n", idx)
			}
			fmt.Println("Copy completed.")
		case 'r':
			fmt.Println("Enter a index number:")
			tmpIdx := -1
			_, idxerr := fmt.Scanf("%d", &tmpIdx)
			if idxerr != nil {
				fmt.Println("Index read failed.")
				continue
			}
			if tmpIdx < 0 {
				fmt.Println("Invalid Index.")
				continue
			}
			if tmpIdx >= len(expTokens) {
				fmt.Printf("Overflow detected.\n")
				continue
			}
			fmt.Println("Your selection:")
			fmt.Printf("[%2d] %s\n", tmpIdx, expTokens[tmpIdx].Params["name"])
			fmt.Println("Please input the new name, end with new line.")
			fmt.Println("To avoid unnecessary problems, USE ASCII Characters!!!")
			tmpIn := bufio.NewReader(os.Stdin)
			newName := readLine(tmpIn)
			newTok := expTokens[tmpIdx]
			newTok.Params["name"] = newName
			tmpfooter := destlines[len(destlines)-1]
			log.Println("New Circuit is: " + newTok.ToXML())
			destlines = append(destlines[0:len(destlines)-1], newTok.ToXML())
			destlines = append(destlines, allinput[newTok.Top+1:newTok.Bottom+1]...)
			destlines = append(destlines, tmpfooter)
			log.Println("Rename and copy OK.")
		case 'p':
			for id, circ := range expTokens {
				fmt.Printf("[%2d] %s\n", id, circ.Params["name"])
			}
		case 'd':
			log.Println("Reprogressing partial parser with new content..")
			destTokens = getRange(destlines, "circuit")
			log.Println("Completed.")
			for id, circ := range destTokens {
				fmt.Printf("[%2d] %s\n", id, circ.Params["name"])
			}
		case 'q':
			fmt.Println("Writing changes into file..")
			dstfile, dsterr = os.Create(dstfilename)
			for _, str := range destlines {
				_, destwrerr := dstfile.WriteString(str + "\n")
				if destwrerr != nil {
					log.Fatalf("Writing to destination file failed. Your file may be corrupted!!!")
				}
			}
			dstfile.Close()
			return
		default:
			continue
		}
	}
}

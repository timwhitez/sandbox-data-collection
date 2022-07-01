package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"encoding/base64"
)


func main(){
	xor := Xor{}
	fi, err := os.Open("Simple.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if len(string(a))<20{
			fmt.Println(string(a))

		}else{
			dec0,_ := base64.StdEncoding.DecodeString(string(a))
			fmt.Println(xor.dec(string(dec0)))
			fmt.Println("")
		}
	}


}
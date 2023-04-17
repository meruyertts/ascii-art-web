package printascii

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func AsciiWeb(inputStr, font string) (string, error) {
	myStr := inputStr
	fileName := ""
	result := ""
	var err error
	fileName = fileNameCheck(font)
	if fileName == "" {
		err = errors.New("only thinkertoy, standard and shadow files are available")
		return "", err
	}
	if !TxtFileCheck(fileName) {
		return "", errors.New("The txt file has been changed")
	}
	if !isASCII(myStr) {
		return "", errors.New("non-ASCII character was entered")
	}
	if len(myStr) == 0 {
		return "", errors.New("string was not entered")
	}

	result, err = splitWord(myStr, fileName)
	if err != nil {
		return "", err
	}
	return result, err
}

func splitWord(myStr, myFile string) (string, error) {
	re := regexp.MustCompile(`\r`)
	newStr := re.Split(myStr, -1)
	var err error
	var myArr [8]string
	var finalStr string
	for i := 0; i < len(newStr); i++ {
		if len(newStr[i]) > 0 {
			myArr, err = printWord(newStr[i], myFile)
			if err != nil {
				return "", err
			}
			for _, i := range myArr {
				finalStr += i + "\n"
			}

		}
		if newStr[i] == "" {
			finalStr += "\n"
		}
	}
	return finalStr, err
}

func fileNameCheck(fName string) string {
	myFile := ""
	switch fName {
	case "standard":
		myFile = "standard.txt"
	case "shadow":
		myFile = "shadow.txt"
	case "thinkertoy":
		myFile = "thinkertoy.txt"
	}

	return myFile
}

func TxtFileCheck(fileName string) bool {
	hashStandard := []byte{225, 148, 241, 3, 52, 66, 97, 122, 184, 167, 142, 28, 166, 58, 32, 97, 245, 204, 7, 163, 240, 90, 194, 38, 237, 50, 235, 157, 253, 34, 166, 191}
	hashShadow := []byte{184, 17, 37, 168, 183, 46, 207, 226, 35, 69, 169, 190, 218, 184, 99, 86, 141, 179, 152, 16, 96, 21, 242, 206, 76, 172, 130, 232, 162, 21, 7, 76}
	hashThinkertoy := []byte{236, 241, 252, 123, 255, 114, 166, 211, 68, 247, 17, 86, 18, 3, 196, 224, 126, 132, 206, 58, 147, 120, 23, 16, 71, 60, 235, 235, 128, 88, 253, 28}
	file, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
		return false
	}
	defer file.Close()
	buf := make([]byte, 30*1024)
	sha256 := sha256.New()
	for {
		n, err := file.Read(buf)
		if n > 0 {
			_, err := sha256.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	sum := sha256.Sum(nil)
	switch fileName {
	case "shadow.txt":
		if string(sum) == string(hashShadow) {
			return true
		}
	case "standard.txt":
		if string(sum) == string(hashStandard) {
			return true
		}
	case "thinkertoy.txt":
		if string(sum) == string(hashThinkertoy) {
			return true
		}
	}
	return false
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func printWord(s, fileName string) ([8]string, error) {
	var err error
	myLine := ""
	myarray := [8]string{}
	for _, char := range s {
		for line := 2; line <= 9; line++ {
			myrune := int(char)
			for i := ' '; i <= '~'; i++ {
				j := (int(i) - ' ') * 9
				if myrune == int(i) {
					firstline, err := readExactLine(fileName, line+j)
					if err != nil {
						log.Print(err)
						return [8]string{}, err
					}
					myLine += firstline
				}
			}
		}
		temp := strings.Split(myLine, "\n")
		for index, s := range temp[:len(temp)-1] {
			myarray[index] += s
		}
		myLine = ""
	}
	return myarray, err
}

func readExactLine(fileName string, lineNumber int) (string, error) {
	inputFile, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	br := bufio.NewReader(inputFile)
	for i := 1; i < lineNumber; i++ {
		_, _ = br.ReadString('\n')
	}
	str, err := br.ReadString('\n')
	if err != nil {
		return "", err
	}

	return str, nil
}

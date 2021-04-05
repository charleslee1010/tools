package toolkit

import (
	"bufio"
//	"io"
	"os"
	"errors"
)

 
func ReadTextFile(fileName string, handler func(int, string) bool) error {
	f, err := os.Open(fileName)
	defer f.Close()
	
	
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	lineNo := 0
	for scanner.Scan() {
		line := scanner.Text()
		
		if !handler(lineNo, line) {
			return errors.New("user abort")
		}
		lineNo ++
	}
	return nil
}
 

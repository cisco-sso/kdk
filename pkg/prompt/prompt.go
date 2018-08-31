package prompt

/* Usage
sp := Prompt{
	Text: "Mount your /keybase directory within KDK? [y/n] ",
	Loop: true,
	Validate: ValidateYorN,
}

result, err := sp.Run()
fmt.Printf("result %v, err %v\n", result, err)
*/

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type Prompt struct {
	Text     string
	Loop     bool
	Validate func(string) error
}

func (sp *Prompt) Run() (string, error) {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Print the description
		fmt.Print(sp.Text)

		// Block and read the input
		scanner.Scan()
		text := scanner.Text()

		// If no validation function exists, return the text immediately
		if sp.Validate == nil {
			return text, nil
		}

		// Otherwise, run the validation function
		if err := sp.Validate(text); err == nil {
			return text, nil
		} else {
			// If the function didn't validate, print why, and continue the loop
			fmt.Println(err)
		}

		// If we are not looping, break after the first iteration
		if sp.Loop == false {
			break
		}
	}
	return "", errors.New("Failed to capture valid input")
}

func ValidateYorN(input string) error {
	if input == "y" || input == "n" {
		return nil
	}
	return errors.New("Input must be 'y' or 'n'")
}

func ValidateDirExists(input string) error {

	if _, err := os.Stat(input); err == nil {
		return nil
	}
	return errors.New("Input directory must exist")
}

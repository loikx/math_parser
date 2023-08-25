package validation

import (
	"regexp"
	"strings"
)

func validationBraces(channel chan error, expression string) {
	cnt := 0
	for _, v := range expression {
		if v == '(' {
			cnt++
		}
		if v == ')' {
			cnt--
		}
	}

	if cnt == 0 {
		channel <- nil
		return
	}

	channel <- ErrInvalidBracesSequence
}

func validationRegular(channel chan error, expression string) {
	regex := regexp.MustCompile("[0-9()\\-+/*]*")
	if regex.MatchString(expression) {
		channel <- nil
		return
	}

	channel <- ErrInvalidRegex
}

func validationOperator(channel chan error, expression string) {
	operators := map[uint8]bool{'+': true, '-': true, '/': true, '*': true}
	for i := 1; i < len(expression); i++ {
		_, prev := operators[expression[i-1]]
		_, cur := operators[expression[i]]
		if prev && cur {
			channel <- ErrInvalidOperators
			return
		}
	}

	channel <- nil
}

func Validate(expression string) error {
	bracesChan := make(chan error)
	regularChan := make(chan error)
	operatorsChan := make(chan error)
	defer close(bracesChan)
	defer close(regularChan)
	defer close(operatorsChan)

	go validationBraces(bracesChan, expression)
	go validationRegular(regularChan, expression)
	go validationOperator(operatorsChan, expression)

	select {
	case brace := <-bracesChan:
		<-regularChan
		<-operatorsChan
		return brace
	case regular := <-regularChan:
		<-bracesChan
		<-operatorsChan
		return regular
	case operators := <-operatorsChan:
		<-bracesChan
		<-regularChan
		return operators
	}

}

func Purify(expression string) string {
	expression = strings.Replace(expression, " ", "", len(expression))

	if expression[0] == '-' {
		expression = "0" + expression
	}

	for i := 1; i < len(expression); i++ {
		if expression[i] == '-' && expression[i-1] == '(' {
			expression = expression[:i] + "0" + expression[i:]
		}
	}

	return expression
}

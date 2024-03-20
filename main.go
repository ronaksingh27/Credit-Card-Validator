package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Function to validate a credit card number using the Luhn algorithm
func isValidLuhn(cardNumber string) bool {
    var sum int
    isSecond := false

    // Iterate through each digit of the card number from right to left
    for i := len(cardNumber) - 1; i >= 0; i-- {
        digit, err := strconv.Atoi(string(cardNumber[i]))
        if err != nil {
            return false // Not a valid digit
        }

        // Double every second digit
        if isSecond {
            digit *= 2
            if digit > 9 {
                digit -= 9
            }
        }

        // Add the digit to the sum
        sum += digit

        // Toggle the flag for alternating digits
        isSecond = !isSecond
    }

    // Check if the sum is divisible by 10
    return sum%10 == 0
}

// Function to identify if a credit card number belongs to Visa
func isVisa(cardNumber string) bool {
    // Visa card numbers start with 4 and have a length of 13 or 16 digits
    return strings.HasPrefix(cardNumber, "4") && (len(cardNumber) == 13 || len(cardNumber) == 16)
}

// Function to identify if a credit card number belongs to MasterCard
func isMasterCard(cardNumber string) bool {
    // MasterCard card numbers start with 5 and the second digit is in the range of 1 to 5 (inclusive),
    // and have a length of 16 digits
    secondDigit, _ := strconv.Atoi(string(cardNumber[1]))
    return strings.HasPrefix(cardNumber, "5") && secondDigit >= 1 && secondDigit <= 5 && len(cardNumber) == 16
}

// Function to identify if a credit card number belongs to American Express
func isAmericanExpress(cardNumber string) bool {
    // American Express card numbers start with 34 or 37 and have a length of 15 digits
    return (strings.HasPrefix(cardNumber, "34") || strings.HasPrefix(cardNumber, "37")) && len(cardNumber) == 15
}

func creditCardHandler(w http.ResponseWriter, r *http.Request){


	requestBody , err := io.ReadAll(r.Body)
	if err != nil {
        http.Error(w, "error reading request body", http.StatusBadRequest)
        return
    }

	type Request struct{
		CardNumber string `json:"card_number"`
	}

	var req Request
	err = json.Unmarshal(requestBody,&req)
	if( err != nil ){
		http.Error(w,"Invalid JSON payload ",http.StatusBadRequest)
		return
	}

	cardNumber := req.CardNumber

	var issuer string
	if isVisa(cardNumber){
		issuer = "VISA"
	}else if isMasterCard(cardNumber){
		issuer = "MASTERCARD"
	}else if( isAmericanExpress(cardNumber)){
		issuer = "AMERICAN EXPRESS"
	}else{
		issuer = "SOME OTHER PAYMENT GATEWAY"
	}




	response := map[string]interface{}{
		"card_number" : cardNumber,
		"is_Valid " : isValidLuhn(cardNumber),
		"issuer" : issuer,
	}

	finalJson , err := json.Marshal(response)
	if( err != nil ){
		 panic(err)
	}

	w.Header().Set("Content-Type","application/json")

	w.Write(finalJson)


}

func main(){
	fmt.Println("Hello");
	r := mux.NewRouter();
	r.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Hello from </h1>"))
	}).Methods(("GET"));

	r.HandleFunc("/validate",creditCardHandler).Methods(("GET"));

	log.Fatal( http.ListenAndServe(":8000",r));
}
package handlers

// import "strings"

// func cleanBody(body string) string {
// 	invalidWords := map[string]string{
// 		"kerfuffle": "****",
// 		"sharbert":  "****",
// 		"fornax":    "****",
// 	}

// 	fields := strings.Fields(body)
// 	res := make([]string, 0, len(fields))

// 	for _, word := range fields {
// 		if censor, invalid := invalidWords[strings.ToLower(word)]; invalid {
// 			res = append(res, censor)
// 		} else {
// 			res = append(res, word)
// 		}
// 	}

// 	return strings.Join(res, " ")
// }

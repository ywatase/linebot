package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

/* --- JSON define */
type ApiLocation struct{}
type ApiContentMetadata struct{}

type ApiContent struct {
	Id              string             `json:id`
	ContentType     int                `json:contentType`
	From            string             `json:from`
	CreatedTime     int                `json:createdTime`
	To              []string           `json:to`
	ToType          int                `json:toType`
	ContentMetadata ApiContentMetadata `json:contentMetadata`
	Text            string             `json:text`
	Location        ApiLocation        `json:location`
}
type ApiResult struct {
	From        string     `json:from`
	FromChannel string     `json:fromCannel`
	To          []string   `json:to`
	ToChannel   int        `json:toChannel`
	EventType   string     `json:eventType`
	Id          string     `json:id`
	Content     ApiContent `json:content`
}
type MessageRecieve struct {
	Result []ApiResult `json:result`
}
type ApiResponse struct {
	Result []ApiResult `json:result`
}

/* ----------------------- */
/* --- controller          */
/* ----------------------- */
func apiRequest(w http.ResponseWriter, r *http.Request) {
	textContent := ApiContent{ContentType: 1, ToType: 1, Text: "hello"}
	result := ApiResult{ToChannel: 1441301333, EventType: "138311609100106403", Content: textContent}
	ret := ApiResponse{Result: []ApiResult{result}}
	request := ""

	// JSON return
	defer func() {
		// result
		outjson, err := json.Marshal(ret)
		if err != nil {
			fmt.Println(err) //TODO: change to log
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(outjson))
	}()

	// type check
	//     if r.Method != "POST" {
	//         ret.Status = 1
	//         ret.Code = "Not POST method"
	//         return
	//     }

	// request body
	rb := bufio.NewReader(r.Body)
	for {
		s, err := rb.ReadString('\n')
		request = request + s
		if err == io.EOF {
			break
		}
	}

	// JSON parse
	var msg MessageRecieve
	b := []byte(request)
	err := json.Unmarshal(b, &msg)
	if err != nil {
		result.Content = ApiContent{ContentType: 1, ToType: 1, Text: "JSON parse error."}
		return
	}
	result.To = []string{msg.Result[0].Content.From}

	// mecab parse
	//     result, err := mecab.Parse(dec.Sentence)
	//     if err == nil {
	//         for _, n := range result {
	//             ret.Result = append(ret.Result, n)
	//         }
	//     }

}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, world")
}

/* ----------------------- */
/* --- main                */
/* ----------------------- */
func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/callback", apiRequest)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

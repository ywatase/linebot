package main

import (
	"bufio"
	"bytes"
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
		err := postMessage(ret)
		if err != nil {
			fmt.Println(err) //TODO: change to log
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		outjson, err := json.Marshal(ret)
		if err != nil {
			fmt.Println(err) //TODO: change to log
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

func postMessage(apiRes ApiResponse) (err error){
	paramBytes, err := json.Marshal(apiRes)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://trialbot-api.line.me/v1/events", bytes.NewReader(paramBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Line-ChannelID", "1462234830")
	req.Header.Set("X-Line-ChannelSecret", "73d0eee6bb5721b60055b7d067c33383")
	req.Header.Set("X-Line-Trusted-User-With-ACL", "u0c08bbc74380a164ae5111714d7a1161")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
//     body, err := ioutil.ReadAll(resp.Body)
//     if err != nil {
//         return err
//     }
	return nil
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

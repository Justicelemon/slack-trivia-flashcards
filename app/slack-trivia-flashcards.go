package slacktriviaflashcards

import (
    "encoding/base64"
    "encoding/json"
    "golang.org/x/net/context"
    "google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
    "net/http"
)

type Question struct {
    Id           int      `json:"id"`
    Answer       string   `json:"answer"`
    Question     string   `json:"question"`
    Value        int      `json:"value"`
    AirDate      string   `json:"airdate"`
    CreatedAt    string   `json:"created_at"`
    UpdatedAt    string   `json:"updated_at"`
    CategoryId   int      `json:"category_id"`
    GameId       int      `json:"game_id"`
    InvalidCount int      `json:"invalid_count"`
    Category     Category `json:"category"`
}

type Category struct {
    Id         int    `json:"id"`
    Title      string `json:"title"`
    CreatedAt  string `json:"created_at"`
    UpdatedAt  string `json:"updated_at"`
    CluesCount int    `json:"clues_count"`
}

type SlackMessage struct {
    ResponseType string `json:"response_type"`
    Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
    Fallback string  `json:"fallback"`
    Fields   []Field `json:"fields"`
}
type Field struct {
    Title string `json:"title"`
    Value string `json:"value"`
    Short bool   `json:"short"`
}

func init() {
    http.HandleFunc("/api/1/trivia", triviaHandler)
    http.HandleFunc("/api/1/decode", decodeHandler)
}

func triviaHandler(w http.ResponseWriter, r *http.Request) {
    //Read the Request Parameter "command"
    c := appengine.NewContext(r)
    command := r.FormValue("command")
    text := r.FormValue("text")
    var err error
    var errResponse string
    var jsonMessage []byte

    //Ideally do other checks for tokens/username/etc
    if command == "/trivia" {
        if text == "random" {
            // Query jservice to get a random Jeopardy question
            var jServiceApiResponse []Question
            err = GetJSON(c, "http://jservice.io/api/random", &jServiceApiResponse)
            if err != nil {
                errResponse = "Unable to talk to JService: " + err.Error()
                goto SEND
            }

            if len(jServiceApiResponse) == 0 {
                errResponse = "JService did not return a question"
                goto SEND
            }

            slackMessage := SlackMessage{
                ResponseType: "in_channel",
                Attachments: []Attachment{
                    Attachment{
                        Fallback: jServiceApiResponse[0].Question,
                        Fields: []Field{
                            Field{
                                Title: "Category",
                                Value: jServiceApiResponse[0].Category.Title,
                                Short: false,
                            },
                            Field{
                                Title: "Question",
                                Value: jServiceApiResponse[0].Question,
                                Short: false,
                            },
                            Field{
                                Title: "Encoded Answer",
                                Value: encode(([]byte(jServiceApiResponse[0].Answer))),
                                Short: false,
                            },
                        },
                    },
                },
            }
            jsonMessage, err = json.MarshalIndent(slackMessage, "", "  ")
            if err != nil {
                errResponse = "Error marshalling the Slack Message: " + err.Error()
                goto SEND
            }
        } else {
            errResponse = "I do not understand your message " + text + " please type one of: random"
            goto SEND
        }
    } else {
        errResponse = "I do not understand your command: " + command + ", Please use /decode {encoded string} or /trivia random instead"
        goto SEND
    }

SEND:
    if errResponse != "" {
        jsonMessage, err := json.Marshal(map[string]string{"response_type": "in_channel", "text": errResponse})
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        w.Write(jsonMessage)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonMessage)
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
    command := r.FormValue("command")
    text := r.FormValue("text")
    var response string
    if command == "/decode" {
        enbyte, err := decode(text)
        if err != nil {
            response = (err.Error())
        } else {
            response = (string(enbyte))
        }
    } else {
        response = "I do not understand your command: " + command + ", Please use /decode {encoded string} or /trivia random instead"
    }
    slackResponse := map[string]string{"response_type": "in_channel", "text": response}
    js, err := json.Marshal(slackResponse)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func encode(src []byte) string {
    return base64.StdEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
    return base64.StdEncoding.DecodeString(src)
}

func GetJSON(c context.Context, url string, target interface{}) error {
    client := &http.Client{
        Transport: &urlfetch.Transport{
            Context: c,
            AllowInvalidServerCertificate: true,
        },
    }
    response, err := client.Get(url)
    if err != nil {
        return err
    }
    defer response.Body.Close()
    json.NewDecoder(response.Body).Decode(target)
    return nil
}

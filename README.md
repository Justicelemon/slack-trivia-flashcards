Slack Trivia Flashcards
=======================

Slack Trivia Flashcards is an API designed to be used on Slack as a slash command. Slack Trivia Flashcards talks to [jService](jservice.io) to grab a random Jeopardy question and encodes the answer before sending it back to Slack. If the users want to see the answer to the quesiton, they can also use Slack Trivia Flashcards to decode the answer. 

Dependencies
-----
- Python 2.7
- Google App Engine SDK for Go

Endpoints
-----

### Endpoint - Trivia

The trivia endpoint is reachable at `/api/1/trivia`. The command must be `/trivia`.

#### Endpoint - Trivia (Random)

Setting the text to `random` will fetch a random question from jService. 

Example: `http://localhost:8080/api/1/trivia?command=/trivia&text=random`

Slack Example: `/trivia random`

Example Response: 

```json
{
  "response_type": "in_channel",
  "attachments": [
    {
      "fallback": "This U.S. general who loved horses \u0026 studied at a cavalry school helped protect the Lipizzaners in WWII",
      "fields": [
        {
          "title": "Category",
          "value": "by george, it's george",
          "short": false
        },
        {
          "title": "Question",
          "value": "This U.S. general who loved horses \u0026 studied at a cavalry school helped protect the Lipizzaners in WWII",
          "short": false
        },
        {
          "title": "Encoded Answer",
          "value": "R2VvcmdlIFBhdHRvbg==",
          "short": false
        }
      ]
    }
  ]
}
```

### Endpoint - Decode

The decode endpoint is reachable at `/api/1/decode`. The command must be `/trivia` and the text should be one of the encoded answers that one of the calls to the trivia endpoint.

Example: `http://localhost:8080/api/1/decode?command=/decode&text=R2VvcmdlIFBhdHRvbg==`

Slack Example: `/decode dGhlIFRvbWIgb2YgdGhlIFVua25vd24gU29sZGllcg==`

Example Response:

```json
{"response_type":"in_channel","text":"George Patton"}
```

Because of the way HTTP urls are encoded, certain encoded answers might give `illegal base64 data at input byte 3` when accessed through a web browser directly. This issue hasn't shown up on Slack yet with the same encoded answers that caused those issues.

Running Locally
-----
Run `goapp serve app/` in the main directory or just `goapp serve` in the app directory to start the local development server. The app should now be accessible at `http://localhost:8080`, though just typing `http://localhost:8080` without any commands leads to the app responding with a 404. The development app server watches for changes in the files, so there's no need to restart `goapp serve` as development is being done on it.


Deploying
-----
1. Sign in to [App Engine](https://appengine.google.com/) with your Google Account and click "Create Application"
2. Fill in the form for creating your application, make sure to note what you entered for "Application Identifier"
3. Change the value of the `application:` setting in app.yaml from `applicationid` to the Application Identifier set in step 2
4. Upload the app by running `goapp deploy app/` and enter your Google email ID and password

The app should now be available at `http://{app_id}.appspot.com`

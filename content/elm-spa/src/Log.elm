port module Log exposing (debug, error, stringHttpError)

import Http


port debug : String -> Cmd msg


port error : String -> Cmd msg


stringHttpError : Http.Error -> String
stringHttpError httpError =
    case httpError of
        Http.BadUrl url ->
            "Bad URL: " ++ url

        Http.Timeout ->
            "Timeout"

        Http.NetworkError ->
            "Network error"

        Http.BadStatus _ ->
            "Bad status"

        Http.BadPayload errorMessage _ ->
            "Bad payload, failed to parse payload: " ++ errorMessage

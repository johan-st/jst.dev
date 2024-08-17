module main

import Browser

main = Browser.Document
    { init = init
    , update = update
    , view = view
    , subscriptions = subscriptions
    }


-- MODEL

type alias Model 
    = Bool

init : Model
init = False

-- UPDATE

type Msg
    = Toggle

update : Msg -> Model -> Model
update msg model =
    case msg of
        Toggle ->
            not model


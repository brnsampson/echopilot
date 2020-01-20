module Main exposing (..)

-- INIT

import Browser
import Html exposing (Html, Attribute, div, input, text, button)
import Html.Events exposing (onClick, onInput)
import Html.Attributes exposing (..)
import Http
import Url.Builder exposing (relative)
import Json.Decode exposing (Decoder, field, string)
import Json.Encode

main = Browser.element {init = init, view = view, update = update, subscriptions = subscriptions}

-- MODEL

type alias Model = { user_string : String, echo_string: String }

init : () -> (Model, Cmd Msg)
init _ =
  (
    { user_string = "", echo_string = "" },
    Cmd.none
  )

-- VIEW
view : Model -> Html Msg
view model =
  div [] [
    input [ placeholder "Text to echo", value model.user_string, onInput UpdateString ] []
    ,button [ onClick Echo ] [ text "Click to echo!" ]
    ,div [] [ text model.echo_string ]
  ]

-- UPDATE

type Msg =
  Echo
  | UpdateString String
  | GotJson (Result Http.Error String)

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  case msg of
    UpdateString string ->
      (
        { model | user_string = string }
        , Cmd.none
      )
    Echo ->
      (
        model
        , echoString model.user_string
      )
    GotJson result ->
      case result of
        Ok data ->
          (
            { model | echo_string = data }
            , Cmd.none
          )
        Err error ->
          case error of
            Http.BadUrl str ->
              (
                { model | echo_string = "Bad Url: " ++ str }
                , Cmd.none
              )
            Http.Timeout ->
              (
                { model | echo_string = "Timeout!" }
                , Cmd.none
              )
            Http.NetworkError ->
              (
                { model | echo_string = "Network Error!" }
                , Cmd.none
              )
            Http.BadStatus status ->
              (
                { model | echo_string = "Bad Status: " ++ String.fromInt status }
                , Cmd.none
              )
            Http.BadBody body ->
              (
                { model | echo_string = "Bad Body: " ++ body }
                , Cmd.none
              )
  

-- SUBSCRIPTIONS
subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none

-- HTTP
echoString : String -> Cmd Msg
echoString value =
  Http.post
    { url = relative ["v1", "api", "echo"] []
    , body = Http.jsonBody (echoJson value)
    , expect = Http.expectJson GotJson jsonDecoder
    }

echoJson : String -> Json.Encode.Value
echoJson value =
  Json.Encode.object [ ( "content", Json.Encode.string value ) ]

jsonDecoder : Decoder String
jsonDecoder =
  field "content" string

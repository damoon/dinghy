module Main exposing (main)

import Browser
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (class, id, style, src, href)
import Html.Events exposing (..)
import Http
import Json.Decode exposing (Decoder, field, string, int, bool, list, decodeString, map, map3, map4, maybe)
import Url
import List exposing (concat)
import Delay
import Process
import Task
import String
import Time

main : Program () Model Msg
main =
  Browser.application
    { init = init
    , view = view
    , update = update
    , subscriptions = subscriptions
    , onUrlChange = UrlChanged
    , onUrlRequest = LinkClicked
    }

type alias Model =
  { key : Nav.Key
  , url : Url.Url
  , dir : Maybe Directory
  , fetching : Fetching
  }
  
type Fetching
  = Loading
  | LoadingSlowly
  | Loaded
  | Failed String


init : () -> Url.Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url key =
  ( Model key url Nothing Loading
  , Cmd.batch [ getRandomCatGif url.path
              , delay 500 (LoadingIsSlow url.path)
              ]
  )


-- UPDATE


type Msg
  = LoadingIsSlow String
  | GotGif (Result Http.Error Directory)
  | LinkClicked Browser.UrlRequest
  | UrlChanged Url.Url
  | Tick Time.Posix

delay : Float -> msg -> Cmd msg
delay time msg =
  Process.sleep time
  |> Task.perform (\_ -> msg)

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of

    Tick newTime ->
      case model.fetching of
        Loading ->
          ( model, Cmd.none )
        LoadingSlowly ->
          ( model, Cmd.none )
        _ ->
          ( { model | fetching = Loading }
          , Cmd.batch [ getRandomCatGif model.url.path
                      , delay 500 (LoadingIsSlow model.url.path) ] )

    LoadingIsSlow path ->
      case model.fetching of
        Loading ->
          if path == model.url.path then
            ( { model | fetching = LoadingSlowly }
            , Cmd.none )
          else
            ( model, Cmd.none )
        _ ->
          ( model, Cmd.none )


    GotGif result ->
      case result of
        Ok dir ->
          ( { model | fetching = Loaded, dir = Just dir }
          , Cmd.none )

        Err txt ->
          ( { model | fetching = Failed (errorToString txt) }
          , Cmd.none )

    LinkClicked urlRequest ->
      case urlRequest of
        Browser.Internal url ->
          ( model
          , Nav.pushUrl model.key (Url.toString url) )

        Browser.External href ->
          ( model
          , Nav.load href )

    UrlChanged url ->
      if url.path == model.url.path then
        ( model, Cmd.none )
      else
        ( { model | url = url, fetching = Loading }
        , Cmd.batch [ getRandomCatGif url.path
                    , delay 500 (LoadingIsSlow url.path)
                    ]
        )

errorToString : Http.Error -> String
errorToString err =
    case err of
        Http.Timeout ->
            "Timeout exceeded"

        Http.NetworkError ->
            "Network error"

        Http.BadStatus num ->
            "Http Status: " ++ String.fromInt num

        Http.BadBody text ->
            "Unexpected response from api: " ++ text

        Http.BadUrl url ->
            "Malformed url: " ++ url

-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
  Time.every 1000 Tick


-- VIEW

view : Model -> Browser.Document Msg
view model =
  { title = "URL Interceptor"
  , body =
      [ div []
          [ h1 []
              [ img [ src "/favicon.png" ] []
              , text "Dinghy"
              ]
          , viewFetching model.fetching
          , viewDirectory model.dir
          ]
      ]
  }


viewFetching : Fetching -> Html Msg
viewFetching state =
  case state of
    Loading ->
      text ""
    LoadingSlowly ->
      errorBox "Loading..."
    Failed err ->
      errorBox err
    Loaded ->
      text ""


errorBox : String -> Html Msg
errorBox err =
  div
    [ id "error" ]
    [ text err ]


viewDirectory : Maybe Directory -> Html Msg
viewDirectory dira =
  case dira of
    Nothing ->
      text ""
    Just dir ->
      div [ ]
      (concat
        [ [ h2 [] (navigation dir.path)
          , br [] []
          ]
        , List.map viewFolder dir.directories
        , List.map viewFile dir.files
        ])

navigation : String -> List (Html Msg)
navigation path =
  let
    elements = String.split "/" path
    links = navigationElements "/" elements
  in
    concat
    [ [ a [ href "/" ] [ text "Root" ] ]
    , links
    ]

navigationElements : String -> List String -> List (Html Msg)
navigationElements previous elements =
  let
    h = List.head elements
    url = case h of
      Nothing ->
        ""
      Just name ->
        previous ++ name ++ "/"
    t = List.tail elements
    ls = case t of
      Nothing ->
        []
      Just xs ->
        navigationElements url xs
  in
  case (h, t) of
    (Nothing, _) ->
      []
    (Just txt, _) ->
      concat
        [ [ text " / "
          , a [ href url ] [ text txt ]
          ]
        , ls
        ]


viewFolder : String -> Html Msg
viewFolder name =
  div 
    [ class "icon" ]
    [ div 
      [ class "icon-inner" ]
      [ a 
        [ href (name++"/")
        , class "folder"
        ]
        [ span
          [ class "fiv-sqo fiv-icon-folder"
          , style "width" "72px"
          , style "height" "72px"
          , style "margin" "3px"
          ]
          []
        , br [] []
        , text name
        ]
      ]
    ]

viewFile : File -> Html Msg
viewFile fi =
  let
    iconclass = "fiv-sqo fiv-icon-" ++ fi.icon
  in
  div 
    [ class "icon" ]
    [ div 
      [ class "icon-inner" ]
      [ a 
        [ href fi.downloadURL ]
        [ span
          [ class iconclass
          , style "width" "72px"
          , style "height" "72px"
          , style "margin" "3px"
          ]
          []
        , br [] []
        , text fi.name
        ]
      ]
    ]

-- HTTP


getRandomCatGif : String -> Cmd Msg
getRandomCatGif path =
  Http.request
    { method = "GET"
    , url = "http://backend:8080" ++ path
    , body = Http.emptyBody
    , headers = [
      Http.header "Accept" "application/json;q=0.9"
    ]
    , expect = Http.expectJson GotGif directoryDecoder
    , timeout = Just 5000
    , tracker = Nothing
    }


directoryDecoder : Decoder Directory
directoryDecoder =
    map3 Directory
        (field "Path" string)
        (field "Directories" directoriesDecoder)
        (field "Files" filesDecoder)

type alias Directory =
  { path        : String
  , directories : List String
  , files       : List File
  }

type alias File =
  { name : String
  , size : Int
  , downloadURL : String
  , icon : String
--  , icon : Maybe String
  }

directoriesDecoder : Decoder (List String)
directoriesDecoder =
  Json.Decode.list string

filesDecoder : Decoder (List File)
filesDecoder =
  list fileDecoder

fileDecoder : Decoder File
fileDecoder =
  map4 File
    (field "Name" string)
    (field "Size" int)
    (field "DownloadURL" string)
    (field "Icon" string)
--    (maybe (field "Icon" string))

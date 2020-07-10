module Main exposing (main)

import Browser
import Browser.Navigation as Nav
import Cmd.Extra exposing (withCmd, withCmds, withNoCmd)
import Html exposing (Html, a, span, text, img, h1, div, br)
import Html.Attributes exposing (href, class, id, src, width, height, alt)
import Html.Events exposing (onClick)
import Json.Decode as JD exposing (decodeString, Decoder, field, string, bool, int, list, map, map3, map7, maybe)
import Json.Encode exposing (Value)
import PortFunnel.WebSocket as WebSocket exposing (Response(..))
import PortFunnels exposing (FunnelDict, Handler(..), State)
import Url
import List exposing (concat)
import Process
import Task
import String
--import Debug


handlers : List (Handler Model Msg)
handlers =
    [ WebSocketHandler socketHandler
    ]


subscriptions : Model -> Sub Msg
subscriptions model =
    PortFunnels.subscriptions Process model


funnelDict : FunnelDict Model Msg
funnelDict =
    PortFunnels.makeFunnelDict handlers getCmdPort


getCmdPort : String -> Model -> (Value -> Cmd Msg)
getCmdPort moduleName _ =
    PortFunnels.getCmdPort Process moduleName False


-- MODEL

type alias Config =
    { backend : String
    , websocket : String
    }

type alias Model =
    { endpoint : String
    , wasLoaded : Bool
    , state : State
    , key : String
    , error : Maybe String
    , nav : Nav.Key
    , url : Url.Url
    , dir : Maybe Directory
    , fetching : Fetching
    , backend : String
    }


type alias Directory =
  { path        : String
  , directories : List String
  , files       : List File
  }


type alias File =
  { name : String
  , path : String
  , size : Int
  , downloadURL : String
  , icon : String
  , thumbnail : Maybe String
  , archive : Bool
  }


type Fetching
  = Loading
  | LoadingSlowly
  | Loaded
  | Failed String


main : Program Config Model Msg
main =
  Browser.application
    { init = init
    , view = view
    , update = update
    , subscriptions = subscriptions
    , onUrlChange = UrlChanged
    , onUrlRequest = LinkClicked
    }


delay : Float -> msg -> Cmd msg
delay time msg =
    Process.sleep time
        |> Task.perform (\_ -> msg)


init : Config -> Url.Url -> Nav.Key -> ( Model, Cmd Msg )
init cfg url key =
    let
        model = { endpoint = cfg.websocket
                , backend = cfg.backend
                , wasLoaded = False
                , state = PortFunnels.initialState
                , key = "socket"
                , error = Nothing
                , url = url
                , nav = key
                , dir = Nothing
                , fetching = Loading
                }
    in
        model
        |> withCmd (delay 0 Startup)


-- UPDATE


type Msg
  = Startup
  | LoadingIsSlow String
  | Extract String
  | Delete String
  | LinkClicked Browser.UrlRequest
  | UrlChanged Url.Url
  | Process Value


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of

    Startup ->
        model
        |> withCmd (
            if model.wasLoaded then
                WebSocket.makeOpenWithKey model.key model.endpoint |> send model
            else
                delay 10 Startup
        )

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

    Extract path ->
      model
      |> withCmd (WebSocket.makeSend model.key ( "ex " ++ path ) |> send model)

    Delete path ->
      model
      |> withCmd (WebSocket.makeSend model.key ( "rm " ++ path ) |> send model)

    LinkClicked urlRequest ->
      case urlRequest of
        Browser.Internal url ->
          if String.startsWith model.backend (Url.toString url) then
            ( model
            , Nav.load (Url.toString url) )
          else
            ( model
            , Nav.pushUrl model.nav (Url.toString url) )

        Browser.External href ->
          ( model
          , Nav.load href )

    UrlChanged url ->
      if url.path == model.url.path then
        ( model, Cmd.none )
      else
        let
            mdl = { model | url = url, fetching = Loading }
        in
            mdl
            |> withCmds [ WebSocket.makeSend mdl.key ( "cd " ++ mdl.url.path ) |> send mdl
                        , delay 500 (LoadingIsSlow mdl.url.path)
                        ]

    Process value ->
        case
            PortFunnels.processValue funnelDict value model.state model
        of
            Err error ->
                { model | error = Just error } |> withNoCmd
            Ok res ->
                res


send : Model -> WebSocket.Message -> Cmd Msg
send model message =
    WebSocket.send (getCmdPort WebSocket.moduleName model) message


doIsLoaded : Model -> Model
doIsLoaded model =
    if not model.wasLoaded && WebSocket.isLoaded model.state.websocket then
        { model
            | wasLoaded = True
        }

    else
        model


socketHandler : Response -> State -> Model -> ( Model, Cmd Msg )
socketHandler response state mdl =
    let
        model =
            doIsLoaded
                { mdl
                    | state = state
                    , error = Nothing
                }
    in
    case response of
        WebSocket.MessageReceivedResponse { message } ->
            let
                result = decodeString directoryDecoder message
            in
            case result of
                Ok dir ->
                    { model | fetching = Loaded, dir = Just dir }
                        |> withNoCmd
                Err txt ->
                    { model | fetching = Failed (JD.errorToString txt) }
                        |> withNoCmd

        WebSocket.ConnectedResponse _ ->
            model
            |> withCmds [ WebSocket.makeSend model.key ( "cd " ++ model.url.path ) |> send model
                        , delay 500 (LoadingIsSlow model.url.path)
                        ]

        WebSocket.ClosedResponse _ ->
            model
                |> withNoCmd

        WebSocket.ErrorResponse error ->
            { model | error = Just (WebSocket.errorToString error) }
                |> withNoCmd

        _ ->
            case WebSocket.reconnectedResponses response of
                [] ->
                    model |> withNoCmd

                [ ReconnectedResponse _ ] ->
                    model
                        |> withCmds [ WebSocket.makeSend model.key ( "cd " ++ model.url.path ) |> send model
                                    , delay 500 (LoadingIsSlow model.url.path)
                                    ]

                _ ->
--                    { model | error = Just (Debug.toString list) }
                    { model | error = Just "unknown error" }
                        |> withNoCmd


-- VIEW


br : Html msg
br =
    Html.br [] []


view : Model -> Browser.Document Msg
view model =
  let
    title = case model.dir of
              Nothing ->
                "dinghy"
              Just dir ->
                "dinghy/" ++ dir.path
  in
  { title = title
  , body =
      [ div []
          [ viewFetching model.fetching
          , div [ id "repo" ] [
            a
              [ href "https://github.com/damoon/dinghy" ]
              [ img [ src "/repo.png", width 32, height 32, alt "open issue" ] [] ]
          ]
          , h1 []
              ( concat [
                [ img [ src "/favicon.png", width 32, height 32, alt "Logo" ] []
                , a [ href "/" ] [ text "Dinghy" ]
                ]
                , navigation model.dir
              ] )
          , viewDirectory model.backend model.dir
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


navigation : Maybe Directory -> List (Html Msg)
navigation maybeDir =
  case maybeDir of
    Nothing ->
      [ text "" ]
    Just dir ->
      let
        elements = String.split "/" dir.path
      in
        navigationElements "/" elements


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
    (Just _, Just []) ->
      []
    (Just txt, _) ->
      concat
        [ [ text " / "
          , a [ href url ] [ text txt ]
          ]
        , ls
        ]


viewDirectory : String -> Maybe Directory -> Html Msg
viewDirectory backend maybeDir =
  case maybeDir of
    Nothing ->
      text ""
    Just dir ->
      div [ ]
      (concat
        [ List.map (viewFolder dir.path) dir.directories
        , List.map (viewFile backend) dir.files
        ])


viewFolder : String -> String -> Html Msg
viewFolder path name =
  let
    delete = img [ src "/delete.png"
                 , class "button"
                 , onClick (Delete ("/"++path++name))
                 ] []
  in
  div 
    [ class "icon" ]
    [ div 
      [ class "icon-inner" ]
      [ a [ href (name++"/")
          , class "folder"]
          [ span[ class "fiv-sqo fiv-icon-folder fiv-icon" ] []
          , br
          , text name
          ]
      , br
      , delete
      ]
    ]


viewFile : String -> File -> Html Msg
viewFile backend fi =
  let
    extract = if fi.archive then
                img [ src "/extract.gif"
                    , class "button"
                    , onClick (Extract fi.path)
                    ] []
              else
                text ""
    delete = img [ src "/delete.png"
                 , class "button"
                 , onClick (Delete fi.path)
                 ] []
  in
  div 
    [ class "icon" ]
    [ div 
      [ class "icon-inner" ]
      [ a [ href (backend ++ "/" ++ fi.downloadURL) ] (viewIcon backend fi)
      , br
      , extract
      , delete
      ]
    ]


viewIcon : String -> File -> List (Html msg)
viewIcon backend fi =
  case fi.thumbnail of
    Nothing ->
      [ span [ class ("fiv-sqo fiv-icon-" ++ fi.icon ++ " fiv-icon") ] []
      , br
      , text fi.name
      ]
    Just url ->
      [ div
        [ class "thumbnail" ]
        [ img [ src (backend ++ "/" ++ url) ] []
        ]
      , text fi.name
      ]


-- Decode

directoryDecoder : Decoder Directory
directoryDecoder =
    map3 Directory
        (field "Path" string)
        (field "Directories" directoriesDecoder)
        (field "Files" filesDecoder)


directoriesDecoder : Decoder (List String)
directoriesDecoder =
  JD.list string

filesDecoder : Decoder (List File)
filesDecoder =
  list fileDecoder

fileDecoder : Decoder File
fileDecoder =
  map7 File
    (field "Name" string)
    (field "Path" string)
    (field "Size" int)
    (field "DownloadURL" string)
    (field "Icon" string)
    (maybe (field "Thumbnail" string))
    (field "Archive" bool)

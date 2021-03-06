<!DOCTYPE HTML>
<html>
<head>
  <meta charset="UTF-8">
  <title>dinghy</title>
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="icon" type="image/png" href="/favicon.png" />
  <link rel="stylesheet" type="text/css" href="//cdn.jsdelivr.net/gh/necolas/normalize.css@8.0.1/normalize.css" />
  <link rel="stylesheet" type="text/css" href="//cdn.jsdelivr.net/npm/file-icon-vectors@1.0.0/dist/file-icon-vivid.min.css" />
  <link rel="stylesheet" type="text/css" href="//cdn.jsdelivr.net/npm/file-icon-vectors@1.0.0/dist/file-icon-square-o.min.css" />
  <link rel="stylesheet" type="text/css" href="//fonts.googleapis.com/css?family=Roboto+Slab">
  <link rel="stylesheet" type="text/css" href="/theme.css">
</head>
<body>
  <div id="elm"></div>
  <script src="/main.js"></script>
  <script src="//cdn.jsdelivr.net/gh/billstclair/elm-websocket-client@4.1.0/example/site/js/PortFunnel.js"></script>
  <script src="//cdn.jsdelivr.net/gh/billstclair/elm-websocket-client@4.1.0/example/site/js/PortFunnel/WebSocket.js"></script>
  <script>
    var fmt = localStorage.getItem("viewFormat") || "GridView";
    var app = Elm.Main.init({
      node:  document.getElementById("elm"),
      flags: { backend: "${BACKEND_URL}"
             , websocket: "${WEBSOCKET_URL}"
             , format: fmt
             }
    });
    PortFunnel.subscribe(app);
    app.ports.saveViewFormat.subscribe(function(value) {
      localStorage.setItem("viewFormat", value);
    });
  </script>
</body>
</html>

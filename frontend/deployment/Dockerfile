# build-env ###################################################
FROM ymtszw/elm_for_ci:0.19.1 AS build-env

RUN yarn global add elm-live

WORKDIR /app

COPY . /app

RUN elm make src/Main.elm --optimize --output=./dist/main.js
RUN uglifyjs ./dist/main.js --compress "pure_funcs=[F2,F3,F4,F5,F6,F7,F8,F9,A2,A3,A4,A5,A6,A7,A8,A9],pure_getters,keep_fargs=false,unsafe_comps,unsafe" | uglifyjs --mangle --output=./dist/main.js

CMD elm-live src/Main.elm --host 0.0.0.0 --port 8000 --dir=./dist --pushstate -- --output=./dist/main.js

###############################################################
FROM nginx:1.19.0-alpine
COPY deployment/nginx.conf /etc/nginx/nginx.conf
COPY deployment/default.conf /etc/nginx/conf.d/default.conf
COPY --from=build-env /app/dist /usr/share/nginx/html

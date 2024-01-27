# build-env ###################################################
FROM ymtszw/elm_for_ci:0.19.1@sha256:0871be293f662293b0e2cb81b2378eca9237242045e68434b4f6f22fb7fbfc7a AS build-env

RUN yarn global add elm-live
RUN apk add gettext

WORKDIR /app

COPY . /app

CMD cat dist/index.html.tpl | envsubst > dist/index.html && \
    elm-live src/Main.elm --host 0.0.0.0 --port 8000 --dir=./dist --pushstate -- --output=./dist/main.js

# optimize-env ###################################################
FROM build-env AS optimize-env

RUN elm make src/Main.elm --optimize --output=./dist/main.js
RUN uglifyjs ./dist/main.js --compress "pure_funcs=[F2,F3,F4,F5,F6,F7,F8,F9,A2,A3,A4,A5,A6,A7,A8,A9],pure_getters,keep_fargs=false,unsafe_comps,unsafe" | uglifyjs --mangle --output=./dist/main.js

###############################################################
FROM nginx:1.25.3-alpine@sha256:363e10bd564f68a966d3fd3e5e59b701f8d0e36b68ef8517d0527bf7e844092c

COPY deployment/envsubst-index.sh /docker-entrypoint.d/15-envsubst-index.sh

COPY deployment/nginx.conf /etc/nginx/nginx.conf
COPY deployment/default.conf /etc/nginx/conf.d/default.conf
COPY --from=optimize-env /app/dist /usr/share/nginx/html

FROM node:14.3.0-alpine3.11

RUN apk add bash git
RUN yarn global add elm elm-test elm-analyse elm-live

WORKDIR /app

COPY . /app

CMD elm-live src/Main.elm -w=true --dir=./dist --pushstate -- --output=./dist/main.js

# https://github.com/rtfeldman/elm-spa-example
# $ elm make src/Main.elm --output elm.js --optimize
# $ uglifyjs elm.js --compress 'pure_funcs="F2,F3,F4,F5,F6,F7,F8,F9,A2,A3,A4,A5,A6,A7,A8,A9",pure_getters=true,keep_fargs=false,unsafe_comps=true,unsafe=true,passes=2' --output=elm.js && uglifyjs elm.js --mangle --output=elm.js

#elmfmt: ##@qa format code and remove unused imports
#	find src | grep \.elm\$$ | xargs -L 1 elm-impfix --replace
#	find src | grep \.elm\$$ | xargs elm-format --yes
#
#analyse: ##@qa run static code analyser
#	elm-analyse . -s

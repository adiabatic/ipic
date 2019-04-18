# Remember that ipic uses the builtin flag package.
# “-h” is properly treated like `gcc -Wall`, not `ls -laFh`.
complete -c ipic -xo a --description "Search for an album"
complete -c ipic -xo b --description "Search for a book"
complete -c ipic -xo f --description "Search for a film"
complete -c ipic  -o h --description "Show this help message"
complete -c ipic -xo i --description "Search for an iOS app"
complete -c ipic -xo m --description "Search for a macOS app"
complete -c ipic -xo t --description "Search for a TV show"

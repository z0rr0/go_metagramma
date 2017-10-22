# GoMetagramma

It is a metagramma game - find the shortest way from a word to other one changing only one symbol.

1. Prepare "database" `dict.json` using file `dict.txt` with words list.

    ```sh
    ./go_metagramma -i dict.txt -o dict.json
    ```

2. Find a way:

    ```sh
    ./go_metagramma -d dict.json -f стол -t хлеб
    
    162164 items are read from dict.json
    0: стол
    1: стой
    2: слой
    3: злой
    4: злей
    5: влей
    6: влев
    7: хлев
    8: хлеб
    duration 575.586798ms
    ```

## Build

It was checked on Go 1.9

```sh
go install github.com/z0rr0/go_metagramma
```

Run test:

```sh
go test -race -v -cover -coverprofile=coverage.out -trace trace.out github.com/z0rr0/go_metagramma
go tool cover -html=coverage.out
```

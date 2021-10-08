# atm
An ATM in the console

## Usage
Start the `./atm` binary and enter one of the following commands:
```
authorize <account_id> <pin>
withdraw <amount>
deposit <amount>
balance
history
logout
end
```

## Development
To compile the code execute execute the following command in the project root directory:
```bash
go build -o atm cmd/atm/main.go
```

A file called `./atm` will appear in your current working directory to run this file simply execute it:
```bash
./atm
```

## Testing
To run all tests in the suite execute the following in the project root directory:
```bash
go test
```

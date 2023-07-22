# Sneaker BOT for Footlocker USA

This repository contains a Sneaker BOT designed for Footlocker USA, which allows users to automate the process of purchasing sneakers from the Footlocker website. The BOT is built using Go programming language and is organized into different modules and components for better maintainability and extensibility.

## Below is the project structure for the Sneaker BOT for Footlocker USA:

```
- Footlocker
  - datadome
    - types.go
  - footlocker
    - types.go
  - logger
    - logger.go
  - shared
    - types.go
  - datadome-solution.go
  - datadome_payload.json
  - footlocker-module.go
  - proxies.txt
  - proxy-rotator.go
  - go.mod
  - go.sum
  - main.go
```

### Project Structure Explanation:

- **`Footlocker/`**: The root directory of the project.

  - **`datadome/`**: Package responsible for handling Datadome integration.

    - **`types.go`**: Defines data structures and types specific to the Datadome module.

  - **`footlocker/`**: Package containing functionalities related to Footlocker.

    - **`types.go`**: Defines data structures and types specific to the Footlocker module.

  - **`logger/`**: Package for logging events and actions during the BOT execution.

    - **`logger.go`**: Defines the logging behavior and implementation.

  - **`shared/`**: Shared package containing common data structures and functions used across modules.

    - **`types.go`**: Defines shared data structures and types.

  - **`datadome-solution.go`**: File containing the solution to bypass Datadome bot detection.

  - **`datadome_payload.json`**: JSON file containing payload data required for Datadome integration.

  - **`footlocker-module.go`**: The main module containing core functionalities of the Footlocker BOT.

  - **`proxies.txt`**: Text file containing a list of proxies that can be used for anonymous purchase attempts.

  - **`proxy-rotator.go`**: Implements the logic for rotating through the list of proxies.

  - **`go.mod`**: Go module file specifying the dependencies and their versions.

  - **`go.sum`**: Go checksum file containing the expected cryptographic hashes of the module sources.

  - **`main.go`**: The entry point of the application. Contains code to run the Footlocker Sneaker BOT.

## Usage

### Prerequisites

To run the Footlocker Sneaker BOT, you need to have the following installed on your system:

- Go programming language (version specified in `go.mod`)
- Internet connection

### Configuration

Before running the BOT, you need to set up the necessary configurations:

1. **`proxies.txt`:** Open the `proxies.txt` file and add the list of proxies you want to use for the bot. Each proxy should be on a separate line. The BOT will rotate through these proxies to prevent IP blocking.

2. **`main.go`:** In the `main.go` file, locate the variable `ProductURL` and set it to the URL of the sneaker product page you want to purchase from Footlocker USA. For example:

   ```go
   // URL of the sneaker product page on Footlocker USA
   ProductURL := "https://www.footlocker.com/us/product/sneaker_product_page"

   // and make sure to edit the task fields required to purchase from FTL. e.g: Cnb, Cardname, Email etc... :
   task := shared.Task{
   	ProfileName:   "ASM",
   	ProductURL:    "https://www.footlocker.com/product/~/38019001.html",
   	Size:          "10.0",
   	ProductID:     "",
   	Quantity:      1,
   	UseProxy:      false,
   	Mode:          "",
   	Aco:           false,
   	Region:        "NY",
   	RegionIsocode: "US-NY",
   	Store:         "",
   	Keywords:      "",
   	Sku:           "37581101",
   	Payment:       "",
   	Profile: shared.Profile{
   		ProfileName: "ASM Profile",
   		FirstName:   "ASM",
   		LastName:    "DEV",
   		Age:         1,
   		BirthDay:    1,
   		BirthMonth:  1,
   		BirthYear:   1,
   		Gender:      "",
   		Email:       "email",
   		Phone:       "2025961737",
   		Address:     "83 Pendergast Street",
   		Address2:    "",
   		Zip:         "11225",
   		City:        "BROOKLYN",
   		Country:     "",
   		CountryISO:  "US",
   		CountryCode: "",
   		State:       "",
   		Cardname:    "Kai Avila",
   		Cnb:         "4109703255583065",
   		Month:       "01",
   		Year:        "2027",
   		Cvv:         "340",
   		CardType:    "visa",
   		Password:    "",
   	},
   	Id: 0,
   }
   ```

### Running the BOT

Once you have configured the proxies and set the `ProductURL`, you can run the Footlocker Sneaker BOT. Open your terminal or command prompt, navigate to the root directory of the project, and execute the following command:

```
go run .
```

The BOT will initiate the automated process of purchasing sneakers from the specified Footlocker USA product page. It will handle the bot detection mechanisms using Datadome integration and automatically rotate through the list of proxies from `proxies.txt` to ensure a seamless and anonymous purchasing experience.

## License

This Sneaker BOT project is licensed under the [MIT License](LICENSE). You are free to modify, distribute, and use the code as per the terms mentioned in the license.

---

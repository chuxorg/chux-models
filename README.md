[![Test chux-models](https://github.com/chuxorg/chux-models/actions/workflows/build_test.yaml/badge.svg)](https://github.com/chuxorg/chux-models/actions/workflows/build_test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/chuxorg/chux-models)](https://goreportcard.com/report/github.com/chuxorg/chux-models)
# chux-models


## Overview
`chux-models` is a `Golang` library designed to provide services with an easy-to-use and consistent interface for performing CRUD operations on various data models, 
such as Products, Articles, Users, and future models. The library enforces the business logic of the chux-* ecosystem and ensures data is validated before being persisted 
into a data store. `chux-models` is dependent on the [`chux-datastore`](https://github.com/chuxorg/chux-datastore) library for handling the connection and read/writes to/from MongoDB.

## Features
- Consistent and straightforward interface for CRUD operations
- Built-in data validation and enforcement of business logic
- Support for multiple data models, including Products, Articles, Users, and more
- Seamless integration with the chux-mongo library for MongoDB connectivity
= Designed for use in the chux-* ecosystem
## Installation
To install chux-models, use the following command:

```sh
go get github.com/chuxorg/chux-models
```
## Usage
To use `chux-models` in your project, import the library and create instances of the desired data models. The library provides functions for creating, updating, deleting, and loading instances of the supported models.

Here's a simple example of how to use the Product model:

```go
package main

import (
    "github.com/chuxorg/chux-models/models"
)

func main() {
    // Create a new Product instance
    product, err := models.NewProduct()
    if err != nil {
        panic(err)
    }

    // Set the Product's properties
    product.Name = "Example Product"
    product.Description = "This is an example product."

    // Save the Product to the data store
    err = product.Save()
    if err != nil {
        panic(err)
    }
}

```
# Makefile

- `make test` - Runs all tests in `chux-models`.
- `make test-models` - Only runs tests in the Models package.
- `make test-release` - Commits, tags, and releases `chux-models`  
   &nbsp; 
   To release and version, pass in major and minor values on the command line. If either major or minor has a value, the patch number is set to zero. If neither major nor minor has a value, the patch number is incremented by 1.
   You can run the target with or without the MAJOR_VALUE and MINOR_VALUE variables, like this:
   &nbsp; 
   To bump the patch version:
   ```shell
   make release-version
   ```
   &nbsp; 
   To set a new major version and reset minor and patch to zero:
   ```shell
   MAJOR_VALUE=2 make release-version
   ```
   &nbsp; 
   To set a new minor version and reset patch to zero:
   ```shell
   MINOR_VALUE=3 make release-version 
   ```
   &nbsp; 
   To set both new major and minor versions and reset patch to zero:
   ```shell
   MAJOR_VALUE=2 MINOR_VALUE=3 make release-version 
   ``` 
## Contributing
Contributions to `chux-models` are welcome! Please submit an issue or a pull request if you have any ideas or suggestions for improvements.



## Methodology
### TDD
Test-driven development (TDD) is used while developing this shared library. TDD involves writing tests for your code before you write the actual implementation. 

This methodology has several advantages:

- Improved code quality: Writing tests first ensures that your code is testable and encourages modular design, making it easier to maintain and refactor in the future.
- Faster development: Once you have a comprehensive test suite, you can confidently make changes to your code without worrying about breaking existing functionality. Tests serve as a safety net, catching any regressions.
- Better collaboration: Tests can act as documentation, showing how your code is expected to behave and be used. This can help other developers understand your code more quickly.
- Easier debugging: When a bug is discovered, you can write a test that reproduces the issue before fixing it. This ensures that the bug won't reappear unnoticed in the future.

To adopt TDD for chux-models library, the following steps where followed:

1. Organize your code into packages and ensure that each package has a clear responsibility.
2. Write test cases for each function or method in your package, starting with simple test cases and gradually increasing complexity.
3. Run the tests and make sure they all pass before moving on to the next function or method.
4. For example, if you have a package called parser, create a parser_test.go file in the same package to hold your tests. 
5. Write tests for each function or method, and use a testing framework like the built-in testing package or a third-party library like github.com/stretchr/testify to make your tests more expressive and easier to write.

The overall goal of TDD is to ensure that chux-models is well-tested, maintainable, and easier for others to understand and use.

License
chux-models is released under the [GNU Public License v3](https://www.gnu.org/licenses/gpl-3.0.en.html)
&nbsp;
```go
There once was a library named chux,
Whose bizobj did strut and did tux,
With CRUD it did play,
And kept chaos at bay,
For devs in the chux-* matrix deluxe!
```
&nbsp;
<font size="2">Copyright Chuck Sailer &copy; 2023 - chux &trade;</font>

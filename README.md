# Token

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Build and Test](https://github.com/alexaandru/token/actions/workflows/ci.yml/badge.svg)](https://github.com/alexaandru/token/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexaandru/token/v2)](https://goreportcard.com/report/github.com/alexaandru/token/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexaandru/token/v2.svg)](https://pkg.go.dev/github.com/alexaandru/token/v2)

This is a simple package that generates randomized base62 encoded tokens based on an integer.
It's ideal for short url services or for any short, unique, randomized tokens you need to use throughout your app.

## Credits

This repo is a fork of [marksalpeter/token](https://github.com/marksalpeter/token).

Special thanks to [@einsteinx2](https://github.com/einsteinx2).
The encode and decode functions were ported from a short url project of his and he graciously
allowed [@marksalpeter](https://github.com/marksalpeter) to publish them.

Special thanks to [@sudhirj](https://github.com/sudhirj) for encorperating lexical sort order of the tokens into the package.

## Why Fork?

I originally forked it before the parent project adopted modules, so that I can add them.
Since then, I've cleaned up the code a tiny bit, got linters to pass and eliminated all
dependencies and reorganized the README a bit.

## How it Works

`Token` is an alias for `uint64`.

The `Token.Encode()` method returns a base62 encoded string based off of the uint64.
The string will always be in the same sort order as the uint64.

`Token` implements the `encoding.TextMarshaler` and `encoding.TextUnmarshaler` interfaces
to encode and decode to and from the base62 string representation of the `uint64`

Basically, the outside world will always see the token as a base62 encoded string,
but in your app you will always be able to use the token as a `uint64` for fast,
indexed, unique, lookups in various databases.

**IMPORTANT:** Remember to always check for collisions when adding randomized tokens to a database

### Example

```Go
package main

import (
	"fmt"
	"github.com/alexaandru/token/v2"
)

type Model struct {
    ID	token.Token `json:"id"`
}

func main() {
	// create a new model
	model := Model {
		ID:	token.New(), // creates a new, random uint64 token
	}
	fmt.Println(model.ID)          // 2751173559858
	fmt.Println(model.ID.Encode()) // Mr1NSSu

	// encode the model as json
	marshaled, err := json.Marshal(&model)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshaled)) // {"id":"Mr1NSSu"}

	// decode the model
	var unmarshaled Model
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		panic(err)
	}
	fmt.Println(unmarshaled.ID)    // 2751173559858

}
```

### References

- https://instagram-engineering.com/sharding-ids-at-instagram-1cf5a71e5a5c
- https://developer.twitter.com/en/docs/basics/twitter-ids.html
- https://github.com/ulid/spec

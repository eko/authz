<p align="center">
  <img height="125" alt="Fiber Schema" src="https://raw.githubusercontent.com/gofiber/schema/main/.github/logo-dark.svg#gh-dark-mode-only" />
  <img height="125" alt="Fiber Schema" src="https://raw.githubusercontent.com/gofiber/schema/main/.github/logo.svg#gh-light-mode-only" />
</p>

# gofiber/schema

Package gofiber/schema converts structs to and from form values.

Per-commit benchmark charts: <https://gofiber.github.io/schema/benchmarks/>

## Example

Here's a quick example: we parse POST form values and then decode them into a struct:

```go
// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var decoder = schema.NewDecoder()

type Person struct {
    Name  string
    Phone string
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        // Handle error
    }

    var person Person

    // r.PostForm is a map of our POST form values
    err = decoder.Decode(&person, r.PostForm)
    if err != nil {
        // Handle error
    }

    // Do something with person.Name or person.Phone
}
```

Conversely, contents of a struct can be encoded into form values. Here's a variant of the previous example using the Encoder:

```go
var encoder = schema.NewEncoder()

func MyHttpRequest() {
    person := Person{"Jane Doe", "555-5555"}
    form := url.Values{}

    err := encoder.Encode(person, form)

    if err != nil {
        // Handle error
    }

    // Use form values, for example, with an http client
    client := new(http.Client)
    res, err := client.PostForm("http://my-api.test", form)
}

```

To define custom names for fields, use a struct tag "schema". To not populate certain fields, use a dash for the name and it will be ignored:

```go
type Person struct {
    Name  string `schema:"name,required"`  // custom name, must be supplied
    Phone string `schema:"phone"`          // custom name
    Admin bool   `schema:"-"`              // this field is never set
}
```

The supported field types in the struct are:

* bool
* float variants (float32, float64)
* int variants (int, int8, int16, int32, int64)
* string
* uint variants (uint, uint8, uint16, uint32, uint64)
* struct
* a pointer to one of the above types
* a slice or a pointer to a slice of one of the above types
* a map with string keys whose values are one of the non-struct types above, a type with a registered converter or an `encoding.TextUnmarshaler` implementation, or a pointer or slice of those (see below)

Unsupported types are simply ignored, however custom types can be registered to be converted.

### Custom converters

`RegisterConverter` takes a `func(string) reflect.Value` that returns an invalid `reflect.Value` when the input does not convert. The decoder then reports a `ConversionError` without saying why. To report the reason, register a `func(string) (reflect.Value, error)` with `RegisterConverterE` instead: a non-nil error fails the conversion and is kept as the error's `Err`, so the message ends in its text, as it does for a failing `encoding.TextUnmarshaler`. `ConversionError` unwraps to `Err`, so `errors.Is` and `errors.As` reach it.

```go
decoder.RegisterConverterE(Level(0), func(s string) (reflect.Value, error) {
    n, err := strconv.Atoi(s)
    if err != nil {
        return reflect.Value{}, fmt.Errorf("level: %w", err)
    }
    return reflect.ValueOf(Level(n)), nil
})
// "level=high" fails with:
// schema: error converting value for "level". Details: level: strconv.Atoi: parsing "high": invalid syntax
```

A type has at most one converter: registering one of either kind replaces the other.

A struct type with a converter or an `encoding.TextUnmarshaler` implementation is decoded as a whole from its field's own key, and so is a multipart file field: a key below one, such as `link.host` or `avatar.filename`, is an unknown key.

### Map fields

A map field is filled from keys of the form `<field>.<key>`, one map entry per key:

```go
type Filter struct {
    Attributes map[string]string   `schema:"attributes"` // attributes.colour=red
    Tags       map[string][]string `schema:"tags"`       // tags.size=m&tags.size=l
    Limits     map[string]int      `schema:"limits"`     // limits.max=10
}
```

The field name matches case-insensitively like any other field, while the key is kept exactly as sent, and everything after the field name is the key (`attributes.a.b` gives the key `a.b`). Each value is decoded as a field of the map's value type would be, so a single-value map takes the last of several values. For a plain built-in value type (a string, bool or number), an empty value adds no entry unless `ZeroEmpty` is set; any other value type, such as a slice, a pointer or a type with a converter or `encoding.TextUnmarshaler`, gets an entry decoded from the empty value as such a field would be. A map is allocated on first use and entries are added to an existing one. A map type with a registered converter or an `encoding.TextUnmarshaler` implementation is decoded as a whole from the field's own key instead.

## Decoding Key/Value Pairs

When values arrive one key/value pair at a time, as a parsed query string's do, `DecodeValues` decodes the pairs without first building the map `Decode` takes:

```go
keys := []string{"name", "phone", "phone"}
values := []string{"Jane Doe", "555-5555", "555-0100"}

err := decoder.DecodeValues(&person, keys, values)
```

`values[i]` is a value of the key `keys[i]`, and a key given more than once has all of its values in the order given, as it would in the map. The two slices must be the same length.

## Limits

Keys usually come from a client, so the decoder bounds what they can make it allocate: a key may name at most 64 fields, and a slice index may be at most 1000 (or `MaxSize`, if lower). One call may also allocate at most 8 MiB in all for the values behind nil pointers on its keys' paths and for the slices it grows, counting the capacity each growth reserves. Maps and file fields are not counted, nor are the slice fields declared on the destination struct itself unless the struct embeds an exported pointer or the call grew eight other slices at a key's first index before them. The first eight slices a call grows at a key's first index grow geometrically, any other slice to exactly the index plus one. A key past a limit is not decoded and is reported under its own name in the returned `MultiError`, also when `IgnoreUnknownKeys` is set.

## Setting Defaults

It is possible to set default values when encoding/decoding by using the `default` tag option. The value of `default` is applied when a field has a zero value, a pointer has a nil value, or a slice is empty.

```go
type Person struct {
    Phone string `schema:"phone,default:+123456"`          // custom name
    Age int     `schema:"age,default:21"`
	Admin bool    `schema:"admin,default:false"`
	Balance float64 `schema:"balance,default:10.0"`
    Friends []string `schema:friends,default:john|bob`
}
```

The `default` tag option is supported for the following types:

* bool
* float variants (float32, float64)
* int variants (int, int8, int16, int32, int64)
* uint variants (uint, uint8, uint16, uint32, uint64)
* string
* a slice of the above types. As shown in the example above, `|` should be used to separate between slice items. 
* a pointer to one of the above types (pointer to slice and slice of pointers are not supported).

> [!NOTE]  
> The default is applied only when the field's key is missing from the source. A key provided as declared keeps its decoded value even if it is the zero value: `balance=0.0` leaves `Balance` at `0`. This also holds for fields promoted from an embedded struct. Use a pointer only when the handler itself needs to tell "not sent" from "sent as zero" after decoding.

<!-- skip-docs -->
## ☕ Supporters

Fiber is an open-source project that runs on donations to pay the bills, e.g., our domain name, hosting, and serverless infrastructure. If you want to support Fiber, please become a [GitHub Sponsor](https://github.com/sponsors/gofiber).

<p align="center">
  <a href="https://www.coderabbit.ai/?utm_source=gofiber&utm_medium=sponsor&utm_content=readme">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://www.coderabbit.ai/images/logo-dark.svg">
      <img width="280" height="52" alt="CodeRabbit" src="https://www.coderabbit.ai/images/logo-orange.svg">
    </picture>
  </a>
</p>
<p align="center">
  <a href="https://blacksmith.sh/?utm_source=gofiber&utm_medium=sponsor&utm_content=readme">
    <img width="280" height="96" alt="Blacksmith" src="https://raw.githubusercontent.com/gofiber/.github/main/assets/sponsors/blacksmith.png">
  </a>
</p>
<p align="center">
  <sub><b>Tool Sponsors</b> - supporting Fiber with free IDE licenses and AI credits</sub>
</p>
<p align="center">
  <a href="https://www.jetbrains.com/?from=gofiber" title="JetBrains - IDE licenses"><img width="36" height="36" alt="JetBrains" src="https://github.com/JetBrains.png?size=72"></a>
  &nbsp;
  <a href="https://openai.com/?utm_source=gofiber&utm_medium=sponsor&utm_content=readme" title="OpenAI - AI credits"><img width="36" height="36" alt="OpenAI" src="https://github.com/openai.png?size=72"></a>
  &nbsp;
  <a href="https://www.anthropic.com/?utm_source=gofiber&utm_medium=sponsor&utm_content=readme" title="Anthropic - AI credits"><img width="36" height="36" alt="Anthropic" src="https://github.com/anthropics.png?size=72"></a>
</p>

<!-- sponsors -->

### 📅 Monthly Sponsors

<table>
<tr><td valign="top"><strong>🔥 Fiber Guardian</strong></td><td><a href="https://www.coderabbit.ai/?utm_source=cr_org&amp;utm_medium=github" title="@coderabbitai"><img src="https://github.com/coderabbitai.png" width="50" alt="@coderabbitai" /></a></td></tr>
<tr><td valign="top"><strong>☕ Fiber Supporter</strong></td><td><a href="https://ndole.studio" title="@NdoleStudio"><img src="https://github.com/NdoleStudio.png" width="34" alt="@NdoleStudio" /></a></td></tr>
<tr><td valign="top"><strong>🪴 Fiber Friend</strong></td><td><a href="https://www.kamaslau.com/works?utm_source=github&amp;utm_medium=profile&amp;utm_campaign=always-on" title="@kamaslau"><img src="https://github.com/kamaslau.png" width="32" alt="@kamaslau" /></a>&nbsp;<a href="https://github.com/simonheisstpeter" title="@simonheisstpeter"><img src="https://github.com/simonheisstpeter.png" width="32" alt="@simonheisstpeter" /></a></td></tr>
</table>

### 🎁 One-time Sponsors

<table>
<tr><td valign="top"><strong>🚀 Fiber Hero</strong></td><td><a href="https://www.thanks.dev" title="@thnxdev"><img src="https://github.com/thnxdev.png" width="40" alt="@thnxdev" /></a></td></tr>
<tr><td valign="top"><strong>🪴 Fiber Friend</strong></td><td><a href="https://github.com/Gl1tchedPixzl" title="@Gl1tchedPixzl"><img src="https://github.com/Gl1tchedPixzl.png" width="26" alt="@Gl1tchedPixzl" /></a></td></tr>
</table>
<!-- sponsors -->
<!-- skip-docs -->

## License

BSD licensed. See the LICENSE file for details.

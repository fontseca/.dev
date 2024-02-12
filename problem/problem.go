package problem

import (
  "encoding/json"
  "fmt"
  "log/slog"
  "net/http"
  "net/url"
  "reflect"
  "strings"
  "unicode"
)

// baseURL contains metadata about how to ue the base URL.
type baseURL struct {
  url      *url.URL // url holds the parsed URL.
  fragment bool     // fragment indicates whether the URL will be used with fragments, instead of URL path.
  empty    bool     // empty indicates whether url is nil.
}

// set resets the baseURL to its initial state.
func (b *baseURL) set() {
  b.url = nil
  b.fragment = false
  b.empty = true
}

// base holds metadata about the default URL
// to be used when invoking Builder.Type method.
var base = baseURL{
  url:      nil,
  fragment: false,
  empty:    true,
}

// doParseURL parses the raw URL string and returns the parsed URL.
// It also returns a boolean flag indicating whether parsing was
// successful. If raw is an empty string, it returns a nil pointer
// to url.URL and false.
func doParseURL(raw string) (parsed *url.URL, ok bool) {
  raw, err := url.JoinPath(raw) // Sanitizes URL.
  if nil != err {
    slog.Error(err.Error())
    return nil, false
  }

  if '/' == raw[len(raw)-1] {
    raw = strings.TrimRight(raw, "/")
  }

  parsed, err = url.Parse(raw)
  if nil != err {
    slog.Error(err.Error())
    return nil, false
  }

  if "" == parsed.String() {
    return nil, false
  }

  return parsed, true
}

// SetGlobalURL sets the global base URI with the provided URI string.
// It also accepts an optional fragment flag to indicate if the URI
// will be used with fragments when invoking the builder's Type function.
// If an error occurs, then an empty URL is used.
func SetGlobalURL(raw string, fragment ...bool) {
  raw = strings.TrimSpace(raw)

  base.set()

  if 0 < len(fragment) {
    base.fragment = fragment[0]
  }

  if "" == raw {
    return
  }

  parsed, ok := doParseURL(raw)
  if ok {
    base.url = parsed
    base.empty = false
  }
}

// snakeToPascalCase converts a string from snake_case format to PascalCase.
// It capitalizes the first letter of each word separated by underscores and
// converts  the rest of the word to lowercase. Note: This implementation may
// not handle non-ASCII characters.
func snakeToPascalCase(s string) string {
  if "" == s {
    return ""
  }
  var f = func(c rune) bool { return '_' == c || unicode.IsSpace(c) }
  var words = strings.FieldsFunc(s, f)
  var builder strings.Builder
  for _, word := range words {
    builder.WriteRune(unicode.ToTitle(rune(word[0])))
    builder.WriteString(strings.ToLower(word[1:]))
  }
  return builder.String()
}

// canonicalSnakeCase returns the canonical snake_case format of the input string s.
// It converts the string s to lowercase and replaces any whitespace with underscores,
// ensuring it adheres to the canonical standard snake_case format. It does not
// perform any additional transformations; it assumes the input string is already in
// a non-canonical form of snake_case, such as "Snake_Case" or "snake case".
func canonicalSnakeCase(s string) string {
  if "" == s {
    return ""
  }
  return strings.ToLower(strings.Join(strings.Fields(s), "_"))
}

// Extension is a map that stores all the additional members that are
// specific to a problem type.
type Extension map[string][]any

// Add adds a new extension to the additional members of the problem.
func (e Extension) Add(key string, value any) {
  e[canonicalSnakeCase(key)] = append(e[canonicalSnakeCase(key)], value)
}

// Set sets the value of the extension for the specified key,
// replacing any existing values.
func (e Extension) Set(key string, value any) {
  e[canonicalSnakeCase(key)] = []any{value}
}

// Del deletes the extension for the specified key.
func (e Extension) Del(key string) {
  delete(e, canonicalSnakeCase(key))
}

// isValidHTTPStatusCode checks if the provided integer code is a valid HTTP status code.
// HTTP status codes are considered valid if they fall within the range of 100 to 599.
func isValidHTTPStatusCode(code int) bool {
  return code >= 100 && code <= 999
}

// Problem represents an RFC 9457 Problem Details object. When serialized in a JSON
// object, this format is identified with the "application/problem+json" media type.
// For more details, refer to RFC 9457: https://www.rfc-editor.org/rfc/rfc9457#name-the-problem-details-json-ob.
// It also implements the error interface for straightforward mobilization.
type Problem interface {
  error

  // Extension returns the Extension object associated with the
  // problem. If the Extension object is nil, then it's created.
  Extension() Extension

  // Emit sends the problem details as an HTTP response through the
  // provided http.ResponseWriter w. The response body is a JSON object
  // and its Content-Type is "application/problem+json".
  Emit(w http.ResponseWriter)
}

// problem implements the Problem interface. It serves to create new problems and
// as a base for custom problems.
type problem struct {
  // typ is a URI reference that identifies the problem type.
  typ string

  // status is the HTTP status code generated by the origin server for this occurrence of the problem.
  status int

  // title is a short, human-readable summary of the problem type.
  title string

  // detail is a human-readable explanation specific to this occurrence of the problem.
  detail string

  // instance is a URI reference that identifies the specific occurrence of the problem. It may or may not yield further information if dereferenced.
  instance string

  // extension contains additional members that are specific to a problem type.
  extension Extension
}

func (p *problem) Error() string {
  return ""
}

func (p *problem) Extension() Extension {
  if nil == p.extension {
    p.extension = Extension{}
  }
  return p.extension
}

// hasOneValueOnly checks if the slice of extensions for a key has only one element.
func (p *problem) hasOneValueOnly(values []any) bool {
  return 1 == len(values)
}

// makeStructFieldFor creates a reflect.StructField with the given name, sample type, and optional omitempty tag.
func (p *problem) makeStructFieldFor(name string, sample reflect.Type, omitempty ...bool) reflect.StructField {
  tagname := fmt.Sprintf(`json:"%s"`, canonicalSnakeCase(name))
  if 1 <= len(omitempty) && omitempty[0] {
    tagname = fmt.Sprintf(`json:"%s,omitempty"`, canonicalSnakeCase(name))
  }
  return reflect.StructField{
    Name: reflect.ValueOf(snakeToPascalCase(name)).Interface().(string),
    Type: sample,
    Tag:  reflect.StructTag(tagname),
  }
}

// appendExtensions appends the additional fields to the problem struct fields.
func (p *problem) appendExtensions(fields *[]reflect.StructField) {
  for extKey, extValues := range p.extension {
    var extType = reflect.TypeOf(extValues)
    if p.hasOneValueOnly(extValues) {
      extType = reflect.TypeOf(extValues[0])
    }
    *fields = append(*fields, p.makeStructFieldFor(extKey, extType))
  }
}

// setExtensionValues sets extension values to the given struct value.
func (p *problem) setExtensionValues(s reflect.Value) {
  for extKey, extValues := range p.extension {
    var value = reflect.ValueOf(extValues)
    if p.hasOneValueOnly(extValues) {
      value = reflect.ValueOf(extValues[0])
    }
    s.FieldByName(snakeToPascalCase(extKey)).Set(value)
  }
}

// setValuesToStructFields sets problem values to the given struct fields.
func (p *problem) setValuesToStructFields(s reflect.Value) {
  s.FieldByName("Type").SetString(p.typ)
  s.FieldByName("Status").SetInt(int64(p.status))
  s.FieldByName("Title").SetString(p.title)
  s.FieldByName("Detail").SetString(p.detail)
  s.FieldByName("Instance").SetString(p.instance)
  p.setExtensionValues(s)
}

// generateStruct generates a struct representing the problem.
func (p *problem) generateStruct() any {
  fields := []reflect.StructField{
    p.makeStructFieldFor("Type", reflect.TypeOf("")),
    p.makeStructFieldFor("Status", reflect.TypeOf(0)),
    p.makeStructFieldFor("Title", reflect.TypeOf("")),
    p.makeStructFieldFor("Detail", reflect.TypeOf(""), true),
    p.makeStructFieldFor("Instance", reflect.TypeOf(""), true),
  }
  p.appendExtensions(&fields)
  s := reflect.New(reflect.StructOf(fields)).Elem()
  p.setValuesToStructFields(s)
  return s.Interface()
}

// sanitize cleans up problem fields.
func (p *problem) sanitize() {
  if !isValidHTTPStatusCode(p.status) {
    p.status = http.StatusOK
  }
  p.detail = strings.TrimSpace(p.detail)
  p.instance = strings.TrimSpace(p.instance)
  p.typ = strings.TrimSpace(p.typ)
  p.title = strings.TrimSpace(p.title)
  if "" == p.typ {
    p.typ = "about:blank"
    if "" == p.title {
      p.title = http.StatusText(p.status)
    }
  }
}

// serialize converts the struct representation of the problem to JSON bytes.
func (p *problem) serialize(s any) []byte {
  data, err := json.Marshal(s)
  if nil != err {
    slog.Error(err.Error())
    return nil
  }
  return data
}

// doEmit writes the serialized problem to the http.ResponseWriter.
func (p *problem) doEmit(data []byte, w http.ResponseWriter) {
  w.Header().Set("Content-Type", "application/problem+json")
  w.WriteHeader(p.status)
  _, err := w.Write(data)
  if nil != err {
    slog.Error(err.Error())
  }
}

func (p *problem) Emit(w http.ResponseWriter) {
  if nil != w {
    p.sanitize()
    s := p.generateStruct()
    if serialized := p.serialize(s); nil != serialized {
      p.doEmit(serialized, w)
    }
  }
}

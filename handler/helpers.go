package handler

import (
  "encoding/json"
  "errors"
  "fmt"
  "fontseca/problem"
  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"
  "github.com/go-playground/validator/v10"
  "io"
  "log/slog"
  "math"
  "net/http"
  "reflect"
  "strconv"
  "strings"
  "testing"
)

func marshal(t *testing.T, value any) []byte {
  var data, err = json.Marshal(value)
  if nil != err {
    t.Fatal(err)
  }
  return data
}

func check(err error, w http.ResponseWriter) bool {
  if nil != err {
    var p *problem.Problem
    if errors.As(err, &p) {
      p.Emit(w)
    } else {
      problem.NewInternal().Emit(w)
    }
    return true
  }
  return false
}

type failure struct {
  Field     string `json:"field"`
  Criterion string `json:"criterion"`
  Parameter string `json:"parameter,omitempty"`
}

func bindJSONRequestBody(c *gin.Context, obj any) (ok bool) {
  c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
  var err = c.ShouldBindJSON(obj)
  if nil != err {
    var syntaxError *json.SyntaxError
    var unmarshalTypeError *json.UnmarshalTypeError
    var validationErrors validator.ValidationErrors

    var p problem.Problem
    p.Status(http.StatusBadRequest)

    switch {
    default:
      slog.Error(err.Error())
      problem.NewInternal().Emit(c.Writer)
      return
    case errors.Is(err, io.EOF):
      p.Title("Empty request body.")
      p.Detail("The request body must not be empty. Please provide a valid JSON object.")
    case strings.HasPrefix(err.Error(), "json: unknown field "):
      p.Title("Unexpected field in request body.")
      p.Detail("The request body contains an unexpected field. Please check the properties of the object.")
      field := strings.TrimPrefix(err.Error(), "json: unknown field \"")
      p.With("unexpected", strings.TrimSuffix(field, "\""))
    case errors.As(err, &validationErrors):
      p.Title("Invalid HTTP request body.")
      p.Status(http.StatusUnprocessableEntity)
      p.Detail("The provided JSON data does not meet the required validation criteria. Please review your input and try again.")
      var l = len(validationErrors)
      for _, e := range validationErrors {
        f := failure{
          Criterion: e.Tag(),
          Parameter: e.Param(),
          Field:     e.Field(),
        }
        if 1 == l {
          p.With("errors", []any{f})
        } else {
          p.With("errors", f)
        }
      }
    case errors.As(err, &syntaxError):
      p.Title("Malformed JSON in request body.")
      p.Detail("The request body contains invalid JSON syntax.")
      p.With("position", syntaxError.Offset)
      p.With("error", syntaxError.Error())
    case errors.As(err, &unmarshalTypeError):
      p.Title("Invalid value type in request body.")
      p.Detail("The request body contains a value that does not match the expected data type.")
      p.With("property", unmarshalTypeError.Field)
      p.With("has_type", unmarshalTypeError.Value)
      p.With("wants_type", unmarshalTypeError.Type.String())
    case errors.Is(err, io.ErrUnexpectedEOF):
      p.Title("Ill-formed JSON in request body.")
      p.Detail("The request body contains incomplete or truncated JSON data.")
    case err.Error() == "http: request body too large":
      p.Title("Request body too large.")
      p.Status(http.StatusRequestEntityTooLarge)
      p.Detail("The size of the request body must not exceed 1MB.")
    }
    p.Emit(c.Writer)
    return false
  }
  return true
}

func validateStruct(s any) error {
  var err = binding.Validator.ValidateStruct(s)
  if nil != err {
    var validationErrors validator.ValidationErrors
    if errors.As(err, &validationErrors) {
      var p problem.Problem
      p.Title("Failed to validate request data.")
      p.Status(http.StatusUnprocessableEntity)
      p.Detail("The provided data does not meet the required validation criteria. Please review your input and try again.")
      var l = len(validationErrors)
      for _, e := range validationErrors {
        f := failure{
          Criterion: e.Tag(),
          Parameter: e.Param(),
          Field:     e.Field(),
        }
        if 1 == l {
          p.With("errors", []any{f})
        } else {
          p.With("errors", f)
        }
      }
      err = &p
    } else {
      slog.Error(err.Error())
    }
  }
  return err
}

func handleStrconvError(err error, targetType, fieldName string) (error, bool) {
  if nil != err {
    var numErr *strconv.NumError
    if errors.As(err, &numErr) {
      switch {
      case errors.Is(err, strconv.ErrSyntax):
        return problem.NewUnparsableValue(targetType, fieldName, numErr.Num), false
      case errors.Is(err, strconv.ErrRange):
        return problem.NewValueOutOfRange(targetType, fieldName, numErr.Num), false
      }
    }
    return err, false
  }
  return nil, true
}

func bindPostForm(c *gin.Context, val any) error {
  if nil == c || nil == val {
    var err = errors.New("got an unacceptable nil parameter")
    slog.Error(err.Error())
    return err
  }
  var t = reflect.TypeOf(val)
  var wrongType = errors.New("type of parameter \"val\" is not a pointer to a struct")
  if reflect.Pointer != t.Kind() {
    slog.Error(wrongType.Error())
    return wrongType
  }
  var s = t.Elem()
  if reflect.Struct != s.Kind() {
    slog.Error(wrongType.Error())
    return wrongType
  }
  var v = reflect.ValueOf(val).Elem()

  for i := 0; i < s.NumField(); i++ {
    var fieldType = s.Field(i)
    if fieldType.IsExported() {
      var fieldValue = v.Field(i)
      if fieldValue.IsValid() && fieldValue.CanSet() {
        var fieldName = strings.Split(fieldType.Tag.Get("json"), ",")[0]
        var value, success = c.GetPostForm(fieldName)
        if !success {
          continue
        }
        value = strings.TrimSpace(value)
        var kind = fieldValue.Kind()
        switch {
        case reflect.String == kind:
          fieldValue.SetString(value)
        case reflect.Int == kind || reflect.Int8 == kind || reflect.Int16 == kind || reflect.Int32 == kind || reflect.Int64 == kind:
          var bitSize = 32
          var notInt = reflect.Int != kind
          if notInt {
            bitSize = int(math.Pow(2, float64(kind)))
          }
          var parsed, err = strconv.ParseInt(value, 10, bitSize)
          if nil != err {
            var targetType = "int"
            if notInt {
              targetType = fmt.Sprintf("int%d", bitSize)
            }
            if err, ok := handleStrconvError(err, targetType, fieldName); !ok {
              return err
            }
          }
          fieldValue.SetInt(parsed)
        case reflect.Float32 == kind || reflect.Float64 == kind:
          var parsed, err = strconv.ParseFloat(value, 128)
          if nil != err {
            var targetType = "float32"
            if reflect.Float64 == kind {
              targetType = "float64"
            }
            if err, ok := handleStrconvError(err, targetType, fieldName); !ok {
              return err
            }
          }
          fieldValue.SetFloat(parsed)
        case reflect.Bool == kind:
          var parsed, err = strconv.ParseBool(value)
          if err, ok := handleStrconvError(err, "bool", fieldName); !ok {
            return err
          }
          fieldValue.SetBool(parsed)
        }
      }
    }
  }
  return nil
}

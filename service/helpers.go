package service

import (
  "fontseca/problem"
  "net/http"
  "net/url"
  "strings"
)

func sanitizeURL(urls ...*string) error {
  if 0 == len(urls) {
    return nil
  }

  for _, u := range urls {
    if nil == u {
      continue
    }

    *u = strings.TrimSpace(*u)
    if "" == *u {
      continue
    }

    uri, err := url.ParseRequestURI(*u)
    if nil != err {
      var p problem.Problem
      p.Title("Unprocessable URL format.")
      p.Status(http.StatusUnprocessableEntity)
      p.Detail("There was an error parsing the requested URL. Please try with a different URL or verify the current one for correctness.")
      p.With("wrong_url", *u)
      return &p
    }
    *u = uri.String()
  }

  return nil
}

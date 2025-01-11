package main

type CookieBuilder interface {
	Encode(name string, value any) (string, error)
}

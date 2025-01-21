package main

type CookieManager interface {
	Encode(name string, value any) (string, error)
	Decode(name, value string, dst any) error
}

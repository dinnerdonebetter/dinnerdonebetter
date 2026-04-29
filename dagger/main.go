// Package main is the Dagger module that runs CI for dinnerdonebetter.
//
// Each public method on Ci corresponds to a single GitHub Actions workflow.
// Workflows in .github/workflows/ are thin shells that invoke `dagger call`
// against this module so the same pipeline runs identically locally and in CI.
package main

type Ci struct{}

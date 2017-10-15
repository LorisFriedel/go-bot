package router

import (
	"testing"
)

func TestStrictMatch(t *testing.T) {
	r, _ := RouteBuilder.Prefix("begin").Contains("middle,").Suffix("end").Pattern("-.-").Build()
	if !r.Match("begin---middle,end") {
		t.FailNow()
	}
}

func TestStrictMatchContainsError(t *testing.T) {
	r, _ := RouteBuilder.Prefix("begin").Contains("middl,").Suffix("end").Pattern("-.-").Build()
	if r.Match("begin---middle,end") {
		t.FailNow()
	}
}

func TestStrictMatchPrefixError(t *testing.T) {
	r, _ := RouteBuilder.Prefix("begn").Contains("middle,").Suffix("end").Pattern("-.-").Build()
	if r.Match("begin---middle,end") {
		t.FailNow()
	}
}

func TestStrictMatchSuffixError(t *testing.T) {
	r, _ := RouteBuilder.Prefix("begin").Contains("middle,").Suffix("ed").Pattern("-.-").Build()
	if r.Match("begin---middle,end") {
		t.FailNow()
	}
}

func TestStrictMatchPatternError(t *testing.T) {
	r, _ := RouteBuilder.Prefix("begin").Contains("middle,").Suffix("end").Pattern("-o-").Build()
	if r.Match("begin---end") {
		t.FailNow()
	}
}

func TestSoftMatch(t *testing.T) {
	r, _ := RouteBuilder.Prefix("begin").Contains("middle,").Suffix("end").Pattern("-.-").Soft().Build()
	if !r.Match("begin") {
		t.FailNow()
	}

	if !r.Match("---") {
		t.FailNow()
	}

	if !r.Match("end") {
		t.FailNow()
	}
}

func TestSoftMatchError(t *testing.T) {
	r, _ := RouteBuilder.Prefix("begin").Pattern("---").Suffix("end").Soft().Build()
	if r.Match("egin") {
		t.FailNow()
	}

	if r.Match("-O-") {
		t.FailNow()
	}

	if r.Match("en") {
		t.FailNow()
	}
}

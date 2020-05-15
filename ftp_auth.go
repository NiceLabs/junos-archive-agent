package main

type AnonymousAuth struct{}

func (*AnonymousAuth) CheckPasswd(string, string) (bool, error) { return true, nil }

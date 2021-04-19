package main

/*
登陆时发生的错误，包括用户名错误，密码错误
*/
type LoginError struct {
	message string
}

func (err LoginError) Error() string {
	return err.message
}

type JwtError struct {
	message string
}

func (err JwtError) Error() string {
	return err.message
}

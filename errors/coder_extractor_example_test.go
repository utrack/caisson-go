package errors_test

import (
	"database/sql"
	"fmt"

	"github.com/utrack/caisson-go/errors"
)

func ExampleCoder() {
	// in internal code

	internalCode := func() error {
		// define global error
		var ErrNotFound = errors.NewCoder("NOT_FOUND").WithHTTPCode(404).WithMessage("Not found (user message)")

		// in the function:
		err := sql.ErrNoRows

		return ErrNotFound.Wrap(err)
	}

	// in handler code
	{
		err := internalCode()
		err = errors.Wrap(err, "some more wrapping")

		errRsp := errors.Code(err)

		fmt.Println("Error:", err)
		fmt.Println("HTTP code:", errRsp.HTTPCode())
		fmt.Println("Message:", errRsp.Message())
		fmt.Println("Code:", errRsp.Type())
		fmt.Println("Unwrap is sql.ErrNoRows:", errRsp.Unwrap() == sql.ErrNoRows)
	}

	// Output:
	// Error: some more wrapping: sql: no rows in result set
	// HTTP code: 404
	// Message: Not found (user message)
	// Code: NOT_FOUND
	// Unwrap is sql.ErrNoRows: true
}

func ExampleDetailWith() {

	// in the internal package:
	// the error is public (ErrUserAlreadyExists);
	// the error details are private (errData).
	type errData struct {
		UserID int
	}

	var ErrUserAlreadyExists = errors.NewCoderDetailer[errData]("USER_ALREADY_EXISTS").WithHTTPCode(409).WithMessage("User already exists")

	// in internal/database layer
	sqlDoCreateUser := func() error {
		err := errors.New("sql: email already exists")
		d := errData{UserID: 31337}

		return ErrUserAlreadyExists.Wrap(err, d)
	}

	// in business-level code
	businessLevel := func() error {
		err := sqlDoCreateUser()

		details := ErrUserAlreadyExists.Details(err)
		if details != nil {
			// user already exists! handle it as such
			fmt.Println(err.Error())
			fmt.Println("user ID:", details.UserID)
			return err
		}
		if err != nil {
			// some other error
			fmt.Println("some unknown error: ", err.Error())
			return err
		}
		return nil
	}

	// in the handler/platform code
	handlerLevel := func() error {
		err := businessLevel()
		errRsp := errors.Code(err)
		fmt.Println("--- platform")
		fmt.Println("Error:", err)
		fmt.Println("HTTP code:", errRsp.HTTPCode())
		fmt.Println("Message:", errRsp.Message())
		fmt.Println("Code:", errRsp.Type())
		return err
	}

	_ = handlerLevel()
	// Output:
	// sql: email already exists
	// user ID: 31337
	// --- platform
	// Error: sql: email already exists
	// HTTP code: 409
	// Message: User already exists
	// Code: USER_ALREADY_EXISTS
}

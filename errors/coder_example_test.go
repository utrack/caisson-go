package errors_test

import (
	"database/sql"
	"fmt"

	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/caisson-go/levels/level3/errorbag"
)

// ExampleCoder_httpCodeSignaling shows how to use Coder to signal HTTP codes.
//
// The HTTP code can be picked up by the platform's handler to send it to the client.
func ExampleCoder_httpCodeSignaling() {
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

// ExampleCoderDetailer_customDetails shows how to use CoderDetailer to enrich errors with custom details.
// In this example, we use it to return additional context for the frontend; handler can pick up the details via [github.com/utrack/caisson-go/pkg/errorbag.ListPairs] and marshal them to the client.
func ExampleCoderDetailer_customDetails() {

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

		details := ErrUserAlreadyExists.ExtractDetail(err)
		if details != nil {
			// user already exists! handle it as such
			fmt.Println("business level error:", err.Error())
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
		fmt.Println("--- handler/platform-level:")
		fmt.Println("Error:", err)
		fmt.Println("HTTP code:", errRsp.HTTPCode())
		fmt.Println("Message:", errRsp.Message())
		fmt.Println("Code:", errRsp.Type())
		fmt.Println("Details:", errorbag.ListPairs(err))
		return err
	}

	_ = handlerLevel()
	// Output:
	// business level error: sql: email already exists
	// user ID: 31337
	// --- handler/platform-level:
	// Error: sql: email already exists
	// HTTP code: 409
	// Message: User already exists
	// Code: USER_ALREADY_EXISTS
	// Details: map[errors.Coded:sql: email already exists errors_test.errData:{31337}]
}

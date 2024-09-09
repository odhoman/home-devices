package common

import (
	"encoding/json"
	"fmt"
	constants "lambdas/common/constants"

	"github.com/aws/aws-lambda-go/events"
)

func ReturnOkResponse(code int, message string) map[string]interface{} {
	return map[string]interface{}{
		"statusCode": code,
		"message":    message,
	}
}

func ReturnOKWithMessageAPIGatewayProxyResponse(code int, message string) events.APIGatewayProxyResponse {
	return ReturnAPIGatewayProxyResponse(code, MessageResponse{message})
}

func ReturnAPIGatewayProxyResponse(code int, body interface{}) events.APIGatewayProxyResponse {

	jsonData, marshalError := json.Marshal(body)

	if marshalError != nil {
		fmt.Printf("Error converting message to JSON. Message: %v - Error: %v", body, marshalError)
		return ReturnDefaultInternalServerErrorResponse()
	}

	return createDefaultAPIGatewayProxyResponse(code, jsonData)

}

func ReturnInternalServerErrorAPIGatewayProxyResponse(errors []string) events.APIGatewayProxyResponse {
	return ReturnErrorResponseAPIGatewayProxyResponse(errors, 500)
}

func ReturnNotFoundErrorAPIGatewayProxyResponse(errors []string) events.APIGatewayProxyResponse {
	return ReturnErrorResponseAPIGatewayProxyResponse(errors, 404)
}

func ReturnBadRequestErrorAPIGatewayProxyResponse(errors []string) events.APIGatewayProxyResponse {
	return ReturnErrorResponseAPIGatewayProxyResponse(errors, 400)
}

func BadRequestErrorAPIGatewayProxyResponseSingleMessage(message string) events.APIGatewayProxyResponse {
	return ReturnBadRequestErrorAPIGatewayProxyResponse([]string{message})
}

func InternalServerErrorAPIGatewayProxyResponseSingleMessage(message string) events.APIGatewayProxyResponse {
	return ReturnInternalServerErrorAPIGatewayProxyResponse([]string{message})
}

func ReturnNotFoundErrorAPIGatewayProxyResponseSingleMessage(message string) events.APIGatewayProxyResponse {
	return ReturnNotFoundErrorAPIGatewayProxyResponse([]string{message})
}

func ReturnDefaultInternalServerErrorResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: constants.InternalServerErrorDefaultBodyResponse,
	}
}

func ReturnErrorResponseAPIGatewayProxyResponse(errors []string, code int) events.APIGatewayProxyResponse {
	jsonData, marshalError := json.Marshal(ReturnErrorResponse2(errors))

	if marshalError != nil {
		fmt.Printf("Error converting message to JSON. Error Message to convert: %v - Conversion Error: %v", errors, marshalError)
		return ReturnDefaultInternalServerErrorResponse()
	}

	return createDefaultAPIGatewayProxyResponse(code, jsonData)
}

func createDefaultAPIGatewayProxyResponse(code int, jsonData []byte) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonData),
	}
}

func ReturnErrorResponse2(errors []string) map[string]interface{} {
	return map[string]interface{}{
		"errors": errors,
	}
}

func ReturnErrorResponse(errors []string, code int) map[string]interface{} {
	return map[string]interface{}{
		"statusCode": code,
		"errors":     errors,
	}
}

func ReturnErrorResponseFromSingleError(error string, code int) map[string]interface{} {
	return ReturnErrorResponse([]string{error}, code)
}

func ReturnErrorResponse500(errors []string) map[string]interface{} {
	return ReturnErrorResponse(errors, 500)
}

func ReturnErrorResponse500SingleMessage(error string) map[string]interface{} {
	return ReturnErrorResponse([]string{error}, 500)
}

func ReturnErrorResponseNotFoundSingleMessage(error string) map[string]interface{} {
	return ReturnErrorResponse([]string{error}, 404)
}

func ReturnErrorResponseBadRequestdSingleMessage(error string) map[string]interface{} {
	return ReturnErrorResponse([]string{error}, 400)
}

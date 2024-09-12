**Home Devices**

This project aims to manage home devices. Implements Lambda functions in **Golang** to:
- Handle CRUD of home devices: create, update, delete, get
- SQS listener Lambda function to process messages from a queue Messages
will contain JSON payloads with information about device-home associations.
Upon receiving a message, update the corresponding device record in
DynamoDB with the new homeId information

Also implements the AWS infra necessary  in order to run the described above: 
- 5 lambda functions
- Api Gateway
- SQS queue
- DynamoDB

Provides necessary for testing, building the Lambda functions, and deploying the stack on AWS.

**Directory Structure**

The project is divided into three main folders:

- The **lambda** folder contains the Lambda functions (written in Golang) that will be executed for CRUD operations and updates on home devices.
- The **lib** folder contains the stack with the creation and configuration of AWS services.
- In the root folder, the **Makefile**  contains all the necessary comands for testing, building the Lambda functions, and deploying the stack on AWS.

**Makefile**

Below are the actions of Makefile:

- **`test_build_and_deploy_stack_all_lambdas`**: Test, build, and deploy the entire stack for all Lambdas.
- **`test_build_and_deploy_stack_for_single_lambda`**: Test, build, and deploy the stack for a single specified Lambda.
- **`test_and_build_all`**: Test and build all Lambdas sequentially.
- **`test_and_build_createDevice`**: Test and build only the `createDevice` Lambda.
- **`test_and_build_deleteDevice`**: Test and build only the `deleteDevice` Lambda.
- **`test_and_build_updateDevice`**: Test and build only the `updateDevice` Lambda.
- **`test_and_build_getDevice`**: Test and build only the `getDevice` Lambda.
- **`test_and_build_homeDeviceListener`**: Test and build only the `homeDeviceListener` Lambda.
- **`test_and_build_single_lambda`**: Test all Lambdas and build a single specified Lambda if tests pass.
- **`build_single_lambda`**: Build a single specified Lambda.
- **`test_all`**: Run tests for all Lambdas in the directory.
- **`run_tests_in_dir`**: Recursively run tests in all subdirectories.

**Operations Performed by the Lambda Functions**

***CreateDevice***

This function is responsible for creating a new device in the HomeDevices table in the DynamoDB database.

****Request Validations****

- **MAC (string) (json:"mac")**:
  - **Type**: String
  - **Validation**:
    - **Required**: This field is mandatory and must be provided.
    - **Min Length**: The MAC address must be at least 12 characters long.
    - **Max Length**: The MAC address cannot exceed 17 characters.
    - **Pattern Match**: The MAC address must conform to the expected MAC address format (e.g., XX:XX:XX:XX:XX:XX where XX represents hexadecimal digits).

- **Name (string) (json:"name")**:
  - **Type**: String
  - **Validation**:
    - **Required**: This field is mandatory and must be provided.
    - **Min Length**: The name must be at least 3 characters long.
    - **Max Length**: The name cannot exceed 50 characters.

- **Type (string) (json:"type")**:
  - **Type**: String
  - **Validation**:
    - **Required**: This field is mandatory and must be provided.
    - **Min Length**: The type must be at least 3 characters long.
    - **Max Length**: The type cannot exceed 20 characters.

- **HomeID (string) (json:"homeId")**:
  - **Type**: String
  - **Validation**:
    - **Required**: This field is mandatory and must be provided.
    - **Min Length**: The HomeID must be at least 5 characters long.
    - **Max Length**: The HomeID cannot exceed 30 characters.

**Unique Condition**

The combination of homeID and MAC must be unique within the table. It is not possible to create two devices with the same data.

**URL**

`POST https://q9n7bpmkr1.execute-api.us-east-1.amazonaws.com/prod/v1/device`

**Request - Response Examples**

- **Succeed Case**: Returns an HTTP 201 response with the complete data of the created device.

  **Example Request**:

  ```json
  {
    "mac": "0A:1B:2C:3D:4E:5F",
    "name": "Living Room Light",
    "type": "light",
    "homeId": "home2345"
  }
  ```

  **Example Response**:

  ```json
  {
    "id": "aab60d33-1188-4c70-8576-db11f0a65479",
    "mac": "0A:1B:2C:3D:4E:5F",
    "name": "Living Room Light",
    "type": "light",
    "homeId": "home2345",
    "createdAt": 1725940243,
    "modifiedAt": 1725940243
  }
  ```

- **Bad Request**:
  - **Validation Error**: Returns an HTTP 400 bad request error with validation errors for each field.

    ```json
    {
      "errors": [
        "Validation failed for field 'Name': required",
        "Validation failed for field 'Type': required"
      ]
    }
    ```

  - **Device Already Exists**: Returns an HTTP 400 bad request error indicating that the device already exists.

    ```json
    {
      "errors": [
        "Device Already Exist"
      ]
    }
    ```

- **Internal Server Error**: Returns a message indicating that there was an error creating the device.

  ```json
  {
    "errors": [
      "Internal Server error creating a new device"
    ]
  }
  ```

***UpdateDevice***

Modifies existing device information in the database.

**Request Validations**

- **MAC (string) (json:"mac")**:
  - **Type**: String
  - **Validation**:
    - **Optional**: This field is not required, but if provided, it must follow the validation rules.
    - **Min Length**: The MAC address must be at least 12 characters long.
    - **Max Length**: The MAC address cannot exceed 17 characters.
    - **Pattern Match**: The MAC address must conform to the expected MAC address format (e.g., XX:XX:XX:XX:XX:XX where XX represents hexadecimal digits).

- **Name (string) (json:"name")**:
  - **Type**: String
  - **Validation**:
    - **Optional**: This field is not required, but if provided, it must follow the validation rules.
    - **Min Length**: The name must be at least 3 characters long.
    - **Max Length**: The name cannot exceed 50 characters.

- **Type (string) (json:"type")**:
  - **Type**: String
  - **Validation**:
    - **Optional**: This field is not required, but if provided, it must follow the validation rules.
    - **Min Length**: The type must be at least 3 characters long.
    - **Max Length**: The type cannot exceed 20 characters.

- **HomeID (string) (json:"homeId")**:
  - **Type**: String
  - **Validation**:
    - **Optional**: This field is not required, but if provided, it must follow the validation rules.
    - **Min Length**: The HomeID must be at least 5 characters long.
    - **Max Length**: The HomeID cannot exceed 30 characters.

At least one of the fields described above must have a value.

**URL**

`PUT https://q9n7bpmkr1.execute-api.us-east-1.amazonaws.com/prod/v1/device/{id}`

**Request - Response Examples**

- **Succeed Case**: Returns an HTTP 201 response with a message indicating that the device was updated.

  **Example Request**:

  ```json
  {
    "mac": "0A:1B:2C:3D:4E:5F",
    "name": "Living Room Light",
    "type": "light",
    "homeId": "home4"
  }
  ```

  **Example Response**:

  ```json
  {
    "message": "Device updated"
  }
  ```

- **Bad Request**:
  - **Validation Error**: Returns an HTTP 400 bad request error with validation errors for each field.

    ```json
    {
      "errors": [
        "MAC address must be between 12 and 17 characters",
        "Type must be between 3 and 20 characters",
        "Home ID must be between 5 and 30 characters"
      ]
    }
    ```

- **Not Found**: Returns an HTTP 404 not found error indicating that the device was not found.

  ```json
  {
    "errors": [
      "Device Not Found"
    ]
  }
  ```

- **Internal Server Error**: Returns a message indicating that there was an error updating the device.

  ```json
  {
    "errors": [
      "Internal Server error updating a device"
    ]
  }
  ```

***DeleteDevice***

Removes a device from DynamoDB.

**URL**

`DELETE https://q9n7bpmkr1.execute-api.us-east-1.amazonaws.com/prod/v1/device/{id}`

**Request - Response Examples**

- **Succeed Case**: Returns an HTTP 200 response with a message indicating that the device was deleted.

  **Example Response**:

  ```json
  {
    "message": "Device deleted"
  }
  ```

- **Not Found**: Returns an HTTP 404 not found

 error indicating that the device was not found.

  ```json
  {
    "errors": [
      "Device Not Found"
    ]
  }
  ```

- **Internal Server Error**: Returns a message indicating that there was an error deleting the device.

  ```json
  {
    "errors": [
      "Internal Server error deleting a device"
    ]
  }
  ```

***GetDevice***

Retrieves details of a device based on a unique identifier.

**URL**

`GET https://q9n7bpmkr1.execute-api.us-east-1.amazonaws.com/prod/v1/device/{id}`

**Request - Response Examples**

- **Succeed Case**: Returns an HTTP 200 response with all the details of the device.

  **Example Response**:

  ```json
  {
    "id": "9a335b29-eec2-4dbc-8fc8-508f5433741e",
    "mac": "0A:1B:2C:3D:4E:5F",
    "name": "Living Room Light",
    "type": "light",
    "homeId": "home3",
    "createdAt": 1725971399,
    "modifiedAt": 1725971399
  }
  ```

- **Not Found**: Returns an HTTP 404 not found error indicating that the device was not found.

  ```json
  {
    "errors": [
      "Device Not Found"
    ]
  }
  ```

- **Internal Server Error**: Returns a message indicating that there was an error retrieving the device.

  ```json
  {
    "errors": [
      "Internal Server error getting the device"
    ]
  }
  ```

**UpdateDevice (SQS Listener)**

This Lambda function listens to SQS messages to process updates to device-home associations. Upon receiving a message, it updates the corresponding device record in DynamoDB with the new homeId information.

**SQS Message Validations**

- **ID (string) (json:"id")**:
  - **Type**: String
  - **Validation**:
    - **Required**: This field must be provided.

- **HomeID (string) (json:"homeId")**:
  - **Type**: String
  - **Validation**:
    - **Required**: This field must be provided.
    - **Min Length**: The HomeID must be at least 5 characters long.
    - **Max Length**: The HomeID cannot exceed 30 characters.

**Processing Examples**

- **Succeed Case**: The homeId is successfully updated in the database for the provided id.

**Errors**

- **Parsing Message**: The logs will indicate that the message from SQS could not be parsed.

- **Validation Error**: There was a validation error in one of the fields in the message from SQS.

- **Internal Server Error**: There was an error trying to update the homeId in the database.
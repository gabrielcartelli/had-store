# Backend Hat Store Training Program

This document provides an overview of the backend implementation for the Hat Store training program. It includes setup instructions, API endpoint descriptions, and other relevant information for developers working on the backend.

## Project Structure

The backend consists of the following components:

- **main.go**: Entry point of the application. Initializes the server and sets up routes.
- **handlers/hats.go**: Contains functions to handle HTTP requests related to hats, including product listing, cart operations, and checkout.
- **models/hat.go**: Defines the Hat struct, representing a hat product with properties like ID, Name, Price, and Description.
- **routes/routes.go**: Sets up the API routes for the application.

## Setup Instructions

1. **Install Go**: Ensure that you have Go installed on your machine. You can download it from the official Go website.

2. **Clone the Repository**: Clone the repository to your local machine using the following command:
   ```
   git clone <repository-url>
   ```

3. **Navigate to the Backend Directory**:
   ```
   cd hat-store-training/backend
   ```

4. **Install Dependencies**: If there are any dependencies, install them using:
   ```
   go mod tidy
   ```

5. **Run the Application**: Start the server by running:
   ```
   go run main.go
   ```

The server will start listening for incoming requests on the specified port.

## API Endpoints

- **GET /hats**: Retrieves a list of available hats.
- **POST /cart/add**: Adds a hat to the shopping cart.
- **POST /cart/update**: Updates the quantity of a hat in the shopping cart.
- **POST /checkout**: Processes the checkout, allowing selection of payment methods such as Pix or boleto.

## Contributing

Contributions to the backend are welcome. Please follow the standard Git workflow for submitting changes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
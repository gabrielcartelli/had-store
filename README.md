# Hat Store Training Program

Welcome to the Hat Store Training Program! This project is designed to help quality analysts understand the functionality of a hat store application, including product listings, shopping cart operations, and checkout processes.

## Project Structure

The project is divided into two main components: **backend** and **frontend**.

### Backend

The backend is developed in Go and includes the following files:

- **main.go**: The entry point of the application that initializes the server and sets up routes.
- **handlers/hats.go**: Contains functions to handle HTTP requests related to hats, including listing hats, adding to the cart, updating the cart, and processing checkout.
- **models/hat.go**: Defines the Hat struct, representing a hat product with properties like ID, Name, Price, and Description.
- **routes/routes.go**: Sets up the API endpoints for the application.

### Frontend

The frontend is developed using HTML, CSS, and JavaScript and includes the following files:

- **index.html**: The main HTML file that structures the webpage and links to CSS and JavaScript files.
- **css/styles.css**: Contains styles for the frontend application, defining layout, colors, and fonts.
- **js/app.js**: Handles user interactions, including adding hats to the cart, updating quantities, and managing the checkout process.

## Setup Instructions

### Backend

1. Navigate to the `backend` directory.
2. Run `go run main.go` to start the server.
3. The API will be available at `http://localhost:8080`.

### Frontend

1. Open `frontend/index.html` in a web browser.
2. The application will display the product listing and allow users to interact with the shopping cart.

## Usage

- Users can browse the available hats, add them to their shopping cart, adjust quantities, and proceed to checkout.
- Payment methods include Pix and boleto.

## Contributing

Feel free to contribute to the project by submitting issues or pull requests. Your feedback and improvements are welcome!

## License

This project is licensed under the MIT License.
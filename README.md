#   Bookshelf App - REST API & Frontend

A working version of this application can be found at [mybookshelfapp.com](http://www.mybookshelfapp.com)

This project is a Bookshelf app built with with Go and MongoDB (Backend), as well as HTML5/CSS3 (Frontend).  This app lets users build personalized lists of books using the Google Books API, and store these lists of books on a MongoDB server.  Features include:
- Near exclusive use of the Go Standard Library.  No third party frameworks.
- Full CRUD features for Book entries.
- Session tracking and authentication with Unique Universal Identifiers (UUIDs) written to browser Cookies.
- Hash encryption of user passwords
- Responsive frontend design for all browser sizes

### SETUP INSTRUCTIONS

To use this code you will require a [Google Developer Console API Key](https://console.developers.google.com/).  Sign-up is free and no credit card is required to access free-tier usage.

To compile this code you must be running Go version 1.11+.  In terminal type the following commands:
```
git clone https://github.com/pgmorgan/goSite.git
cd goSite
touch dev.env
vim dev.env
```
Insert the following lines in `dev.env`, resplacing all `<content>` with your own information:
```
PORT=<port number>
MONGODB_URL=<mongodb connection string>
GOOGLE_DEV_API_KEY=<api key>
```
To run the web server return to the root of the repository and type:
```
go run main.go
```
or alternatively
```
go build -o goSite
./goSite
```
<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="https://i.imgur.com/6wj0hh6.jpg" alt="Project logo"></a>
</p>

<h3 align="center">Go-Snip</h3>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![GitHub Issues](https://img.shields.io/github/issues/blue-davinci/go-snip.svg)](https://github.com/blue-davinci/go-snip/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/blue-davinci/go-snip.svg)](https://github.com/blue-davinci/go-snip/pulls)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

</div>

---

<p align="center"> A Fullstack Go + HTML + Vanilla JS Application for Snippet Generation and Sharing.
    <br> 
</p>

## 📝 Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Deployment](#deployment)
- [Usage](#usage)
- [Built Using](#built_using)
- [TODO](./TODO.md)
- [Contributing](../CONTRIBUTING.md)
- [Authors](#authors)
- [Acknowledgments](#acknowledgement)

## 🧐 About <a name = "about"></a>

      Go-Snip is a platform designed for sharing snippets of code, stories, posts, and more. It's a community where users can share their knowledge and learn from each other.
      With Go-Snip, you can easily save and organize your snippets for personal use or share them with the community. 
      Go-Snip also offers login capabilities, allowing users to create their own personal space for storing and managing snippets. You can keep your snippets private or share them with the Go-Snip community.

## 🏁 Getting Started <a name = "getting_started"></a>
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See [deployment](#deployment) for notes on how to deploy the project on a live system.
     
### Prerequisites

Before you begin, ensure you have met the following requirements:

- You have installed the latest version of [Go](https://golang.org/dl/), which is the programming language used in this project.
- You have installed [MySQL](https://dev.mysql.com/downloads/installer/), which is the database used in this project.
- You have installed [Git](https://git-scm.com/downloads), which is used for version control.

To confirm that you have installed Go and Git correctly, open a new terminal and run the following commands:

```bash
go version
git --version
```
### Installing <a name = "installing"></a>

Follow these steps to get a development environment running:

1. **Clone the repository** from GitHub:

    ```bash
    git clone https://github.com/blue-davinci/gosnip.git
    ```

2. **Navigate to the project directory**: Use the `cd` command to navigate to the projet directory:

    ```bash
    cd gosnip
    ```

3. **Install the project dependencies**: The Go tools will automatically download and install the dependencies listed in the `go.mod` file when you build or run the project. To download the dependencies without building or running the project, you can use the `go mod download` command

    ```bash
    go mod download
    ```

4. **Set up the database**. First, start your MySQL server. You can,
view the Migrations and execute them like so:

    ```bash
    migrate -path=migrations/test_db -database='mysql://`your_username`:`your_password`@/`your_database`?parseTime=true' up
    ```
  Migration tool used :- `golang migrate`
  If you need to understand them, you can alsways open them and check what they do. Database creation, table creations and sample data insertions.

5. **Build the project**: You can build the project using the `go build` command:

    ```bash
    go build
    ```
  An alternative is using the MAKEFILE, to just build immediately:
  ```bash
make build/api
```


6. **Run the project**: You can run the project using the `go run` after navigating to `cmd\web\` directory:

    ```bash
    go run ./cmd/web
    ```
  Alternatively, run the make commands for the same:
  ```bash
    make build/exe
  ```

Remember to replace `your_username` and `your_database` with your actual MySQL username and the name of the database you want to create. If your MySQL server is password-protected, you'll be prompted to enter your password after running the `mysql` commands.

## 🔧 Running the tests <a name = "tests"></a>

This project uses Go's built-in testing package. To run the tests, follow these steps:

1. Open a terminal in the project's root directory.

2. Run the following command to execute the tests:

    ```bash
    go test ./...
    ```

  This command will recursively run all tests in the project.

  To run tests for the `handlers` package, you would use:

    ```bash
    go test ./cmd/web
    ```

### Break down into end to end tests <a name = "test_break"></a>

End to end tests are types of tests that simulate a user's journey through encompassing the routimg, middleware and handlers within our application.

In our application, end-to-end testing is performed using the httptest.NewTLSServer() function. This function creates an instance of httptest.Server that allows us to make HTTPS requests, effectively simulating a real server for testing purposes.

We start by creating a new instance of our application struct, which currently includes the following mock application dependancies:
```
		logger:         a mock logger,
		snippets:       a mock interface pointing to our snippet model mock
		users:          a mock interface pointing to our user model mock
		templateCache:  holds a cache of our templates,
		formDecoder:    a data object that holds our form data,
		sessionManager: our session manager,
```

The test server's network address is stored in the ts.URL field. We use this address to make HTTP requests against the server. For example, we can use ts.Client().Get() to send a GET request to the /ping endpoint. This method returns an http.Response struct, which contains the server's response.

```
rs, err := ts.Client().Get(ts.URL + "/ping")
if err != nil {
    t.Fatal(err)
}

// We can then check the value of the response status code and body 
assert.Equal(t, rs.StatusCode, http.StatusOK)

defer rs.Body.Close()

body, err := io.ReadAll(rs.Body)
if err != nil {
    t.Fatal(err)
}

body = bytes.TrimSpace(body)
assert.Equal(t, string(body), "OK")
```

***Sample Output**: sample output of the above test:
```
=== RUN   TestPing
Running:  TestPing
{"level":"INFO","time":"0000-00-00Z90:00:000","message":"127.0.0.1:55504 - HTTP/1.1 GET /ping"}
--- PASS: TestPing (0.01s)
PASS
ok      github.com/blue-davinci/gosnip/cmd/web  0.972s
```

### Test Structure

- The individual tests can be found as files ending with `*_test.go`
- `testutils_test.go` holds utility functions that we can are using for our tests
- `internal\models\mocks` holds the mock models for our snippets and user code.

## 🎈 Usage <a name="usage"></a>

Once you have the application up and running (refer to the [installation instructions](#installing)), you can start using it by navigating to `localhost:<port>` in your web browser, where `<port>` is the port number the application is running on.

Here's a brief overview of the application's functionality:

- **Account Creation**: Click on the 'Sign Up' button to create a new account. You'll need to provide some basic information, such as your username and password.

- **Logging In**: Once you've created an account, you can log in by clicking on the 'Log In' button and entering your account details.

- **Creating Snippets**: After logging in, you can create a new snippet by clicking on the 'Create Snippet' button. You'll need to provide a title and the content for your snippet.

- **Viewing Account Details**: You can view your account details by clicking on your username in the top-right corner of the screen.

- **Viewing Snippets**: All users, including those who aren't logged in, can view snippets that have been posted by navigating to the 'Snippets' page.

Remember to replace `<port>` with the actual port number that your application is running on.
- **Security Features:** Protections included are:
                    - CSRF Protection
                    - Session management
                    - XSS protection
                    - Self Signed TLS Certificates
                    - HTTPS server
                    - Connection Timeouts + Other HTTPS Settings

## 🚀 Deployment <a name = "deployment"></a>

{{To Do}}

## ⛏️ Built Using <a name = "built_using"></a>
- [Golang](https://golang.org/) - Language
- [MySQL](https://dev.mysql.com/) - Database
- [NodeJs](https://nodejs.org/en/) - Server Environment
- HTML + JS.. (refer to the [installation instructions](#installing))


## ✍️ Authors <a name = "authors"></a>

- [@ble_davinci](https://github.com/blue-davinci) - project


## 🎉 Acknowledgements <a name = "acknowledgement"></a>

- Hat tip to anyone whose code was used

## 📚 References <a name = "references"></a>
- [Go Documentation](https://golang.org/doc/): Official Go documentation and tutorials.
- [PostgreSQL Documentation](https://www.postgresql.org/docs/): Official PostgreSQL documentation.
- [Go database/sql tutorial](http://go-database-sql.org/): Tutorial on using Go's database/sql package.

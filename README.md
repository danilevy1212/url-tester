**url-tester**

A versatile Go application for checking the availability of URLs using HTTP HEAD requests.

**Description**

`url-tester` is a Go-based tool that verifies the accessibility of URLs either passed as an argument or read from a file. It performs HTTP HEAD requests and reports the status for each URL. It supports parallel processing for efficiency.

**Usage**

1. **Compile:**

   ```bash
   go build
   ```

1. **Basic usage help**

    ```bash
    ./url-tester --help   
    ```

1. **Run with a Single URL:**

    ```bash
    ./url-tester <url>
    ```

1. **Run with a File**

    ```bash
    ./url-tester --file <file_path>
    ```

    Replace <file_path> with the path to the file containing the list of URLs. Each URL should be on a new line.

**Output**

The program outputs the status of each URL. For URLs that return a status other than 200 OK, it prints the URL and its HTTP status code.

Example: `Status '404': 'http://example.com/notfound'`

**License**

- This project is licensed under the [MIT License](LICENSE).

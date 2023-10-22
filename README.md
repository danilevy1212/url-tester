**url-tester**

A simple Go program to check the status of URLs using HTTP HEAD requests.

**Description**

This Go project reads a list of URLs from a file and checks their availability using HTTP HEAD requests. If the URL returns a status other than 200 OK, it prints an error message.

**Usage**

1. **Compile:**

   ```
   go build
   ```

2. **Run:**

   ```
   ./url-tester <file_path>
   ```

   Replace `<file_path>` with the path to the file containing the list of URLs you want to check.

**Output**

- For each URL in the list, the program prints either a success message (`<URL>` was found: status 200) or an error message if the URL is inaccessible (`<URL>` was not found: status `<HTTP Status Code>`).

**License**

- This project is licensed under the [MIT License](LICENSE).

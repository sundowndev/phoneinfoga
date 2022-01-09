### Running a scan

Use the `scan` command with the `-n` (or `--number`) option.

```
phoneinfoga scan -n "+1 (555) 444-1212"
phoneinfoga scan -n "+33 06 79368229"
phoneinfoga scan -n "33679368229"
```

Special chars such as `( ) - +` will be escaped so typing US-based numbers stay easy : 

```
phoneinfoga scan -n "+1 555-444-3333"
```

!!! note "Note that the country code is essential. You don't know which country code to use ? [Find it here](https://www.countrycode.org/)"

<!--
#### Input & output file

Check several numbers at once and send results to a file.

```
phoneinfoga scan -i numbers.txt -o results.txt
```

Input file must contain one phone number per line. Invalid numbers will be skipped.

#### Footprinting

```
phoneinfoga scan -n +42837544833 -s footprints
```

#### Custom format reconnaissance

You don't know where to search and what custom format to use ? Let the tool try several custom formats based on the country code for you.

```
phoneinfoga recon -n +42837544833 
```
-->

## Available scanners

See page related to [scanners](scanners.md).

## Launching the server

PhoneInfoga integrates a REST API along with a web client that you can deploy anywhere. The API has been written in Go and web client in Vue.js. The application is stateless, so it doesn't require any persistent storage.

See **[API documentation](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/sundowndev/phoneinfoga/master/api/docs/swagger.yaml)**.

```shell
phoneinfoga serve
phoneinfoga serve -p 8080 # default port is 5000
```

You should then be able to visit the web client from your browser at `http://localhost:<port>`.

![](./images/screenshot.png)

**Running the REST API only**

You can choose to only run the REST API without the web client :

```
phoneinfoga serve --no-client
```

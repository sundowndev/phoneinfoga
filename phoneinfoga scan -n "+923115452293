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

PhoneInfoga embed a bunch of scanners that will provide information about the given phone number. Some of them will request external services, and so might require authentication. By default, unconfigured scanners won't run. The information gathered can then be used for a deeper manual analysis.

See page related to [scanners](scanners.md).

## Launching the web server

PhoneInfoga integrates a REST API along with a web client that you can deploy anywhere. The API has been written in Go and web client in Vue.js. The application is stateless, so it doesn't require any persistent storage.

See **[API documentation](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/sundowndev/phoneinfoga/master/web/docs/swagger.yaml)**.

```shell
phoneinfoga serve # uses default port 5000
phoneinfoga serve -p 8080 # use port 8080
```

Equivalent commands via docker:

```shell
docker run --rm -it -p 5000:5000 sundowndev/phoneinfoga serve # same as `phoneinfoga serve`
docker run --rm -it -p 8080:8080 sundowndev/phoneinfoga serve -p 8080 # same as `phoneinfoga serve -p 8080`
```

You should then be able to visit the web client from your browser at `http://localhost:<port>`.

![](./images/screenshot.png)

**Running the REST API only**

You can choose to only run the REST API without the web client:

```shell
phoneinfoga serve --no-client
```

Equivalent docker command:

```shell
docker run --rm -it -p 5000:5000 sundowndev/phoneinfoga serve --no-client
```

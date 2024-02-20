# Scanners

PhoneInfoga provide several scanners to extract as much information as possible from a given phone number. Those scanners may require authentication, so they're automatically skipped when no authentication credentials are found.

## Configuration

Note that all scanners use environment variables for configuration values. You can define an environment variable inline or put them in a file called `.env` in the current directory. The tool will parse it automatically. To specify another filename, use the flag `--env-file`.

**Example**

```shell
# .env.local
NUMVERIFY_API_KEY="value"
GOOGLECSE_CX="value"
GOOGLE_API_KEY="value"
```

```shell
phoneinfoga scan -n +4176418xxxx --env-file=.env.local
```

### Scanner options

When using the **REST API**, you can also specify those values on a per-request basis. Each scanner supports its own options, see below. For details on how to specify those options, see [API docs](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/sundowndev/phoneinfoga/master/web/docs/swagger.yaml#/Numbers/RunScanner). For readability and simplicity, options are named exactly like their environment variable equivalent.

!!! warning
    Scanner options will override environment variables for the current request.

## Building your own scanner

PhoneInfoga can now be extended with plugins! You can build your own scanner and PhoneInfoga will use it to scan the given phone number.

```shell
$ phoneinfoga scan -n +4176418xxxx --plugin ./custom_scanner.so
```

!!! info
    Plugins are written with the [Go programming language](https://golang.org/). To get started, [see this example plugin](https://github.com/sundowndev/phoneinfoga/tree/master/examples/plugin).

## Local

The local scan is probably the simplest scan of PhoneInfoga. By default, the tool statically parse the phone number and convert it to several formats, it also tries to recognize the country and the carrier. This information are passed to all scanners in order to provide further analysis. The local scanner simply return those information to the end user, so they can exploit it as well.

??? info "Configuration"

    There is no configuration required for this scanner.

??? example "Output example"

    ```shell
    $ phoneinfoga scan -n +4176418xxxx
    
    Results for local
    Raw local: 076418xxxx
    Local: 076 418 xx xx
    E164: +4176418xxxx
    International: 4176418xxxx
    Country: CH
    ```

## Numverify

Numverify provide standard but useful information such as country code, location, line type and carrier. This scanners requires an API-key which you can get on their website after creating an account. You can use a free API key as long as you don't exceed the monthly quota. **This is an [apilayer](https://apilayer.com/marketplace/number_verification-api) key, not numverify itself.**

[Read documentation](https://apilayer.com/marketplace/number_verification-api#details-tab)

??? info "Configuration"

    1. Go to the [Api layer website](https://apilayer.com/) and create an account
    2. Go to "Number Verification API" in the marketplace, click on "Subscribe for free", then choose whatever plan you want
    3. Copy the new API token and use it as an environment variable

    | Environment variable |   Option   | Default | Description                                          |
    |----------------------|------------|---------|-------------------------------------------------------|
    | NUMVERIFY_API_KEY    |   NUMVERIFY_API_KEY  |         | API key to authenticate to the Numverify API.        |

??? example "Output example"

    ```shell
    $ NUMVERIFY_API_KEY=<key> phoneinfoga scan -n +4176418xxxx
    
    Results for numverify
    Valid: true
    Number: 4176418xxxx
    Local format: 076418xxxx
    International format: +4176418xxxx
    Country prefix: +41
    Country code: CH
    Country name: Switzerland (Confederation of)
    Carrier: Sunrise Communications AG
    Line type: mobile
    ```

## Googlesearch

Googlesearch uses the Google search engine and [Google Dorks](https://en.wikipedia.org/wiki/Google_hacking) to search phone number's footprints everywhere on the web. It allows you to search for scam reports, social media profiles, documents and more. **This scanner does only one thing:** generating several Google search links from a given phone number. You then have to manually open them in your browser to see results. So the tool may generate links that do not return any result. This is a design choice we made to avoid technical limitation around [Google scraping](https://en.wikipedia.org/wiki/Search_engine_scraping).

You can however, use this scanner through the REST API in addition with another tool to fetch the result automatically. If you wish to retrieve results automatically, see [Googlecse scanner](#googlecse) instead.

??? info "Configuration"

    There is no configuration required for this scanner.

??? example "Output example"

    ```shell
    $ phoneinfoga scan -n +4176418xxxx
    
    Results for googlesearch
    Social media:
        URL: https://www.google.com/search?q=site%3Afacebook.com+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Atwitter.com+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Alinkedin.com+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Ainstagram.com+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Avk.com+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    Disposable providers:
        URL: https://www.google.com/search?q=site%3Ahs3x.com+intext%3A%224176418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Areceive-sms-now.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Asmslisten.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Asmsnumbersonline.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Afreesmscode.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Acatchsms.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Asmstibo.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Asmsreceiving.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Agetfreesmsnumber.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Asellaite.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Areceive-sms-online.info+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Areceivesmsonline.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Areceive-a-sms.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Asms-receive.net+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Areceivefreesms.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Areceive-sms.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Areceivetxt.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Afreephonenum.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Afreesmsverification.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Areceive-sms-online.com+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Asmslive.co+intext%3A%224176418xxxx%22+OR+intext%3A%22076418xxxx%22
    Reputation:
        URL: https://www.google.com/search?q=site%3Awhosenumber.info+intext%3A%22%2B4176418xxxx%22+intitle%3A%22who+called%22
    
        URL: https://www.google.com/search?q=intitle%3A%22Phone+Fraud%22+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Afindwhocallsme.com+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%224176418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Ayellowpages.ca+intext%3A%22%2B4176418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Aphonenumbers.ie+intext%3A%22%2B4176418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Awho-calledme.com+intext%3A%22%2B4176418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Ausphonesearch.net+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Awhocalled.us+inurl%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Aquinumero.info+intext%3A%22076418xxxx%22+OR+intext%3A%224176418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Auk.popularphotolook.com+inurl%3A%22076418xxxx%22
    Individuals:
        URL: https://www.google.com/search?q=site%3Anuminfo.net+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Async.me+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Awhocallsyou.de+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Apastebin.com+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Awhycall.me+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Alocatefamily.com+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    
        URL: https://www.google.com/search?q=site%3Aspytox.com+intext%3A%22076418xxxx%22
    General:
        URL: https://www.google.com/search?q=intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22+OR+intext%3A%22076+418+xx+xx%22
    
        URL: https://www.google.com/search?q=%28ext%3Adoc+OR+ext%3Adocx+OR+ext%3Aodt+OR+ext%3Apdf+OR+ext%3Artf+OR+ext%3Asxw+OR+ext%3Apsw+OR+ext%3Appt+OR+ext%3Apptx+OR+ext%3Apps+OR+ext%3Acsv+OR+ext%3Atxt+OR+ext%3Axls%29+intext%3A%224176418xxxx%22+OR+intext%3A%22%2B4176418xxxx%22+OR+intext%3A%22076418xxxx%22
    ```

## Googlecse

Google custom search is a Google product allowing users to create Programmable Search Engines for programmatic usage.
This scanner takes an existing search engine you created to perform search queries on a given phone number.

Custom Search JSON API provides 100 search queries (~50 scans) per day for free. If you need more, you may sign up for billing in the API Console. **Additional requests cost $5 per 1000 queries (~500 scans), up to 10k queries per day (~5000 scans)**.

Follow the steps below to create a new search engine : 

1. Go to [GCP console](https://console.cloud.google.com/apis/api/customsearch.googleapis.com/metrics) and enable the custom search API.
2. Go to the [credentials page](https://console.cloud.google.com/apis/credentials) and create a new API token. You can restrict this token to the Custom Search API.
3. [Follow this link](https://programmablesearchengine.google.com/controlpanel/all) and click on "Add" to create a new search engine.
4. Fill the form and make sure you select "Search the entire web".
5. Use the Search Engine ID and the API token to configure the scanner as per the configuration tab below.

??? info "Configuration"

    |  Environment variable |  Option  | Default  | Description                                                 |
    |-----------------------|----------|----------|-------------------------------------------------------------|
    | GOOGLECSE_CX          |    GOOGLECSE_CX    |          | Search engine ID.            |
    | GOOGLE_API_KEY        |  GOOGLE_API_KEY |          | API key to authenticate to the Google API.  |
    | GOOGLECSE_MAX_RESULTS |          |   10     | Maximum results for each request. Each 10 results requires an additional request. This value cannot go above 100.  |

??? example "Output example"

    ```shell
    $ phoneinfoga scan -n +1241325xxxx
 
    Results for googlecse
    Homepage: https://cse.google.com/cse?cx=<redacted>
    Result count: 1
    Items:
        Title: Info about +1241325xxxx
        URL: https://example.com/1241325xxxx
    ```

## OVH

OVH, besides being a web and cloud hosting company, is a telecom provider with several VoIP numbers in Europe. Thanks to their API-key free REST API, we are able to tell if a number is owned by OVH Telecom or not.

??? info "Configuration"

    There is no configuration required for this scanner.

??? example "Output example"

    ```shell
    $ phoneinfoga scan -n +3336517xxxx

    Results for ovh
    Found: true
    Number range: 036517xxxx
    City: Abbeville
    ```

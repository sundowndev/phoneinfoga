# Scanners

PhoneInfoga provide several scanners to extract as much information as possible from a given phone number. Those scanners may require authentication, so they're automatically skipped when no authentication credentials are found. Note that all scanners use environment variables to find credentials.

## Local

The local scan is probably the simplest scan of PhoneInfoga. By default, the tool statically parse the phone number and convert it to several formats, it also tries to recognize the country and the carrier. Those information are passed to all scanners in order to provide further analysis. The local scanner simply return those information to the end user so they can exploit it as well.

=== "Configuration"

    There is no configuration required for this scanner.

=== "Example"

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

Numverify provide standard but useful information such as country code, location, line type and carrier. This scanners requires an API-key which you can get on their website after creating an account. You can use a free API key as long as you don't exceed the monthly quota.

=== "Configuration"

    | Environment variable | Default | Description                                                            |
    |----------------------|---------|------------------------------------------------------------------------|
    | NUMVERIFY_API_KEY    |         | API key to authenticate to the Numverify API.                          |
    | NUMVERIFY_ENABLE_SSL | false   | Whether to use HTTPS or plain HTTP for requests to the Numverify API.  |
=== "Example"

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
    Location:
    Carrier: Sunrise Communications AG
    Line type: mobile
    ```

## Googlesearch

Googlesearch uses the Google search engine and [Google Dorks](https://en.wikipedia.org/wiki/Google_hacking) to search phone number's footprints everywhere on the web. It allows you to search for scam reports, social media profiles, documents and more. **This scanner does only one thing:** generating several Google search links from a given phone number. You then have to manually open them in your browser to see results. So the tool may generate links that do not return any result. This is a design choice we made to avoid technical limitation around [Google scraping](https://en.wikipedia.org/wiki/Search_engine_scraping).

You can however, use this scanner through the REST API in addition with another tool to fetch the result automatically.

=== "Configuration"

    There is no configuration required for this scanner.

=== "Example"

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

## OVH

OVH, besides being a web and cloud hosting company, is a telecom provider with several VoIP numbers in Europe. Thanks to their API-key free REST API, we are able to tell if a number is owned by OVH Telecom or not.

=== "Configuration"

    There is no configuration required for this scanner.

=== "Example"

    ```shell
    $ phoneinfoga scan -n +3336517xxxx

    Results for ovh
    Found: true
    Number range: 036517xxxx
    City: Abbeville
    ```

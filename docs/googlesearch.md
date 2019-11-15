# Dealing with Google captcha

### Using Google API key

If you have a Google search console API key, all you have to do is to edit the file `config.py` and fill it with your credentials. If you leave fields empty, the tool will automatically use the default search feature described below.

#### How to create a Google Custom Search Engine API key and CX id

**CX id** : 

- Go to [https://cse.google.com/cse/create/new](https://cse.google.com/cse/create/new) to create a new search engine
- Fill the form with a fake domain site like `example.com`
- Select English as language
- Give any name to your search engine and click on Create button
- Go to [https://cse.google.com/cse/all](https://cse.google.com/cse/all) again and click on the search engine you just created.
- Select all entries in "Sites to search" and delete them
- Turn "Search the entire web" to ON
- Click on the "Search engine ID" button and copy your search engine id. This is the value for `google_cx_id` field in config.py file

**CSE API key** :

- Go to [https://console.developers.google.com/apis/credentials](https://console.developers.google.com/apis/credentials)
- Click on "Create credentials" and select API key
- Copy the API key and click on close button. This is the value for `google_api_key` field in the config.py file
- **Be sure to restrict the API key** to "Custom Search API"

### Using the webdriver

By default, PhoneInfo uses Selenium to handle Google search feature. When running OSINT scans, you will usually be blacklisted very easily by Google, which will ask the tool to complete a captcha. Nothing more simple, just complete the captcha that appears on the firefox window. Then press ENTER in the CLI to tell the tool it can continue the scanning process.

Still having issues with Google captcha ? Please [open an issue](https://github.com/sundowndev/PhoneInfoga/issues).
**Be careful, the cookie contain your IP address.**

#### Using Docker

When you run the tool with docker, the geckodriver and Selenium are also containerized. Make sure you launched all services with docker-compose.

After you successfully launched all services, you should have the following setup :

```
$ docker-compose ps
        Name                    Command           State            Ports         
---------------------------------------------------------------------------------
phoneinfoga             python phoneinfoga.py     Exit 0                         
phoneinfoga_firefox_1   /opt/bin/entry_point.sh   Up       0.0.0.0:5900->5900/tcp
selenium-hub            /opt/bin/entry_point.sh   Up       0.0.0.0:4444->4444/tcp
```

The web driver (`phoneinfoga_firefox_1`) is mounted on port 5900. It includes a VNC server which allows you to control the browser through Docker. To connect to the VNC server, download a VNC client and connect to `127.0.0.1` on port `5900` with password `secret`.

Here's an example of VNC URL :

```
vnc://127.0.0.1:5900
```

Also make sure you disabled any `read-only` or `view-only` mode so you can interact with the browser in order to complete the captcha.

##### VNC clients

- [Chicken, an open-source VNC client for MacOS X](https://sourceforge.net/projects/chicken/)
- [Vinagre, an open-source VNC client for GNOME (Linux)](https://github.com/GNOME/vinagre)

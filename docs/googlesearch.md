# Dealing with Google captcha

### Using Google API key

If you have a Google search console API key, all you have to do is to edit the file `config.py` and fill it with your credentials. If you leave fields empty, the tool will automatically use the default search feature described below.

#### How to create a Google Custom Search Engine API key and CX id

**CX id** : 

- Go to https://cse.google.com/cse/create/new to create a new search engine
- Fill the form with a fake domain site like `example.com`
- Select English as language
- Give any name to your search engine and click on Create button
- Go to https://cse.google.com/cse/all again and click on the search engine you just created.
- Select all entries in "Sites to search" and delete them
- Turn "Search the entire web" to ON
- Click on the "Search engine ID" button and copy your search engine id. This is the value for `google_cx_id` field in config.py file

**CSE API key** :

- Go to https://console.developers.google.com/apis/credentials
- Click on "Create credentials" and select API key
- Copy the API key and click on close button. This is the value for `google_api_key` field in the config.py file
- **Be sure to restrict the API key** to "Custom Search API"

### Using the webdriver

By default, PhoneInfo uses Selenium to handle Google search feature. When running OSINT scans, you will usually be blacklisted very easily by Google, which will ask the tool to complete a captcha. Nothing more simple, just complete the captcha that appears on the firefox window. Then press ENTER in the CLI to tell the tool it can continue the scanning process.

Still having issues with Google captcha ? Please [open an issue](https://github.com/sundowndev/PhoneInfoga/issues).
**Be careful, the cookie contain your IP address.**
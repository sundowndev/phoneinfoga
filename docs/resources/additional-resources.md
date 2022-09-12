# Additional resources

### Understanding phone numbers

- [whitepages.fr/phonesystem](http://whitepages.fr/phonesystem/)
- [Formatting-International-Phone-Numbers](https://support.twilio.com/hc/en-us/articles/223183008-Formatting-International-Phone-Numbers)
- [National_conventions_for_writing_telephone_numbers](https://en.wikipedia.org/wiki/National_conventions_for_writing_telephone_numbers)

### Open data

- [api.ovh.com/console/#/telephony](https://api.ovh.com/console/#/telephony)
- [countrycode.org](https://countrycode.org/)
- [countryareacode.net](http://www.countryareacode.net/en/)
- [directory.didww.com/area-prefixes](http://directory.didww.com/area-prefixes)
- [numinfo.net](http://www.numinfo.net/)
- [gist.github.com/Goles/3196253](https://gist.github.com/Goles/3196253)

## Footprinting

!!! info
    Both free and premium resources are included. Be careful, the listing of a data source here does not mean it has been verified or is used in the tool. Data might be false. Use it as an OSINT framework.

### Reputation / fraud

- scamcallfighters.com
- signal-arnaques.com
- whosenumber.info
- findwhocallsme.com
- yellowpages.ca
- phonenumbers.ie
- who-calledme.com
- usphonesearch.net
- whocalled.us
- quinumero.info

### Disposable numbers

- receive-sms-online.com
- receive-sms-now.com
- hs3x.com
- twilio.com
- freesmsverification.com
- freeonlinephone.org
- sms-receive.net
- smsreceivefree.com
- receive-a-sms.com
- receivefreesms.com
- freephonenum.com
- receive-smss.com
- receivetxt.com
- temp-mails.com
- receive-sms.com
- receivesmsonline.net
- receivefreesms.com
- sms-receive.net
- pinger.com (=> textnow.com)
- receive-a-sms.com
- k7.net
- kall8.com
- faxaway.com
- receivesmsonline.com
- receive-sms-online.info
- sellaite.com
- getfreesmsnumber.com
- smsreceiving.com
- smstibo.com
- catchsms.com
- freesmscode.com
- smsreceiveonline.com
- smslisten.com
- sms.sellaite.com
- smslive.co

### Individuals

- Facebook
- Twitter
- Instagram
- Linkedin
- True People
- Fast People
- Background Check
- Pipl
- Spytox
- Makelia
- IvyCall
- PhoneSearch
- 411
- USPhone
- WP Plus
- Thats Them
- True Caller
- Sync.me
- WhoCallsMe
- ZabaSearch
- DexKnows
- WeLeakInfo
- OK Caller
- SearchBug
- numinfo.net

### Google dork examples

```
insubject:"+XXXXXXXXX" OR insubject:"+XXXXX" OR insubject:"XXXXX XXX XXX"
insubject:"XXXXXXXXX" OR intitle:"XXXXXXXXX"
intext:"XXXXXXXXX" AND (ext:doc OR ext:docx OR ext:odt OR ext:pdf OR ext:rtf OR ext:sxw OR ext:psw OR ext:ppt OR ext:pptx OR ext:pps OR ext:csv OR ext:txt OR ext:html)
site:"hs3x.com" "+XXXXXXXXX"
site:signal-arnaques.com intext:"XXXXXXXXX" intitle:" | Phone Fraud"
```
